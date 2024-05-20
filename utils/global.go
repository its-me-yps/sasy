package utils

import (
	"os"
	"path"
)

// Directories to ignore
var IgnoreDirs = []string{".sasy"}

var (
	SasyPath   string
	WorkindDir string
)

func Init() error {
	if wd, err := os.Getwd(); err != nil {
		return err
	} else {
		WorkindDir = wd
	}
	SasyPath = path.Join(WorkindDir, ".sasy")
	return nil
}
