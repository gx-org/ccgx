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

// Package gotc provides functions to invoke the Go toolchain.
package gotc

import (
	"fmt"
	"os"
	"os/exec"
)

// Check that Go is installed.
func Check() error {
	cmd := exec.Command("go", "version")
	if _, err := cmd.Output(); err != nil {
		return fmt.Errorf("invalid Go installation: %v", err)
	}
	return nil
}

// ModInit runs the go mod init command.
func ModInit(modName string) error {
	cmd := exec.Command("go", "mod", "init", modName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// ModTidy runs the go mod tidy command.
func ModTidy() error {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
