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
	gxmodule "github.com/gx-org/gx/build/module"
	"github.com/gx-org/gx/golang/binder/ccbindings"
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
func Pack(mod *gxmodule.Module, path string) error {
	return gxCommand(mod, "github.com/gx-org/gx/golang/packager", "--gx_package_module="+path)
}

// Bind a GX package by generating C++ header files to a given target.
func Bind(mod *gxmodule.Module, path, target string) error {
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
	for _, binder := range bnd.Files() {
		bindingPath := binder.BuildFilePath(target, pkg.IR())
		if err := os.MkdirAll(filepath.Dir(bindingPath), 0755); err != nil {
			return fmt.Errorf("cannot create target folder: %v", err)
		}
		f, err := os.Create(bindingPath)
		if err != nil {
			return fmt.Errorf("cannot create target file: %v", err)
		}
		defer f.Close()
		if err := binder.WriteBindings(f); err != nil {
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
	for _, pkg := range pkgs {
		if err := Pack(mod, pkg); err != nil {
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
			"github.com/gx-org/gx/golang/binder/cgx",
		},
	}
	if err := filepath.WalkDir(mod.Root(), files.collectGXImports); err != nil {
		return err
	}
	src, err := writeGoSource(path, basename, &files)
	if err != nil {
		return err
	}
	filePath := filepath.Join(path, basename+".a")
	return gotc.BuildArchive(src, filePath)
}
