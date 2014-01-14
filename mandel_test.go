package main

import "testing"

func TestMandel(t *testing.T) {
	m, err := newMandelImg(160, 120, pal256Gray,
		complex(-2.0, -1.2), complex(1.0, 1.2), 16, 2)
	if err != nil {
		t.Fatal(err)
	}
	if m.Bounds().Dx() != 160 || m.Bounds().Dy() != 120 {
		t.Fatalf("Bad bounds: %v", m.Bounds())
	}
	if len(m.histo) != 16+1 {
		t.Fatalf("Bad histo len: %v", len(m.histo))
	}
	if len(m.cnhisto) != 16 {
		t.Fatalf("Bad cnhisto len: %v", len(m.cnhisto))
	}
	l := len(m.Palette)
	for i, v := range m.cnhisto {
		t.Logf("%02d: %d, %f, %d", i, m.histo[i],
			v, int(float64(l)*v))
	}
}
