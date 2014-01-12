mandel
======

Serve images of the Mandelbrot set over HTTP.

"mandel" is a simple web application, written in Go, that renders
images of the all-familiar [Mandelbrot
Set](http://en.wikipedia.org/wiki/Mandelbrot_set) and serves them over
HTTP. Features:

- Zoom-in on any area of the set.
- Select palette to use when rendering the set
- Select image size (WxH in pixels)
- Select maximum iteration count
- Color using the [histogram
  method](http://en.wikipedia.org/wiki/Mandelbrot_set#Histogram_coloring)
  that produces coloring independent of the maximum iteration count
  setting.
- Self-contained binary with no external support files. To install to
  another server, just copy the binary and run it.
  
## Example Renderings

[A few pictures](https://github.com/npat-efault/mandel/wiki/Example-Renderings)
rendered with "mandel".

## Install

In order to build the "mandel" command you need to have the package
"github.com/npat-efault/bundle" and the command
"github.com/npat-efault/bundle/mkbundle" installed.

First install the prerequisites:

```
  $ mkdir -p $GOPATH/src/github.com/npat-efault/
  $ cd $GOPATH/src/github.com/npat-efault
  $ git clone https://github.com/npat-efault/bundle
  $ cd ./bundle
  $ ./all.sh install
```

Then install "mandel":

```
  $ cd $GOPATH/src/github.com/npat-efault
  $ git clone https://github.com/npat-efault/mandel
  $ cd ./mandel
  $ ./all.sh install
```

## Run

Run like this:

```
  $ mandel :8080
```

Then direct your browser to http://localhost:8080/

