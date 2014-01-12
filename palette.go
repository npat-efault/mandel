// Functions for generating and displaying gradient palettes

package main

import (
	"image/color"
)

// colPt specifies a "color point" in an interpolated gradient palette
type colPt struct {
	// Index of the point in the palette.
	Idx int
	// Color of the point.
	Col color.RGBA
}

// linterp performs linear interpolation between the first (p[0]) and
// the last (pal[len(p) - 1]) colors in palette "pal". It fills the
// remaining (pal[1:len(p) - 1]) palette slots with the intepolated
// colors. Colors in slots pal[0] and pal[len(p) - 1] MUST be of type
// color.RGBA. The remaining slots will also be filled with colors of
// type color.RGBA
func linterp(pal color.Palette) {
	n := len(pal)
	if n <= 2 {
		return
	}
	s, e := pal[0].(color.RGBA), pal[n-1].(color.RGBA)
	dr := int(e.R) - int(s.R)
	dg := int(e.G) - int(s.G)
	db := int(e.B) - int(s.B)
	da := int(e.A) - int(s.A)
	for i := 1; i < n-1; i++ {
		c := color.RGBA{}
		c.R = uint8(int(s.R) + i*dr/(n-1))
		c.G = uint8(int(s.G) + i*dg/(n-1))
		c.B = uint8(int(s.B) + i*db/(n-1))
		c.A = uint8(int(s.A) + i*da/(n-1))
		pal[i] = c
	}
}

// linGrad fills palette "pal" with a linearly interpolated gradient
// passing though the points in "pts". "pts" specifies colors at
// specific palette indexes; it must be sorted by index, and all
// indexes must be < len(pal). linGrad does not extrapolate: If pts[0]
// does not specify a color for palette index 0, then the first
// palette slots (pal[0:pts[0].Idx]) will not be filled. The same is
// true for the end of the palette. All palette colors will be of type
// color.RGBA.
//
//   []colPt
//    +--0, {0, 0, 0, 0xff}
//    | 255, {0xff, 0, 0, 0xff}--+
//    | 510, {0, 0, 0, 0xff} ----+---------------------+
//    |                          |                     |
//    |                          |                     |
//  +-V-+---+              +---+-V-+---+         +---+-V-+
//  | 0 | 1       ....      FE | FF| FE    ....    1 | 0 | P(Red)
//  +---+---+              +---+---+---+         +---+---+
//    0   1                 254 255 256           509 510
//
func linGrad(pts []colPt, pal color.Palette) {
	// TODO(npat): Sort pts? Sanity check pts?
	n := len(pts)
	pal[pts[0].Idx] = pts[0].Col
	for i := 0; i < n-1; i++ {
		pal[pts[i+1].Idx] = pts[i+1].Col
		linterp(pal[pts[i].Idx : pts[i+1].Idx+1])
	}
}

// linGrad2 generates a linearly interpolated gradient palette of the
// requested size ("size") by distributing the given colors ("pts")
// across the palette slots and interpolating between them. The first
// given color (pts[0]) will be placed at the start of the palette (at
// slot 0), the last given color (pts[len(pts)-1]) will be placed at
// the end of the palette (at slot size-1) and the remaining given
// colors will be spread across the palette at equal distances. The
// remailing (size - len(pts)) palette slots will be filled by
// interpolating between the given colors. At least 2, and no more
// than "size", colors must be given. All colors in the palette will
// be of type color.RGBA. Returns the generated palette.
func linGrad2(pts []color.RGBA, size int) color.Palette {
	n := len(pts)
	if n < 2 || n > size {
		return nil
	}
	pal := make(color.Palette, size)
	pal[0] = pts[0]
	p := 0
	for i := 1; i < n-1; i++ {
		c := i * size / (n - 1)
		pal[c] = pts[i]
		linterp(pal[p : c+1])
		p = c
	}
	pal[size-1] = pts[n-1]
	linterp(pal[p:])
	return pal
}

