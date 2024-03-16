package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func FormatPath(p string) string {
	absPath, err := filepath.Abs(p)
	if err != nil {
		return p
	}

	if !filepath.IsAbs(absPath) {
		return p
	}

	relPath, err := filepath.Rel(".", absPath)
	if err != nil {
		return p
	}

	if relPath == p {
		return p
	}

	return fmt.Sprintf(".%s", string(filepath.Separator)+relPath)
}
