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
	"io/fs"
	"log"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"github.com/gx-org/ccgx/internal/cmd/debug"
	"github.com/gx-org/ccgx/internal/exec"
	"github.com/gx-org/ccgx/internal/gotc"
	gxmodule "github.com/gx-org/ccgx/internal/module"
	"golang.org/x/mod/module"
)

type gxFiles struct {
	mod  *gxmodule.Module
	list map[string]bool
}

func (fls *gxFiles) visit(path string, dir fs.DirEntry, err error) error {
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
	fls.list[gxPath] = true
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
	return gxCommand(mod, "github.com/gx-org/gx/golang/binder/genbind",
		"--language=cc",
		"--gx_package="+path,
		"--target_folder="+target,
	)
}

// Packages returns the list of GX packages in the current module.
func Packages(mod *gxmodule.Module) ([]string, error) {
	files := gxFiles{mod: mod, list: make(map[string]bool)}
	if err := filepath.WalkDir(mod.Path(), files.visit); err != nil {
		return nil, err
	}
	pkgs := slices.Collect(maps.Keys(files.list))
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

func installLinkToModule(cache *gotc.Cache, targetPath string, dep *module.Version) error {
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

// LinkAllDeps creates links to dependencies.
// Returns the path where the links where created.
func LinkAllDeps(mod *gxmodule.Module, cache *gotc.Cache) (string, error) {
	depsPath := filepath.Join(mod.Path(), "gxdeps")
	if err := os.MkdirAll(depsPath, 0755); err != nil {
		return "", err
	}
	for _, dep := range mod.Deps() {
		if err := installLinkToModule(cache, depsPath, dep); err != nil {
			return "", err
		}
	}
	return depsPath, nil
}
