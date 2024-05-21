//go:build tools
// +build tools

package tools

// tool dependencies
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

//go:generate go build -v -o=../bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
