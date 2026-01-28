#!/usr/bin/env bash
# Install/update Go developer tools (latest) into .tools/bin

set -euo pipefail

#mkdir -p .tools/bin .tools/gopath

# LSP
go install golang.org/x/tools/gopls@latest
# golangci-lint LSP wrapper
go install github.com/nametake/golangci-lint-langserver@latest

# Formatting
go install golang.org/x/tools/cmd/goimports@latest

# Debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# Static analysis
go install honnef.co/go/tools/cmd/staticcheck@latest

# Protobuf / gRPC generators
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Shell formatting
go install mvdan.cc/sh/v3/cmd/shfmt@latest

# Find symbol information in Go source
go install github.com/rogpeppe/godef@latest

# Test runner (it is in the go tool)
# go install gotest.tools/gotestsum@latest

# Go enum generator (It's in the go tool)
# go install github.com/abice/go-enum
