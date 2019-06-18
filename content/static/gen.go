// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

//go:generate go run makestatic.go

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"time"
	"unicode"
)

var files = []string{
	"analysis/call3.png",
	"analysis/call-eg.png",
	"analysis/callers1.png",
	"analysis/callers2.png",
	"analysis/chan1.png",
	"analysis/chan2a.png",
	"analysis/chan2b.png",
	"analysis/error1.png",
	"analysis/help.html",
	"analysis/ident-def.png",
	"analysis/ident-field.png",
	"analysis/ident-func.png",
	"analysis/ipcg-func.png",
	"analysis/ipcg-pkg.png",
	"analysis/typeinfo-pkg.png",
	"analysis/typeinfo-src.png",
	"callgraph.html",
	"codewalk.html",
	"codewalkdir.html",
	"dirlist.html",
	"doc/copyright.html",
	"doc/root.html",
	"error.html",
	"example.html",
	"godoc.html",
	"godocs.js",
	"images/cloud-download.svg",
	"images/go-logo-blue.svg",
	"images/google-logo.svg",
	"images/home-gopher.png",
	"images/footer-gopher.jpg",
	"images/minus.gif",
	"images/play-link.svg",
	"images/plus.gif",
	"images/treeview-black-line.gif",
	"images/treeview-black.gif",
	"images/treeview-default-line.gif",
	"images/treeview-default.gif",
	"images/treeview-gray-line.gif",
	"images/treeview-gray.gif",
	"implements.html",
	"jquery.js",
	"jquery.treeview.css",
	"jquery.treeview.edit.js",
	"jquery.treeview.js",
	"methodset.html",
	"opensearch.xml",
	"package.html",
	"packageroot.html",
	"play.js",
	"playground.js",
	"search.html",
	"searchcode.html",
	"searchdoc.html",
	"searchtxt.html",
	"style.css",
}

// Generate reads a set of files and returns a file buffer that declares
// a map of string constants containing contents of the input files.
func Generate() ([]byte, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n\n%v\n\npackage static\n\n", license, warning)
	fmt.Fprintf(buf, "var Files = map[string]string{\n")
	for _, fn := range files {
		b, err := ioutil.ReadFile(fn)
		if err != nil {
			return b, err
		}
		fmt.Fprintf(buf, "\t%q: ", fn)
		appendQuote(buf, b)
		fmt.Fprintf(buf, ",\n\n")
	}
	fmt.Fprintln(buf, "}")
	return format.Source(buf.Bytes())
}

// appendQuote is like strconv.AppendQuote, but we avoid the latter
// because it changes when Unicode evolves, breaking gen_test.go.
func appendQuote(out *bytes.Buffer, data []byte) {
	out.WriteByte('"')
	for _, b := range data {
		if b == '\\' || b == '"' {
			out.WriteByte('\\')
			out.WriteByte(b)
		} else if b <= unicode.MaxASCII && unicode.IsPrint(rune(b)) && !unicode.IsSpace(rune(b)) {
			out.WriteByte(b)
		} else {
			fmt.Fprintf(out, "\\x%02x", b)
		}
	}
	out.WriteByte('"')
}

const warning = `// Code generated by "makestatic"; DO NOT EDIT.`

var license = fmt.Sprintf(`// Copyright %d The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.`, time.Now().UTC().Year())
