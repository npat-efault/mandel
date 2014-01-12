package main

import (
	"image/color"
	"testing"
)

func TestLinterp(t *testing.T) {
	var do = func(size int, s, e uint8) {
		p := make(color.Palette, size)
		p[0] = color.RGBA{s, s, s, s}
		p[size-1] = color.RGBA{e, e, e, e}
		linterp(p)
		delta := (int(e) - int(s)) / (size - 1)
		sd := 1
		if e < s {
			sd = -1
		}
		pc := p[0].(color.RGBA)
		var i, dc int
		var cc color.RGBA
		for i = 1; i < size; i++ {
			cc = p[i].(color.RGBA)
			dc = int(cc.R) - int(pc.R)
			if dd := (dc - delta) * sd; dd > 1 || dd < 0 {
				break
			}
			dc = int(cc.G) - int(pc.G)
			if dd := (dc - delta) * sd; dd > 1 || dd < 0 {
				break
			}
			dc = int(cc.B) - int(pc.B)
			if dd := (dc - delta) * sd; dd > 1 || dd < 0 {
				break
			}
			dc = int(cc.A) - int(pc.A)
			if dd := (dc - delta) * sd; dd > 1 || dd < 0 {
				break
			}
			pc = cc
		}
		if i < size {
			t.Errorf("linterp: size=%d, s=%d, e=%d",
				size, s, e)
			t.Errorf("%d:%v, %d:%v", i, cc, i-1, pc)
			t.Fatalf("delta=%d, dc=%d", delta, dc)
		}
	}
	do(2, 0, 1)
	do(3, 0, 2)
	do(4, 0, 11)
	do(32, 12, 111)
	do(255, 10, 128)
	do(255, 128, 10)
	do(100000, 255, 0)
	do(100000, 0, 255)
}

func TestPalGray256(t *testing.T) {
	for i := 0; i < 256; i++ {
		c := palGray256[i].(color.RGBA)
		i8 := uint8(i)
		if c.R != i8 || c.G != i8 || c.B != i8 || c.A != 0xff {
			t.Fatalf("%d:%v", i, c)
		}
	}
}

func TestPalGrayR256(t *testing.T) {
	for i := 255; i >= 0; i-- {
		c := palGrayR256[255-i].(color.RGBA)
		i8 := uint8(i)
		if c.R != i8 || c.G != i8 || c.B != i8 || c.A != 0xff {
			t.Fatalf("%d:%v", 255-i, c)
		}
	}
}

func TestLinGrad(t *testing.T) {
	cpts := []colPt{
		{0, color.RGBA{0, 0, 0, 0xff}},
		{255, color.RGBA{255, 0, 0, 0xff}},
		{256, color.RGBA{0, 0, 0, 0xff}},
		{511, color.RGBA{0, 255, 0, 0xff}},
		{512, color.RGBA{0, 0, 0, 0xff}},
		{767, color.RGBA{0, 0, 255, 0xff}}}
	p := make(color.Palette, 768)
	linGrad(cpts, p)
	for i := 0; i < 255; i++ {
		c := p[i].(color.RGBA)
		if c.R != uint8(i) || c.G != 0 || c.B != 0 {
			t.Fatalf("%d:%v", i, c)
		}
	}
	for i := 0; i < 255; i++ {
		c := p[i+256].(color.RGBA)
		if c.R != 0 || c.G != uint8(i) || c.B != 0 {
			t.Fatalf("%d:%v", i+256, c)
		}
	}
	for i := 0; i < 255; i++ {
		c := p[i+512].(color.RGBA)
		if c.R != 0 || c.G != 0 || c.B != uint8(i) {
			t.Fatalf("%d:%v", i+512, c)
		}
	}
}

func TestLinGrad2(t *testing.T) {
	cpts := []color.RGBA{
		{0, 0xff, 0, 0xff},
		{0xff, 0, 0xff, 0xff},
		{0, 0xff, 0, 0xff}}
	p := linGrad2(cpts, 511)
	if len(p) != 511 {
		t.Fatalf("len(p) = %d", len(p))
	}
	for i := 0; i < 256; i++ {
		c := p[i].(color.RGBA)
		if c.R != uint8(i) ||
			c.G != uint8(255-i) ||
			c.B != uint8(i) {
			t.Fatalf("%d:%v", i, c)
		}
	}
	for i := 0; i < 255; i++ {
		c := p[i+256].(color.RGBA)
		if c.R != uint8(254-i) ||
			c.G != uint8(i+1) ||
			c.B != uint8(254-i) {
			t.Fatalf("%d:%v", i+256, c)
		}
	}
}

func TestGrayPal(t *testing.T) {
	p := grayPal(2, 0xff, false)
	c := p[0].(color.RGBA)
	if c.R != 0 || c.G != 0 || c.B != 0 || c.A != 0xff {
		t.Fatalf("0:%v", c)
	}
	c = p[1].(color.RGBA)
	if c.R != 0xff || c.G != 0xff || c.B != 0xff || c.A != 0xff {
		t.Fatalf("0:%v", c)
	}
	p = grayPal(2, 0xff, true)
	c = p[0].(color.RGBA)
	if c.R != 0xff || c.G != 0xff || c.B != 0xff || c.A != 0xff {
		t.Fatalf("0:%v", c)
	}
	c = p[1].(color.RGBA)
	if c.R != 0 || c.G != 0 || c.B != 0 || c.A != 0xff {
		t.Fatalf("0:%v", c)
	}

}
