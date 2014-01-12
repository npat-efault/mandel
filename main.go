// Command mandel starts a simple web application that renders images
// of the Mandelbrot Set and serves them over HTTP.
//
// Usage is:
//
//     mandel <laddr>
//
// Where "<laddr>" is the TCP local network address to listen for HTTP
// connections to. Example:
//
//     mandel :8080
//
package main

import (
	"fmt"
	"html/template"
	"image/color"
	"image/png"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"
)

const (
	// Image size (in pixels)
	minSx = 320
	maxSx = 10240
	dflSx = 640
	minSy = 256
	maxSy = 8192
	dflSy = 512
	// Number of iterations to perform for every pixel, when
	// rendering the image
	minIter = 16
	maxIter = 10240
	dflIter = 64
	// Function domain:
	// (Real: [minX .. maxX], Imag: [minY .. maxY])
	minX  = -2.0
	maxX  = 1.0
	dflX0 = minX
	dflX1 = maxX
	minY  = -1.2
	maxY  = 1.2
	dflY0 = minY
	dflY1 = maxY
	// Default palette
	dflPal = "Gray"
)

var palettes = map[string]color.Palette{
	"Blue 1":       pal256Blue1,
	"Blue 2":       pal256Blue2,
	"Gold 1":       pal256Gold1,
	"Gold 2":       pal256Gold2,
	"Gray":         pal256Gray,
	"Gray Reverse": pal256GrayR}

var templates *template.Template

func renderTmpl(w http.ResponseWriter, t string, d interface{}) {
	err := templates.ExecuteTemplate(w, t+".html", d)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
	}
}

func valInt(r *http.Request, p string, min, max, dfl int) int {
	s := r.FormValue(p)
	v, err := strconv.Atoi(s)
	if err != nil {
		v = dfl
	} else if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

func valFloat64(r *http.Request, p string,
	min, max, dfl float64) float64 {
	s := r.FormValue(p)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		v = dfl
	} else if v < min {
		v = min
	} else if v > max {
		v = max
	}
	return v
}

func valPalette(r *http.Request, p string,
	valid map[string]color.Palette, dfl string) string {
	s := r.FormValue(p)
	_, ok := valid[s]
	if !ok {
		s = dfl
	}
	return s
}

type params struct {
	Sx, Sy         int
	Iter           int
	X0, Y0, X1, Y1 float64
	Pal            string
	Palettes       map[string]color.Palette
}

func (p *params) URL() template.URL {
	s := fmt.Sprintf(
		"sx=%d&sy=%d&iter=%d&x0=%g&y0=%g&x1=%g&y1=%g&pal=%s",
		p.Sx, p.Sy, p.Iter,
		p.X0, p.Y0, p.X1, p.Y1,
		p.Pal)
	return template.URL(s)
}

func getParams(r *http.Request) *params {
	p := &params{}
	// Parse "sx" and "sy" (img size) parameters
	p.Sx = valInt(r, "sx", minSx, maxSx, dflSx)
	p.Sy = valInt(r, "sy", minSy, maxSy, dflSy)
	// Parse "iter" (# of iterations) parameter
	p.Iter = valInt(r, "iter", minIter, maxIter, dflIter)
	// Parse x0, x1, y0, y1 (coordinates) parameters
	p.X0 = valFloat64(r, "x0", minX, maxX, dflX0)
	p.X1 = valFloat64(r, "x1", minX, maxX, dflX1)
	p.Y0 = valFloat64(r, "y0", minY, maxY, dflY0)
	p.Y1 = valFloat64(r, "y1", minY, maxY, dflY1)
	// Parse pal (palette name) parameter
	p.Pal = valPalette(r, "pal", palettes, dflPal)
	p.Palettes = palettes
	return p
}

func mandelHandler(w http.ResponseWriter, r *http.Request) {
	// Parse params and calculate image
	p := getParams(r)
	img, _ := NewMandelImg(p.Sx, p.Sy, p.Palettes[p.Pal],
		complex(p.X0, p.Y0), complex(p.X1, p.Y1), p.Iter, 100.0)
	// Allow caching forever
	t := time.Now().Add(365 * 24 * time.Hour)
	w.Header().Set("Expires", t.Format(http.TimeFormat))
	// Encode and send image
	png.Encode(w, img)
}

func handler(w http.ResponseWriter, r *http.Request) {
	p := getParams(r)
	renderTmpl(w, "main", p)
}

func Usage(cmd string) {
	fmt.Fprintf(os.Stderr, "Usage is: %s <local addr>\n", cmd)
}

func main() {
	if len(os.Args) != 2 {
		Usage(path.Base(os.Args[0]))
		os.Exit(1)
	}
	templates = parseEntries(_bundleIdx, "templates/", ".html")
	http.Handle("/js/", serveEntries(_bundleIdx, "js/", "/js/"))
	http.Handle("/css/", serveEntries(_bundleIdx, "css/", "/css/"))
	http.HandleFunc("/mandel", mandelHandler)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(os.Args[1], nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
