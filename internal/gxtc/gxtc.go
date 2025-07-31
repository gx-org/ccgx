// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package gxtc provides functions to invoke the GX toolchain.
package gxtc

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/gx-org/ccgx/internal/cmd/debug"
	"github.com/gx-org/ccgx/internal/exec"
	"github.com/gx-org/ccgx/internal/gotc"
	"github.com/gx-org/gx/build/builder"
	"github.com/gx-org/gx/build/importers"
	"github.com/gx-org/gx/build/importers/localfs"
	"github.com/gx-org/gx/build/ir"
	gxmodule "github.com/gx-org/gx/build/module"
	"github.com/gx-org/gx/golang/binder/bindings"
	"github.com/gx-org/gx/golang/binder/ccbindings"
	"github.com/gx-org/gx/golang/packager/goembed"
	"github.com/gx-org/gx/golang/packager/pkginfo"
	"github.com/gx-org/gx/stdlib"
	gomodule "golang.org/x/mod/module"
)

type gxFiles struct {
	mod  *gxmodule.Module
	list []string
}

func (fls *gxFiles) collectGXDeps(path string, dir fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !strings.HasSuffix(path, ".gx") {
		return nil
	}
	gxPath, err := fls.mod.GXPathFromOS(path)
	if err != nil {
		return err
	}
	if gxPath == "" {
		return nil
	}
	fls.list = append(fls.list, gxPath)
	return nil
}

func (fls *gxFiles) collectGXImports(path string, dir fs.DirEntry, err error) error {
	if err != nil {
		return err
	}
	if !strings.HasSuffix(path, ".gx") {
		return nil
	}
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
	if err != nil {
		return err
	}
	for _, imp := range file.Imports {
		impPath, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			return fmt.Errorf("%s: import path %q is invalid: %v", path, imp.Path.Value, err)
		}
		fls.list = append(fls.list, impPath)
	}
	return nil
}

