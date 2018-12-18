bdisplay
========

# Development Setup

## Install Goimports

This is used for code formatting.

```
$ go get -u golang.org/x/tools/cmd/goimports
```

## Install Golint

This is used for code linting.

```
$ go get -u golang.org/x/lint/golint
```

## Set Up the Git Pre-Commit Hook

This pre-commit hook makes sure all code is properly formatted and linted before it is committed.

```
$ cp githooks/pre-commit .git/hooks/pre-commit
```

Install `goimports`, `golint` and `go vet` support in your editor of choice as needed.
