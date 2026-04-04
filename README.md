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
`GRAB_HOME` or the config file). So, in this example, the repo
`github.com/jmhodges/grab` was stored in it with the domain (`github.com`) as
the top-level directory, and `jmhodges` and `grab` created as subdirectories
down the path.

Also, the input to `grab` doesn't have to contain a scheme (e.g. `https://`) to
work.

Install
-------

Grab can be installed by running `go install github.com/jmhodges/grab@latest`.

Configuration
-------------

Grab can be configured with the `GRAB_HOME` environment variable and further
configurations (including that setting) can be set with a TOML config file at
`grab/config.toml` inside your platforms configuration directory.

(On Windows, that's in `%AppData%`. On macOS, that's `$HOME/.config/` or
`$HOME/Library/Application Support/grab` if `.config/grab/config.toml` doesn't exist.
On Linux and BSD, that's `$XDG_CONFIG_HOME` or `$HOME/.config/` if that's not
set.)

The config settings in `config.toml` and `config.toml` itself are all optional.

```toml
home = "/home/user/src" # prefer the env var GRAB_HOME if shell-variables are needed
```

- **`home`**: The directory to download repos into. Equivalent to `GRAB_HOME`.

Environment variables override config file values:

- `GRAB_HOME` overrides `home`.
