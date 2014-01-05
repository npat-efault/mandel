package main

import (
	"image"
	"image/color"
	"math/cmplx"
)

func drawMandel(img *image.Gray, c0, c1 complex128, iter int) {
	var dx, dy, x, y float64
	var p0, p1 image.Point
	var px, py int
	var c, z complex128

	// Deltas for stepping on the complex plane
	dx = (real(c1) - real(c0)) / float64(img.Bounds().Dx())
	dy = (imag(c1) - imag(c0)) / float64(img.Bounds().Dy())
	// Image bounds
	p0 = img.Bounds().Min
	p1 = img.Bounds().Max
	// x, y are on the complex plane (world coordinates)
	// px, py are on the image (viewport coordinates)
	y = imag(c0)
	for py = p0.Y; py < p1.Y; py++ {
		x = real(c0)
		for px = p0.X; px < p1.X; px++ {
			c = complex(x, y)
			z = complex(0, 0)
			var i = 0
			for i = 0; i < iter; i++ {
				z = z*z + c
				if cmplx.Abs(z) > 2 {
					break
				}
			}
			img.SetGray(px, py,
				color.Gray{uint8((i * 256) / iter)})
			x += dx
		}
		y += dy
	}
}
