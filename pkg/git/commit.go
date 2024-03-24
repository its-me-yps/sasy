package git

import (
	"fmt"
	"os"
)

func CommitHandler() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	// vcPath := path.Join(wd, ".sasy")
	// dbPath := path.Join(vcPath, "objects")
	dir, _ := os.Open(wd)
	files, _ := dir.Readdir(0)
	for _, file := range files {
		fmt.Println(file.Name())
	}
	return nil
}
