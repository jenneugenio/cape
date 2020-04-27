package mage

import (
	"go/parser"
	"go/token"
	"strconv"
)

// FindImports returns a list of all imported packages inside of a given file
func FindImports(file string) ([]string, error) {
	out := []string{}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ImportsOnly)
	if err != nil {
		return out, err
	}

	for _, pkg := range f.Imports {
		path, err := strconv.Unquote(pkg.Path.Value)
		if err != nil {
			return out, err
		}

		out = append(out, path)
	}

	return out, nil
}
