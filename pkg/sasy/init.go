package sasy 

import (
	"fmt"
	"os"
	"path"
)

func InitHandler() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	gitPath := path.Join(wd, ".git")

	if err := os.Mkdir(gitPath, os.ModePerm); err != nil {
		return err
	}
	subdirs := []string{"objects", "refs"}
	for i := range subdirs {
		if err = os.Mkdir(path.Join(gitPath, subdirs[i]), os.ModePerm); err != nil {
			return err
		}
	}
	fmt.Println("Initialized empty sasy repository in", wd)
	return nil
}
