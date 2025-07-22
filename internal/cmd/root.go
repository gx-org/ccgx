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

package cmd

import (
	"github.com/gx-org/ccgx/internal/cmd/bind"
	"github.com/gx-org/ccgx/internal/cmd/debug"
	"github.com/gx-org/ccgx/internal/cmd/link"
	"github.com/gx-org/ccgx/internal/cmd/mod"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "ccgx",
	Short:         "Generate C++ bindings for GX.",
	SilenceUsage:  true,
	SilenceErrors: false,
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debug.Debug, "debug", "d", false, "print debug information")
	rootCmd.AddCommand(mod.Cmd)
	rootCmd.AddCommand(link.Cmd())
	rootCmd.AddCommand(bind.Cmd())
}
