//go:build linux || freebsd
// +build linux freebsd

package system

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
)

// prepareFiles creates files for testing in the temp directory
func prepareFiles(t *testing.T) (string, string, string, string) {
	dir, err := os.MkdirTemp("", "bhojpur-system-test")
	if err != nil {
		t.Fatal(err)
	}

	file := filepath.Join(dir, "exist")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	invalid := filepath.Join(dir, "doesnt-exist")

	symlink := filepath.Join(dir, "symlink")
	if err := os.Symlink(file, symlink); err != nil {
		t.Fatal(err)
	}

	return file, invalid, symlink, dir
}

func TestLUtimesNano(t *testing.T) {
	file, invalid, symlink, dir := prepareFiles(t)
	defer os.RemoveAll(dir)

	before, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}

	ts := []syscall.Timespec{{Sec: 0, Nsec: 0}, {Sec: 0, Nsec: 0}}
	if err := LUtimesNano(symlink, ts); err != nil {
		t.Fatal(err)
	}

	symlinkInfo, err := os.Lstat(symlink)
	if err != nil {
		t.Fatal(err)
	}
	if before.ModTime().Unix() == symlinkInfo.ModTime().Unix() {
		t.Fatal("The modification time of the symlink should be different")
	}

	fileInfo, err := os.Stat(file)
	if err != nil {
		t.Fatal(err)
	}
	if before.ModTime().Unix() != fileInfo.ModTime().Unix() {
		t.Fatal("The modification time of the file should be same")
	}

	if err := LUtimesNano(invalid, ts); err == nil {
		t.Fatal("Doesn't return an error on a non-existing file")
	}
}
