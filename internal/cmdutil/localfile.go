// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package cmdutil

import (
	"fmt"
	"io/fs"

	"github.com/larksuite/cli/extension/fileio"
	"github.com/larksuite/cli/internal/validate"
	"github.com/larksuite/cli/internal/vfs"
)

// OpenLocalFile opens a regular file in the process filesystem namespace.
// Absolute and relative paths are accepted. It is the shared replacement for
// direct os.Open/os.ReadFile use in commands that intentionally read local
// paths outside the workspace sandbox.
func OpenLocalFile(path string) (fs.File, error) {
	localPath, err := validate.LocalInputPath(path)
	if err != nil {
		return nil, &fileio.PathValidationError{Err: err}
	}

	info, err := vfs.Stat(localPath)
	if err != nil {
		return nil, err
	}
	if !info.Mode().IsRegular() {
		return nil, &fileio.PathValidationError{Err: fmt.Errorf("%q is not a regular file", path)}
	}

	f, err := vfs.Open(localPath)
	if err != nil {
		return nil, err
	}
	openedInfo, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, err
	}
	if !openedInfo.Mode().IsRegular() {
		_ = f.Close()
		return nil, &fileio.PathValidationError{Err: fmt.Errorf("%q is not a regular file", path)}
	}
	return f, nil
}
