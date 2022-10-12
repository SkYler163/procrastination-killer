// +build tools

package tools

// tool dependencies
import (
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

//go:generate go build -v -o=../bin/golangci-lint github.com/golangci/golangci-lint/cmd/golangci-lint
