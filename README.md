Grab
====

Grab is a tool that downloads source code repositories into a convenient
directory layout created from the repo's URL's domain and path. It supports Git,
Mercurial (hg), Subversion, and Bazaar repositories.


```bash
$ grab github.com/jmhodges/grab # https://github.com/jmhodges/grab also works.

$ ls ~/src/github.com/jmhodges/grab
LICENSE   README.md go.mod    go.sum    grab      main.go
```

By default, grab downloads into `$HOME/src` (overridable with the env var
`GRAB_HOME`). The repo `github.com/jmhodges/grab` was stored in it with the
domain (`github.com`) as the top-level directory, and `jmhodges` and `grab`
created as subdirectories down the path.

Also, the input to `grab` doesn't have to contain a scheme (e.g. `https://`) to
work.

Install
-------

Grab can be installed by running `go install github.com/jmhodges/grab`.
