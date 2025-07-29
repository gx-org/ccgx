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
	"github.com/gx-org/gx/build/module"
	"github.com/spf13/cobra"
)

func cmdTidy() *cobra.Command {
	return &cobra.Command{
		Use:   "tidy",
		Short: "update gx.mod",
		RunE:  cTidy,
	}
}

func cTidy(cmd *cobra.Command, args []string) error {
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
