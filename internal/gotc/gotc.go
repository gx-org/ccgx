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
	"path/filepath"
	"strings"

	"golang.org/x/mod/module"
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

const goModCache = "GOMODCACHE"

// Cache stores the path where Go caches modules.
type Cache struct {
	path string
}

// NewCache reads where Go caches modules.
func NewCache() (*Cache, error) {
	cmd := exec.Command("go", "env")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	for line := range strings.Lines(string(out)) {
		spLine := strings.Split(line, "=")
		if len(spLine) != 2 {
			continue
		}
		if spLine[0] != goModCache {
			continue
		}
		cachePath := strings.TrimSpace(spLine[1])
		cachePath = strings.TrimPrefix(cachePath, "'")
		cachePath = strings.TrimSuffix(cachePath, "'")
		return &Cache{path: cachePath}, nil
	}
	return nil, fmt.Errorf("Go variable environment %s not found", goModCache)
}

// ModuleOSPath returns the OS path of a given module path and its version.
func (cache *Cache) OSPath(mod *module.Version) (string, error) {
	dir, file := filepath.Split(mod.Path)
	modCache := filepath.Join(cache.path, dir, fmt.Sprintf("%s@%s", file, mod.Version))
	dirStat, err := os.Stat(modCache)
	if err != nil {
		return "", fmt.Errorf("cannot find GX module path at %s: %v\nPlease run ccgx tidy.\n", modCache, err)
	}
	if !dirStat.IsDir() {
		return "", fmt.Errorf("%s not a directory: %v\nPlease run ccgx tidy.\n", modCache, err)
	}
	return modCache, nil
}

func BuildArchive(root, src, target string) error {
	cmd := exec.Command("go", "build",
		"-buildmode=c-archive",
		"-o",
		target,
		src)
	const cflagsKey = "CGO_CFLAGS"
	var cflagsValue string
	envs := os.Environ()
	for _, env := range envs {
		keyvals := strings.Split(env, "=")
		if len(keyvals) != 2 {
			cmd.Env = append(cmd.Env, env)
			continue
		}
		key, val := keyvals[0], keyvals[1]
		if key != cflagsKey {
			cmd.Env = append(cmd.Env, env)
			continue
		}
		cflagsValue = val
	}
	cmd.Env = append(cmd.Env, cflagsKey+"="+cflagsValue+" -I "+root)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
