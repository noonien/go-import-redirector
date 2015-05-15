go-import-redirector
====================
[![Docker Build Status](http://hubstatus.container42.com/noonien/go-import-redirector)](https://registry.hub.docker.com/u/noonien/go-import-redirector)
[![License: MIT](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/noonien/go-import-redirector/blob/master/LICENSE)

HTTP server that provides the [necessary metadata](https://golang.org/cmd/go/#hdr-Remote_import_paths)
for `go get` to work on custom domains.

Usage
-----

    go-import-redirector [-addr address] [-vcs vcs] [-parts parts] <import> <repo>

This starts a go-import-redirector instance that listens on the specified
address (defaults to ':http') and responds to requests to URLs under the
specified import path with a meta tag that specifies the correct path to a
repository that `go get` can use.

When invoked as:

    go-import-redirector / http://github.com/noonien

the response to <host>/kube-http-proxy will include the following meta tag:

    <meta name="go-import" content="<host>/kube-http-proxy git https://github.com/noonien/kube-http-proxy">

If <repo> contains a wildcard character ('*'), it is replaced with the root repository path.
As an example, when invoked as:

    go-import-redirector -parts 2 /local ssh://git@git.mux.ro/*.git

the response to <host>/local/internal/my-project will include the following
meta tag:

     <meta name="go-import" content="<host>/local/internal/my-project git ssh://git@git.mux.ro/internal/my-project.git">


Valid options:
  - addr - address on which the serve should listen (defaults to ':http').
  - vcs - version constrol system to use, git hg and svn are supported (defaults to 'git').
  - parts - how many URL parts after the import root represent a repository root. (defaults to 1).


Docker
------

A docker container is also provided at `noonien/go-import-redirector` and can
be run like so:

    docker run -p 80:80 noonien/go-import-redirector [-addr address] [-vcs vcs] [-parts parts] <import> <repo>
