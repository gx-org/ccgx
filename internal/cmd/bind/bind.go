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
	"github.com/gx-org/ccgx/internal/gotc"
	"github.com/gx-org/ccgx/internal/gxtc"
	gxmodule "github.com/gx-org/gx/build/module"
	"github.com/spf13/cobra"
)

var cmake bool

// Cmd is the implementation of the mod command.
func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "Create links to dependencies, then generate C++ header files",
		RunE:  cBind,
	}
	cmd.PersistentFlags().BoolVarP(&cmake, "cmake", "", false, "generate CMakeLists.txt")
	return cmd
}

func cBind(cmd *cobra.Command, args []string) error {
	mod, err := gxmodule.Current()
	if err != nil {
		return err
	}
	cache, err := gotc.NewCache()
	if err != nil {
		return err
	}
	if err := gxtc.LinkAllDeps(mod, cache); err != nil {
		return err
	}
	var fs []gxtc.BinderCallback
	if cmake {
		fs = append(fs, gxtc.WriteCMakeLists)
	}
	if err := gxtc.BindAll(mod, fs); err != nil {
		return err
	}
	if err := gxtc.CompileCArchive(mod); err != nil {
		return err
	}
	return nil
}
