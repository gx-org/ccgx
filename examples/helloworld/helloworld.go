// Package helloworld encapsulates GX source files
// into a Go package.
//
// Automatically generated from google3/third_party/gxlang/gx/golang/packager/package.go.
//
// DO NOT EDIT
package helloworld

import (
	"embed"

	"github.com/gx-org/gx/build/builder"
	"github.com/gx-org/gx/build/importers/embedpkg"

	_ "github.com/gx-org/xlapjrt/gx"
)

//go:embed helloworld.gx
var srcs embed.FS

var inputFiles = []string{
	"helloworld.gx",
}

func init() {
	embedpkg.RegisterPackage("/helloworld", Build)
}

var _ embedpkg.BuildFunc = Build

// Build GX package.
func Build(bld *builder.Builder) (builder.Package, error) {
	return bld.BuildFiles("", "helloworld", srcs, inputFiles)
}

