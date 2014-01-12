// Calculate and render the mandelbrot set.

package main

import (
	"errors"
	"image"
	"image/color"
	"math/cmplx"
)

// MandelImg is a Mandelbrot-set image. It implements the image.Image
// interface.
type MandelImg struct {
	// Function domain
	C0, C1 complex128
	// Escape after MaxIter
	MaxIter int
	// Escape radius
	Radius float64
	// Palette used to map pixels to colors
	Palette color.Palette
	// Width & Height in pixels
	w, h int
	// Pixel array. Keeps iteration-count for every pixel
	pix []int
	// Histogram: histo[i] is # of pixels with i iterations
	histo []int
	// Cummulative-normalized histogram: cnhisto[i] is # of pixels
	// with iterations <= i, mapped to [0.0 .. 1.0]
	cnhisto []float64
}

// NewMandel calculates and returns a new Mandelbrot-set
// image. Returns non-nil error if invalid parameters are given.
func NewMandelImg(width, height int, p color.Palette, c0, c1 complex128,
	iter int, radius float64) (*MandelImg, error) {
	if iter <= 0 || radius <= 0 || width <= 0 || height <= 0 {
		err := errors.New("NewMandelImg: Invalid parameters")
		return nil, err
	}
	m := &MandelImg{}
	m.C0, m.C1 = c0, c1
	m.MaxIter = iter
	m.Radius = radius
	m.w = width
	m.h = height
	m.Palette = p
	m.pix = make([]int, width*height)
	m.histo = make([]int, iter+1)
	m.cnhisto = make([]float64, iter)
	m.calcPix()
	m.calcHisto()
	return m, nil
}

func (m *MandelImg) ColorModel() color.Model { return color.RGBAModel }

func (m *MandelImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, m.w, m.h)
}

func (m *MandelImg) At(x, y int) color.Color {
	if !m.pixIn(x, y) {
		return color.RGBA{}
	}
	iter := m.pix[m.pixOffset(x, y)]
	if iter == m.MaxIter {
		return m.Palette[0]
	} else {
		l := len(m.Palette)
		idx := int(m.cnhisto[iter] * float64(l-1))
		return m.Palette[idx]
	}
}

// Opaque scans the image's palette and returns true if all colors are
// fully opaque.
func (m *MandelImg) Opaque() bool {
	for _, c := range m.Palette {
		_, _, _, a := c.RGBA()
		if a != 0xffff {
			return false
		}
	}
	return true
}

// setIter sets the iteration count for the pixel at the given
// coordinates to "iter". It also updates the histogram by
// incrementing histo[iter].
func (m *MandelImg) setIter(x, y int, iter int) {
	if !m.pixIn(x, y) {
		return
	}
	of := m.pixOffset(x, y)
	if iter > m.MaxIter {
		iter = m.MaxIter
	} else if iter < 0 {
		iter = 0
	}
	m.pix[of] = iter
	m.histo[iter]++
}

// pixOffset returns the pix-array index of the pixel at the given
// coordinates.
func (m *MandelImg) pixOffset(x, y int) int {
	return y*m.w + x
}

// pixIn returns true if the pixel at the given coordinates is inside
// the image.
func (m *MandelImg) pixIn(x, y int) bool {
	return x >= 0 && x < m.w && y >= 0 && y < m.h
}

// calcPix calculates pixel values for the image as well as the image
// histogram.
func (m *MandelImg) calcPix() {
	// Deltas for stepping on the complex plane
	dx := (real(m.C1) - real(m.C0)) / float64(m.w)
	dy := (imag(m.C1) - imag(m.C0)) / float64(m.h)
	// x, y are on the complex plane (world coordinates)
	// px, py are on the image (viewport coordinates)
	for y, py := imag(m.C0), 0; py < m.h; y, py = y+dy, py+1 {
		for x, px := real(m.C0), 0; px < m.w; x, px = x+dx, px+1 {
			c := complex(x, y)
			z := complex(0, 0)
			var i = 0
			for i = 0; i < m.MaxIter; i++ {
				z = z*z + c
				if cmplx.Abs(z) > m.Radius {
					break
				}
			}
			m.setIter(px, py, i)
		}
	}
}

// calcHisto calculates the cumulative-normalized histogram for the
// image. The image histogram (the non-cumulative one) must have
// already been calculated before calling calcHisto (i.e. calcPix must
// be called before calcHisto).
func (m *MandelImg) calcHisto() {
	// Calculate cumulative histogram: cnhisto[i] is # of pixels
	// with iter-count <= i. Exclude pixels with iter-count ==
	// MaxIter (i.e. pixels IN the mandelbrot set) from the
	// calculation, for two reasons: (a) these pixels are always
	// colored with Palette[0], and (b) including them in cnhisto
	// would push all other pixels towards the low end of the
	// palette.
	m.cnhisto[0] = float64(m.histo[0])
	total := m.histo[0]
	for i := 1; i < m.MaxIter; i++ {
		m.cnhisto[i] = float64(m.histo[i]) + m.cnhisto[i-1]
		total += m.histo[i]
	}
	// Normalize cnhisto: map cnhisto[i] to range [0.0 .. 1.0]
	for i := 0; i < m.MaxIter; i++ {
		m.cnhisto[i] /= float64(total)
	}
}

// Repalette creates a copy of the image with a different palette and
// returns a pointer to it. The two images (the original and returned
// copy) share the same data (the same pixel array, and the same
// histogram arrays).
func (m *MandelImg) Repalette(p color.Palette) *MandelImg {
	mn := *m
	mn.Palette = p
	return &mn
}
