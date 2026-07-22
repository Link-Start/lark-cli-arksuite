// Copyright (c) 2026 Lark Technologies Pte. Ltd.
// SPDX-License-Identifier: MIT

package cmdutil

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/larksuite/cli/extension/fileio"
)

func TestOpenLocalFile_AcceptsAbsoluteAndParentRelativePaths(t *testing.T) {
	root := t.TempDir()
	workDir := filepath.Join(root, "work")
	if err := os.Mkdir(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, "input.txt")
	if err := os.WriteFile(path, []byte("content"), 0o600); err != nil {
		t.Fatal(err)
	}
	TestChdir(t, workDir)

	for _, input := range []string{path, filepath.Join("..", "input.txt")} {
		f, err := OpenLocalFile(input)
		if err != nil {
			t.Fatalf("OpenLocalFile(%q) error = %v", input, err)
		}
		got, readErr := io.ReadAll(f)
		closeErr := f.Close()
		if readErr != nil || closeErr != nil || string(got) != "content" {
			t.Fatalf("OpenLocalFile(%q) content=%q read=%v close=%v", input, got, readErr, closeErr)
		}
	}
}

func TestOpenLocalFile_RejectsInvalidAndNonRegularInputs(t *testing.T) {
	for _, input := range []string{"input\n.txt", t.TempDir()} {
		if _, err := OpenLocalFile(input); !errors.Is(err, fileio.ErrPathValidation) {
			t.Fatalf("OpenLocalFile(%q) error = %v, want ErrPathValidation", input, err)
		}
	}
}
