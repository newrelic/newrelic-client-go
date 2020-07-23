// +build tools

package main

import (
	// build/test.mk
	_ "github.com/stretchr/testify/assert"
	_ "gotest.tools/gotestsum"

	// build/lint.mk
	_ "github.com/client9/misspell/cmd/misspell"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/llorllale/go-gitlint/cmd/go-gitlint"
	_ "github.com/psampaz/go-mod-outdated"
	_ "golang.org/x/tools/cmd/goimports"

	// build/document.mk
	_ "github.com/git-chglog/git-chglog/cmd/git-chglog"
	_ "golang.org/x/tools/cmd/godoc"

	// build/release.mk
	_ "github.com/goreleaser/goreleaser"
)
