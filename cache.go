// Cache last rendered images

package main

import (
	"container/list"
)

// cacheSize is the number of images to keep in the cache
const cacheSize = 10

type cache struct {
	// List of *mandelImg (the cache itself)
	l *list.List
	// Channel to receive add-image requests from
	chAdd chan *mandelImg
	// Channel to receive lookup-image requests from
	chLookup chan lookupReq
}

// lookupReq is the cache-lookup request structure (send on
// cache.chLookup)
type lookupReq struct {
	// Image parameters
	p *params
	// Chan to send reply to
	ch chan *mandelImg
}

func (c *cache) match(p *params, m *mandelImg) bool {
	return p.Sx == m.Bounds().Dx() &&
		p.Sy == m.Bounds().Dy() &&
		p.X0 == real(m.C0) &&
		p.X1 == real(m.C1) &&
		p.Y0 == imag(m.C0) &&
		p.Y1 == imag(m.C1) &&
		p.Iter == m.MaxIter
}

func (c *cache) search(p *params) *mandelImg {
	for e := c.l.Front(); e != nil; e = e.Next() {
		ce := e.Value.(*mandelImg)
		if c.match(p, ce) {
			return ce
		}
	}
	return nil
}

func (c *cache) add(m *mandelImg) {
	p := params{
		Sx:   m.Bounds().Dx(),
		Sy:   m.Bounds().Dy(),
		X0:   real(m.C0),
		X1:   real(m.C1),
		Y0:   imag(m.C0),
		Y1:   imag(m.C1),
		Iter: m.MaxIter,
	}
	if c.search(&p) != nil {
		return
	}
	if c.l.Len() >= cacheSize {
		c.l.Remove(c.l.Front())
	}
	c.l.PushBack(m)
}

// ReqLookup requests a cache-lookup for an image with the given
// parameters. Returns a pointer to the image, if found in the cache,
// nil otherwise. When the cache is searched for an image, the image's
// palette is not considered: The returned image may have a different
// palette than the one specified in the request parameters. Because
// of this, you must always Repalette images received from the cache
// to make sure they are rendered with the correct palette.
func (c *cache) ReqLookup(p *params) *mandelImg {
	ch := make(chan *mandelImg)
	r := lookupReq{p, ch}
	c.chLookup <- r
	return <-ch
}

// ReqAdd requests that the given image is added to the cache. It is
// ok to request the addition of an image already in the cache
// (nothing happens in this case). The oldest cache entry may be
// evicted as a result of calling ReqAdd (if the cache size has
// reached cacheSize).
func (c *cache) ReqAdd(m *mandelImg) {
	c.chAdd <- m
}

// NewCache creates and initializes an image cache and starts the
// goroutine that receives and processes cache-lookup and cache-add
// requests. Returns a pointer to the newly created cache.
func newCache() *cache {
	c := new(cache)
	c.l = list.New()
	c.chLookup = make(chan lookupReq)
	c.chAdd = make(chan *mandelImg)
	go func(c *cache) {
		for {
			select {
			case lup := <-c.chLookup:
				lup.ch <- c.search(lup.p)
			case img := <-c.chAdd:
				c.add(img)
			}
		}
	}(c)

	return c
}
