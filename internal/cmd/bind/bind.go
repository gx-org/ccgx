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
	"ccgx/internal/module"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Cmd is the implementation of the mod command.
var Cmd = &cobra.Command{
	Use:   "bind",
	Short: "Generate C++ header files",
	RunE:  cBind,
}

func installLinkToGX(targetPath string) error {
	mod, err := module.Current()
	if err != nil {
		return err
	}
	gxVersion := mod.GXVersion()
	if gxVersion == "" {
		return fmt.Errorf("cannot find GX version")
	}
	gxModPath, err := gotc.ModuleOSPath(module.GXModulePath, gxVersion)
	if err != nil {
		return err
	}
	return os.Symlink(gxModPath, filepath.Join(targetPath, "gx"))
}

func cBind(cmd *cobra.Command, args []string) error {
	mod, pkgs, err := gxtc.Packages()
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
	if err := installLinkToGX(depsPath); err != nil {
		return err
	}
	for _, pkg := range pkgs {
		target := filepath.Join(depsPath, pkg)
		if err := gxtc.Bind(pkg, target); err != nil {
			return err
		}
	}
	return nil
}