func gxCommand(mod *gxmodule.Module, gxcmd string, args ...string) error {
	version := mod.VersionOf("github.com/gx-org/gx")
	if version == "" {
		version = "latest"
	}
	osArgs := []string{"run", gxcmd + "@" + version}
	osArgs = append(osArgs, args...)
	if debug.Debug {
		cmdS := append([]string{"DEBUG", "go"}, osArgs...)
		log.Println(strings.Join(cmdS, " "))
	}
	cmd := exec.Command("go", osArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Pack a GX package.
func Pack(mod *gxmodule.Module, targetRoot string, pkgPath string) error {
	pkgInfo, err := pkginfo.Load(mod, pkgPath)
	if err != nil {
		return err
	}
	pkgPaths := strings.Split(pkgPath, "/")
	targetFolder := filepath.Join(targetRoot, filepath.Join(pkgPaths...))
	targetFile := filepath.Join(targetFolder, pkgInfo.GoPackageName()+"_gx.go")
	if err := os.MkdirAll(filepath.Dir(targetFile), os.ModePerm); err != nil {
		return err
	}
	w, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer w.Close()
	if err := goembed.Write(w, pkgInfo); err != nil {
		return err
	}
	for _, gxSrc := range pkgInfo.SourceFiles() {
		gxDst := filepath.Join(targetFolder, filepath.Base(gxSrc))
		if err := copy(gxSrc, gxDst); err != nil {
			return err
		}
	}
	return err
}

// BinderCallback is a function called after bindings have been generated for a package.
type BinderCallback func(target string, pkg *ir.Package, headerPath, ccPath string) error

const cmakeSource = `
cmake_minimum_required (VERSION 3.24)
project (%s_bindings)

include_directories (${CMAKE_CURRENT_LIST_DIR}/../../..)

add_library (%s_bindings STATIC ${CMAKE_CURRENT_LIST_DIR}/%s)
`

// WriteCMakeLists writes CMakeLists.txt for a given package.
func WriteCMakeLists(target string, pkg *ir.Package, headerPath, ccPath string) error {
	path := filepath.Join(target, "CMakeLists.txt")
	text := fmt.Sprintf(cmakeSource,
		pkg.Name.Name,
		pkg.Name.Name,
		filepath.Base(ccPath),
	)
	return os.WriteFile(path, []byte(text), 0755)
}

func writeBinderSourceFile(binder bindings.File, target string, pkg *ir.Package) (string, error) {
	bindingPath := binder.BuildFilePath(target, pkg)
	if err := os.MkdirAll(filepath.Dir(bindingPath), 0755); err != nil {
		return "", fmt.Errorf("cannot create target folder: %v", err)
	}
	f, err := os.Create(bindingPath)
	if err != nil {
		return "", fmt.Errorf("cannot create target file: %v", err)
	}
	defer f.Close()
	if err := binder.WriteBindings(f); err != nil {
		return "", err
	}
	return bindingPath, nil
}

// Bind a GX package by generating C++ header files to a given target.
func Bind(mod *gxmodule.Module, path, target string, fs ...BinderCallback) error {
	localImporter, err := localfs.NewWithModule(mod)
	if err != nil {
		return fmt.Errorf("cannot create local importer: %v", err)
	}
	bld := builder.New(importers.NewCacheLoader(
		stdlib.Importer(nil),
		localImporter,
	))
	pkg, err := bld.Build(path)
	if err != nil {
		return err
	}
	bnd, err := ccbindings.New(pkg.IR())
	if err != nil {
		return err
	}
	ccFiles := bnd.Files()
	headerPath, err := writeBinderSourceFile(ccFiles[0], target, pkg.IR())
	if err != nil {
		return err
	}
	ccPath, err := writeBinderSourceFile(ccFiles[1], target, pkg.IR())
	if err != nil {
		return err
	}
	for _, f := range fs {
		if err := f(target, pkg.IR(), headerPath, ccPath); err != nil {
			return err
		}
	}
	return nil
}

func unique(ss []string) []string {
	m := make(map[string]bool)
	for _, s := range ss {
		m[s] = true
	}
	return slices.Collect(maps.Keys(m))
}

// Packages returns the list of GX packages in the current module.
func Packages(mod *gxmodule.Module) ([]string, error) {
	files := gxFiles{mod: mod}
	if err := filepath.WalkDir(mod.Root(), files.collectGXDeps); err != nil {
		return nil, err
	}
	pkgs := unique(files.list)
	sort.Strings(pkgs)
	return pkgs, nil
}

// PackAll looks for all GX packages and generates a matching Go package to encapsulte the GX source code.
func PackAll(mod *gxmodule.Module) error {
	pkgs, err := Packages(mod)
	if err != nil {
		return err
	}
	packagerRoot, err := DepsPath(mod)
	if err != nil {
		return err
	}
	packagerRoot = filepath.Join(packagerRoot, "packager")
	for _, pkg := range pkgs {
		if err := Pack(mod, packagerRoot, pkg); err != nil {
			return err
		}
	}
	return nil
}

func installLinkToModule(cache *gotc.Cache, targetPath string, dep *gomodule.Version) error {
	gxModPath, err := cache.OSPath(dep)
	if err != nil {
		return err
	}
	targetLink := filepath.Join(targetPath, dep.Path)
	folder := filepath.Dir(targetLink)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}
	if _, err := os.Stat(targetLink); err == nil {
		if err := os.Remove(targetLink); err != nil {
			return err
		}
	}
	return os.Symlink(gxModPath, targetLink)
}

// DepsPath returns the path where dependencies are linked.
// It is created if it does not exist.
func DepsPath(mod *gxmodule.Module) (string, error) {
	depsPath := filepath.Join(mod.Root(), "gxdeps")
	if err := os.MkdirAll(depsPath, 0755); err != nil {
		return "", err
	}
	return depsPath, nil
}

// LinkAllDeps creates links to dependencies.
// Returns the path where the links where created.
func LinkAllDeps(mod *gxmodule.Module, cache *gotc.Cache) error {
	if err := gotc.ModTidy(); err != nil {
		return err
	}
	depsPath, err := DepsPath(mod)
	if err != nil {
		return err
	}
	for _, dep := range mod.Deps() {
		if err := installLinkToModule(cache, depsPath, dep); err != nil {
			return err
		}
	}
	return err
}

func writeGoSource(path, name string, files *gxFiles) (string, error) {
	deps := unique(files.list)
	var imports strings.Builder
	for _, dep := range deps {
		if !strings.HasPrefix(dep, "github.com") {
			continue
		}
		fmt.Fprintf(&imports, "import _ %s\n", strconv.Quote(dep))
	}
	cArchiveSource := fmt.Sprintf(`package main

%s

import _ "helloworld/gxdeps/packager/helloworld"

import "fmt"

import "C"

//export InitModuleHelloworld
func InitModuleHelloworld() {
	fmt.Println("Hello from here")
}

func main() {}
`, imports.String())
	srcFile := filepath.Join(path, name+".go")
	return srcFile, os.WriteFile(srcFile, []byte(cArchiveSource), 0644)
}

const basename string = "carchive"

// CompileCArchive creates a Go file with all the GX/Go dependencies and
// a main function. This file is then compiled into a static binary C library.
func CompileCArchive(mod *gxmodule.Module, path string) error {
	files := gxFiles{
		mod: mod,
		list: []string{
			"google3/third_party/golang/github_com/gomlx/gopjrt/v/v0/pjrt/pjrt",
			"github.com/gx-org/gx/golang/binder/cgx",
			"github.com/gx-org/xlapjrt/cgx",
		},
	}
	if err := filepath.WalkDir(mod.Root(), files.collectGXImports); err != nil {
		return err
	}
	src, err := writeGoSource(path, basename, &files)
	if err != nil {
		return err
	}
	if err := gotc.ModTidy(); err != nil {
		return err
	}
	filePath := filepath.Join(path, basename+".a")
	return gotc.BuildArchive(mod.Root(), src, filePath)
}
