package utils

import "os"

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}
