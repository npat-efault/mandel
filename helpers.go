// Helper functions for doing stuff with bundle entries

package main

import (
	"github.com/npat-efault/bundle"
	"html/template"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// parseEntries parses (as HTML templates) all bundle entries with
// names starting with the given prefix and ending with the given
// extension (suffix). Each template is given the name of the
// respective entry without the prefix. Returns the first template
// parsed. All other templates are associated with it. Panics on
// error.
func parseEntries(idx bundle.Index,
	pref, ext string) *template.Template {
	var tmpl *template.Template

	for _, e := range idx.Dir(pref) {
		if !strings.HasSuffix(e.Name, ext) {
			continue
		}
		b, err := e.Decode(0)
		if err != nil {
			panic("parseEntries: " +
				e.Name + ": " + err.Error())
		}
		s := string(b)
		name := e.Name[len(pref):]

		var t *template.Template
		if tmpl == nil {
			tmpl = template.New(name)
		}
		if name == tmpl.Name() {
			t = tmpl
		} else {
			t = tmpl.New(name)
		}
		_, err = t.Parse(s)
		if err != nil {
			panic("parseEntries: " +
				e.Name + ": " + err.Error())
		}
	}
	if tmpl == nil {
		panic("parseEntries: No entries to parse")
	}
	return tmpl
}

// entryHandler implements http.Handler and serves bundle entries over
// HTTP. See doc of function serveEntries for details.
type entryHandler struct {
	idx     bundle.Index
	pref    string
	baseURL string
}

func (h *entryHandler) ServeHTTP(
	w http.ResponseWriter, r *http.Request) {
	/* Locate entry */
	p := h.pref + r.URL.Path[len(h.baseURL):]
	e := h.idx.Entry(p)
	if e == nil {
		http.NotFound(w, r)
		return
	}
	// Try to set Content-Type based on extension
	ctype := mime.TypeByExtension(filepath.Ext(e.Name))
	if ctype != "" {
		w.Header().Set("Content-Type", ctype)
	}
	// Open and send entry
	erd, err := e.Open(0)
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	io.Copy(w, erd)
}

// serveEntries returns a new entryHandler that serves bundle entries
// with names starting with the given prefix. A request with URL.Path
// "baseURL"foo will be replied by the data of the entry "pref"foo.
func serveEntries(idx bundle.Index,
	pref, baseURL string) *entryHandler {
	return &entryHandler{
		idx:     idx,
		pref:    pref,
		baseURL: baseURL}
}
