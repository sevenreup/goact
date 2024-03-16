package utils

import (
	"fmt"
	"os"
	"strings"
)

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func FormatPath(p string) string {
	if strings.HasPrefix(p, "./") {
		return p
	}
	return fmt.Sprintf("./%s", p)
}
