# ghglob

[![Go status](https://github.com/reviewdog/reviewdog/workflows/Go/badge.svg)](https://github.com/reviewdog/reviewdog/actions)

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
go get github.com/haya14busa/ghglob/cmd/ghglob
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
```

## Packages

| Package | GoDoc |
| ------- | ----- |
| ghglob | [![GoDoc - ghglob](https://godoc.org/github.com/haya14busa/ghglob?status.svg)](https://godoc.org/github.com/haya14busa/ghglob) |
| ghmatcher | [![GoDoc - ghmatcher](https://godoc.org/github.com/haya14busa/ghglob/ghmatcher?status.svg)](https://godoc.org/github.com/haya14busa/ghglob/ghmatcher) |
