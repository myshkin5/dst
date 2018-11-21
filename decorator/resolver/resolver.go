package resolver

import (
	"errors"
	"go/ast"
	"go/types"
	"strings"
)

/*

// Consider this file... B and C could be local identifiers from a different file in this package,
// or from the imported package "a". If only one is from "a" and it is removed, we should remove the
// import when we restore the AST. Thus the node resolver interface needs to be able to resolve the
// package using the full info from go/types.

package main

import (
	. "a"
)

func main() {
	B()
	C()
}
*/

// PackageResolver resolves a package path to a package name.
type PackageResolver interface {
	ResolvePackage(path, dir string) (string, error)
}

// IdentResolver resolves an identifier to a package path. Returns an empty string if the node is
// not an identifier.
type IdentResolver interface {
	ResolveIdent(id *ast.Ident, info *types.Info, file *ast.File, dir string) (string, error)
}

var PackageNotFoundError = errors.New("package not found")

// Guess is a map of package path -> package name. Names are resolved from this map, and if a name
// doesn't exist in the map, the package name is guessed from the last part of the path (after the
// last slash).
type Guess map[string]string

func (r Guess) ResolvePackage(importPath, fromDir string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	if !strings.Contains(importPath, "/") {
		return importPath, nil
	}
	return importPath[strings.LastIndex(importPath, "/")+1:], nil
}

// Map is a map of package path -> package name. Names are resolved from this map, and if a name
// doesn't exist in the map, an error is returned. Note that Guess is not a NodeResolver, so can't
// properly resolve identifiers in dot import packages.
type Map map[string]string

func (r Map) ResolvePackage(importPath, fromDir string) (string, error) {
	if n, ok := r[importPath]; ok {
		return n, nil
	}
	return "", PackageNotFoundError
}
