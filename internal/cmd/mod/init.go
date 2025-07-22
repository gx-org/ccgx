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

package mod

import (
	"github.com/gx-org/ccgx/internal/gotc"
	"github.com/gx-org/ccgx/internal/gxtc"
	"github.com/gx-org/ccgx/internal/module"
	"github.com/spf13/cobra"
)

var initOverwriteFile bool

func cmdInit() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "create gx.mod",
		RunE:  cInit,
		Args:  cobra.ExactArgs(1),
	}
}

func cInit(cmd *cobra.Command, args []string) error {
	if err := gotc.ModInit(args[0]); err != nil {
		return err
	}
	mod, err := module.Current()
	if err != nil {
		return err
	}
	if err := gxtc.PackAll(mod); err != nil {
		return err
	}
	if err := gotc.ModTidy(); err != nil {
		return err
	}
	return nil
}
