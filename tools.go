// +build tools

package main

import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/jackc/tern"
	_ "github.com/magefile/mage"
)