// grayPal generates a linearly interpolated gayscale palette of the
// given size. If "reverse" is false, the first color in the generated
// palette will be color.RGBA{0, 0, 0, 0} and the last color.RGBA{a,
// a, a, a}. If "reverse" is true, the colors are reversed: The first
// color in the palette will be color.RGBA{a, a, a, a} and the last
// color.RGBA{0, 0, 0, 0}. All colors in the palette will be of type
// color.RGBA. Returns the generated palette.
func grayPal(size int, a uint8, reverse bool) color.Palette {
	s := color.RGBA{0, 0, 0, a}
	e := color.RGBA{a, a, a, a}
	if reverse {
		s, e = e, s
	}
	return linGrad2([]color.RGBA{s, e}, size)
}

// pal256Gray is a linearly interpolated grayscale palette of 256
// colors. The first color is {0, 0, 0, 0xff} and the last {0xff,
// 0xff, 0xff, 0xff}. All colors are of type color.RGBA
var pal256Gray = grayPal(256, 0xff, false)

// pal256GrayR is a linearly interpolated grayscale palette of 256
// colors. The first color is {0xff, 0xff, 0xff, 0xff} and the last
// {0, 0, 0, 0xff}. All colors are of type color.RGBA
var pal256GrayR = grayPal(256, 0xff, true)

var pal256Gold1 = make(color.Palette, 256)
var pal256Gold2 = make(color.Palette, 256)
var pal256Blue1 = make(color.Palette, 256)
var pal256Blue2 = make(color.Palette, 256)
var pal256BRG = make(color.Palette, 256)

func init() {
	linGrad([]colPt{
		{0, color.RGBA{0x00, 0x00, 0x00, 0xff}},
		{220, color.RGBA{0x77, 0x55, 0x00, 0xff}},
		{245, color.RGBA{0xff, 0xff, 0x00, 0xff}},
		{255, color.RGBA{0xff, 0xff, 0xff, 0xff}}},
		pal256Gold1)

	linGrad([]colPt{
		{0, color.RGBA{0x00, 0x00, 0x00, 0xff}},
		{75, color.RGBA{0x77, 0x22, 0x00, 0xff}},
		{100, color.RGBA{0xff, 0xff, 0x00, 0xff}},
		{125, color.RGBA{0xff, 0xff, 0xff, 0xff}},
		{150, color.RGBA{0x77, 0x22, 0x00, 0xff}},
		{200, color.RGBA{0x00, 0x00, 0x00, 0xff}},
		{225, color.RGBA{0x77, 0x22, 0x00, 0xff}},
		{240, color.RGBA{0xff, 0xff, 0x00, 0xff}},
		{255, color.RGBA{0xff, 0xff, 0xff, 0xff}}},
		pal256Gold2)

	linGrad([]colPt{
		{0, color.RGBA{0x00, 0x00, 0x00, 0xff}},
		{220, color.RGBA{0x00, 0x00, 0x55, 0xff}},
		{245, color.RGBA{0x44, 0x44, 0xff, 0xff}},
		{255, color.RGBA{0xff, 0xff, 0xff, 0xff}}},
		pal256Blue1)

	linGrad([]colPt{
		{0, color.RGBA{0x00, 0x00, 0x00, 0xff}},
		{50, color.RGBA{0x22, 0x22, 0x55, 0xff}},
		{100, color.RGBA{0x10, 0x10, 0x55, 0xff}},
		{150, color.RGBA{0xff, 0xff, 0xff, 0xff}},
		{175, color.RGBA{0x55, 0x55, 0xff, 0xff}},
		{200, color.RGBA{0x22, 0x22, 0x77, 0xff}},
		{225, color.RGBA{0x00, 0x00, 0x20, 0xff}},
		{240, color.RGBA{0x22, 0x22, 0x77, 0xff}},
		{255, color.RGBA{0xff, 0xff, 0xff, 0xff}}},
		pal256Blue2)
}
