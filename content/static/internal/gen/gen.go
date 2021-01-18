// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package gen is used by content/static/makestatic.go
// to generate content/static/static.go.
//
// This is a separate package so that it can be tested without
// build constraints. cmd/golangorg and other binaries should not
// depend on it.
package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"strings"
	"unicode"

	"golang.org/x/website/internal/markdown"
)

var files = []string{
	"analysis/call-eg.png",
	"analysis/call3.png",
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
	"codewalk.html",
	"codewalkdir.html",
	"dirlist.html",
	"doc/code.html",
	"doc/conduct.html",
	"doc/contrib.html",
	"doc/copyright.html",
	"doc/devel/pre_go1.html",
	"doc/devel/release.html",
	"doc/devel/weekly.html",
	"doc/docs.html",
	"doc/download.js",
	"doc/gopath_code.html",
	"doc/hats.js",
	"doc/install.html",
	"doc/install-source.html",
	"doc/manage-install.html",
	"doc/modules/images/multiple-modules.png",
	"doc/modules/images/single-module.png",
	"doc/modules/images/source-hierarchy.png",
	"doc/modules/images/v2-branch-module.png",
	"doc/modules/images/v2-module.png",
	"doc/modules/images/version-number.png",
	"doc/mvs/buildlist.svg",
	"doc/mvs/downgrade.svg",
	"doc/mvs/exclude.svg",
	"doc/mvs/get-downgrade.svg",
	"doc/mvs/get-upgrade.svg",
	"doc/mvs/replace.svg",
	"doc/mvs/upgrade.svg",
	"doc/root.html",
	"doc/security.html",
	"doc/tutorial/add-a-test.html",
	"doc/tutorial/call-module-code.html",
	"doc/tutorial/compile-install.html",
	"doc/tutorial/create-module.html",
	"doc/tutorial/getting-started.html",
	"doc/tutorial/greetings-multiple-people.html",
	"doc/tutorial/handle-errors.html",
	"doc/tutorial/images/function-syntax.png",
	"doc/tutorial/index.html",
	"doc/tutorial/random-greeting.html",
	"error.html",
	"example.html",
	"godoc.html",
	"godocs.js",
	"images/cloud-download.svg",
	"images/footer-gopher.jpg",
	"images/go-logo-blue.svg",
	"images/home-gopher.png",
	"images/minus.gif",
	"images/play-link.svg",
	"images/plus.gif",
	"jquery.js",
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

var markdownFiles = []string{
	"doc/mod.md",
	"doc/modules/developing.md",
	"doc/modules/gomod-ref.md",
	"doc/modules/major-version.md",
	"doc/modules/managing-dependencies.md",
	"doc/modules/managing-source.md",
	"doc/modules/publishing.md",
	"doc/modules/release-workflow.md",
	"doc/modules/version-numbers.md",
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
			return nil, err
		}
		fmt.Fprintf(buf, "\t%q: ", fn)
		appendQuote(buf, b)
		fmt.Fprintf(buf, ",\n\n")
	}

	for _, fn := range markdownFiles {
		src, err := ioutil.ReadFile(fn)
		if err != nil {
			return nil, err
		}
		gen, err := markdown.Render(src)
		if err != nil {
			return nil, fmt.Errorf("%s: %v", fn, err)
		}
		htmlName := strings.TrimSuffix(fn, ".md") + ".html"
		fmt.Fprintf(buf, "\t%q: ", htmlName)
		appendQuote(buf, gen)
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

const license = `// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.`
