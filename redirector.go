// Copyright 2015 George Jiglau <george@mux.ro>. All rights reserved.
//
// This file is part of go-import-redirector.
//
// Use of this source code is governed by the MIT license that can be found
// in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

var (
	addr       = flag.String("addr", ":http", "address to serve on")
	vcs        = flag.String("vcs", "git", "version control system")
	repoParts  = flag.Int("parts", 1, "how many parts of the url represents the repo root")
	importPath string
	repoPath   string
)

func usage() {
	fmt.Fprintln(os.Stderr, "usage: go-import-redirector [options] <import> <repo>")
	fmt.Fprintln(os.Stderr, "options:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "examples")
	fmt.Fprintln(os.Stderr, "  go-import-redirector / http://github.com/noonien")
	fmt.Fprintln(os.Stderr, "  go-import-redirector -parts 2 /local ssh://git@git.mux.ro/*.git")
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
	}
	importPath = strings.TrimSuffix(flag.Arg(0), "/") + "/"
	repoPath = flag.Arg(1)
	if !strings.HasPrefix(importPath, "/") {
		fmt.Fprintln(os.Stderr, "import path has to start with a `/`")
		os.Exit(1)
	}
	if !strings.Contains(repoPath, "://") {
		log.Fatal("repo path must be a full URL, no schema defined")
	}

	http.HandleFunc(importPath, redirect)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var tmpl = template.Must(template.New("main").Parse(`<!DOCTYPE html>
<html>
<head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
  <meta name="go-import" content="{{.ImportRoot}} {{.VCS}} {{.RepoURL}}">
</head>
<body>Nothing to see here</body>
</html>
`))

type data struct {
	ImportRoot string
	VCS        string
	RepoURL    string
}

func redirect(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimSuffix(req.URL.Path, "/")
	if !strings.HasPrefix(path, importPath) {
		http.NotFound(w, req)
		return
	}
	path = strings.TrimPrefix(path, importPath)

	var repoLen, part int
	for part = 1; part <= *repoParts; part++ {
		n := strings.IndexRune(path[repoLen:], '/')
		if n == -1 {
			repoLen = len(path) + 1
			break
		}

		repoLen += n + 1
	}

	if part < *repoParts {
		http.NotFound(w, req)
		return
	}
	repoRoot := path[:repoLen-1]
	importRoot := req.Host + importPath + repoRoot

	var repoURL string
	if strings.Contains(repoPath, "*") {
		repoURL = strings.Replace(repoPath, "*", repoRoot, 1)
	} else {
		repoURL = strings.TrimSuffix(repoPath, "/") + "/" + repoRoot
	}

	d := &data{
		ImportRoot: importRoot,
		VCS:        *vcs,
		RepoURL:    repoURL,
	}
	var buf bytes.Buffer
	err := tmpl.Execute(&buf, d)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(buf.Bytes())

	log.Printf("Redirecting %s -> %s", importRoot, repoURL)
}
