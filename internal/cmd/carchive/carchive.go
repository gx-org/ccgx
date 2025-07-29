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

// Package carchive provides the Cobra carchive command.
// The carchive command creates a carchive.go with all the Go/GX dependencies
// to build a binary c archive that can be linked with the final executable.
package carchive

import (
	"github.com/gx-org/ccgx/internal/gxtc"
	gxmodule "github.com/gx-org/gx/build/module"
	"github.com/spf13/cobra"
)

// Cmd is the implementation of the mod command.
func Cmd() *cobra.Command {
	return &cobra.Command{
		Use:   "carchive",
		Short: "Create a c archive",
		Long:  "First, create a carchive.go file which includes all the Go/GX dependencies and a main function. This file is then compile using `go build -buildmode=c-archive` to produce a binary static .a library file.",
		RunE:  cArchive,
	}
}

func cArchive(cmd *cobra.Command, args []string) error {
	mod, err := gxmodule.Current()
	if err != nil {
		return err
	}
	depsPath, err := gxtc.DepsPath(mod)
	if err != nil {
		return err
	}
	return gxtc.CompileCArchive(mod, depsPath)
}
