# ghglob

[![Go status](https://github.com/haya14busa/ghglob/workflows/Go/badge.svg)](https://github.com/haya14busa/ghglob/actions)
[![releases](https://img.shields.io/github/release/haya14busa/ghglob.svg)](https://github.com/haya14busa/ghglob/releases)
[![codecov](https://codecov.io/gh/haya14busa/ghglob/branch/master/graph/badge.svg)](https://codecov.io/gh/haya14busa/ghglob)

**ghglob** is glob, or more like pattern matcher based on GitHub Actions's
[Filter pattern spec](https://help.github.com/en/articles/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet).

Support multiple patterns including negation (`!<pattern>`).

## Spec

> - `*` matches zero or more characters, but does not match the / character
> - `**` matches zero or more of any character
> - `?` matches zero or one of the proceeding character
> - `+` matches one or more of the proceeding character
> - `[]` matches any character listed, or included in ranges. Ranges can only include a-zA-Z0-9. e.g [123abc] or [0-9a-f]
> - `!` at the start of a pattern makes it negate previous positive patterns. It has no special meaning if not the first character
>
> -- https://help.github.com/en/articles/workflow-syntax-for-github-actions#filter-pattern-cheat-sheet


## CLI

### Installation

```shell
# Install latest version. (Install it into ./bin/ by default).
$ curl -sfL https://raw.githubusercontent.com/haya14busa/ghglob/master/install.sh| sh -s

# Specify installation directory ($(go env GOPATH)/bin/) and version.
$ curl -sfL https://raw.githubusercontent.com/haya14busa/ghglob/master/install.sh| sh -s -- -b $(go env GOPATH)/bin [vX.Y.Z]

# In alpine linux (as it does not come with curl by default)
$ wget -O - -q https://raw.githubusercontent.com/haya14busa/ghglob/master/install.sh| sh -s [vX.Y.Z]

$ go get github.com/haya14busa/ghglob/cmd/ghglob

# homebrew / linuxbrew
$ brew install haya14busa/tap/ghglob
$ brew upgrade haya14busa/tap/ghglob

# Go
$ go get github.com/haya14busa/ghglob/cmd/ghglob
```

### Example usages

```
$ ghglob **/*.go'
cmd/ghglob/main.go 99
ghmatcher/ghmatcher.go
ghmatcher/ghmatcher_test.go
ghglob_test.go
ghglob.go
_testdir/main.go
cmd/ghglob/main.go

# Support negation pattern.
$ git ls-files | ghglob '**.go' '!**_test.go'
_testdir/main.go
cmd/ghglob/main.go
ghglob.go
ghmatcher/ghmatcher.go

# List .js or .jsx files.
$ ghglob "**.jsx?"

# As CLI arguments. e.g. https://github.com/koalaman/shellcheck
$ shellcheck $(ghglob '**/*.sh')

# gofmt except test data.
$ gofmt -d -s $(ghglob '**.go' '!_testdir/**')

# Run golint against only changed files.
$ git diff --name-only master | ghglob '**.go' | xargs -i golint {}
```

## Packages

| Package | GoDoc |
| ------- | ----- |
| ghglob | [![GoDoc - ghglob](https://godoc.org/github.com/haya14busa/ghglob?status.svg)](https://godoc.org/github.com/haya14busa/ghglob) |
| ghmatcher | [![GoDoc - ghmatcher](https://godoc.org/github.com/haya14busa/ghglob/ghmatcher?status.svg)](https://godoc.org/github.com/haya14busa/ghglob/ghmatcher) |
