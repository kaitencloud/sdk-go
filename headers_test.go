// Copyright 2026 KAITEN INC
// SPDX-License-Identifier: Apache-2.0

package sdk

import (
	"bufio"
	"io/fs"
	"os"
	"strings"
	"testing"
)

// licenseHeader opens every Go file of the module, the generated clients
// included: `task generate` and `task generate:platform` write it too. It is
// the licensing pack's LICENSE_HEADER.txt, in Go comment syntax.
var licenseHeader = []string{
	"// Copyright 2026 KAITEN INC",
	"// SPDX-License-Identifier: Apache-2.0",
}

func TestLicenseHeaders(t *testing.T) {
	module := os.DirFS(".")
	err := fs.WalkDir(module, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != "." && (strings.HasPrefix(d.Name(), ".") || d.Name() == "testdata") {
				return fs.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		missing, err := lacksLicenseHeader(module, path)
		if err != nil {
			return err
		}
		if missing {
			t.Errorf("%s does not open with the license header:\n%s", path, strings.Join(licenseHeader, "\n"))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

func lacksLicenseHeader(module fs.FS, path string) (bool, error) {
	f, err := module.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for _, want := range licenseHeader {
		if !scanner.Scan() || scanner.Text() != want {
			return true, scanner.Err()
		}
	}
	return false, nil
}
