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
	"image"
	"image/png"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	// Image size (in pixels)
	minSx = 320
	maxSx = 5120
	dflSx = 640
	minSy = 240
	maxSy = 4096
	dflSy = 480
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
)

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

type params struct {
	Sx, Sy         int
	Iter           int
	X0, Y0, X1, Y1 float64
}

func (p *params) URL() template.URL {
	s := fmt.Sprintf(
		"sx=%d&sy=%d&iter=%d&x0=%g&y0=%g&x1=%g&y1=%g",
		p.Sx, p.Sy, p.Iter,
		p.X0, p.Y0, p.X1, p.Y1)
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
	return p
}

func mandelHandler(w http.ResponseWriter, r *http.Request) {
	p := getParams(r)
	img := image.NewGray(image.Rect(0, 0, p.Sx, p.Sy))
	drawMandel(img,
		complex(p.X0, p.Y0), complex(p.X1, p.Y1), p.Iter)
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
