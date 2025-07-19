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

// Package bind provides the Cobra bind command
package bind

import (
	"ccgx/internal/gotc"
	"ccgx/internal/gxtc"
	gxmodule "ccgx/internal/module"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"golang.org/x/mod/module"
)

// Cmd is the implementation of the mod command.
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bind",
		Short: "Generate C++ header files",
		RunE:  cBind,
	}
}

func installLinkToModule(mod *gxmodule.Module, targetPath string, dep *module.Version) error {
	gxModPath, err := gotc.ModuleOSPath(dep)
	if err != nil {
		return err
	}
	targetLink := filepath.Join(targetPath, dep.Path)
	folder := filepath.Dir(targetLink)
	if err := os.MkdirAll(folder, 0755); err != nil {
		return err
	}
	return os.Symlink(gxModPath, targetLink)
}

func cBind(cmd *cobra.Command, args []string) error {
	mod, err := gxmodule.Current()
	if err != nil {
		return err
	}
	pkgs, err := gxtc.Packages(mod)
	if err != nil {
		return err
	}
	if len(pkgs) == 0 {
		return nil
	}
	depsPath := filepath.Join(mod.Path(), "gxdeps")
	if err := os.MkdirAll(depsPath, 0755); err != nil {
		return err
	}
	for _, dep := range mod.Deps() {
		if err := installLinkToModule(mod, depsPath, dep); err != nil {
			return err
		}
	}
	for _, pkg := range pkgs {
		target := filepath.Join(depsPath, pkg)
		if err := gxtc.Bind(mod, pkg, target); err != nil {
			return err
		}
	}
	return nil
}
