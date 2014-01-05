mandel
======

Serve images of the Mandelbrot set over HTTP.

"mandel" is a simple web application, written in Go, that renders
images of the all-familiar Mandelbrot set and serves them over
HTTP. It allows zooming-in on any area of the set.

## Install

In order to build "mandel" you need to have package
"github.com/npat-efault/bundle" and command
"github.com/npat-efault/bundle/mkbundle" installed.

First install the prerequisites:

```
  $ mkdir -p $GOPATH/github.com/npat-efault/
  $ cd $GOPATH/github.com/npat-efault
  $ git clone https://github.com/npat-efault/bundle
  $ cd ./bundle
  $ ./all.sh install
```

Then install "mandel":

```
  $ cd $GOPATH/github.com/npat-efault
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

