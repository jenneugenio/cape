// +build tools

package main

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/jackc/tern"
	_ "github.com/magefile/mage"
	_ "github.com/markbates/pkger/cmd/pkger"
	_ "helm.sh/helm/v3/cmd/helm"
	_ "sigs.k8s.io/kind/cmd/kind"
)
