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

package gxtc

import (
	"fmt"
	"io"
	"os"
)

// copy a file from src to dst.
func copy(src, dst string) error {
	srcStat, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("cannot copy %s: %v", src, err)
	}
	if !srcStat.Mode().IsRegular() {
		return fmt.Errorf("cannot copy non-regular %s file %s", srcStat.Mode().String(), srcStat.Name())
	}
	dstStat, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("cannot copy %s to %s: %v", src, dst, err)
		}
	} else {
		if !(dstStat.Mode().IsRegular()) {
			return fmt.Errorf("destination %s %s is non-regular", dstStat.Name(), dstStat.Mode().String())
		}
		if os.SameFile(srcStat, dstStat) {
			return nil
		}
	}
	if err = os.Link(src, dst); err == nil {
		return nil
	}
	return copyContent(src, dst)
}

func copyContent(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("cannot open source %s: %v", src, err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("cannot destination %s: %v", dst, err)
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return fmt.Errorf("copy error: %v", err)
	}
	return out.Sync()
}
