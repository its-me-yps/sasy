package sasy

import (
	"fmt"
	"os"
	"path"
	"sasy/utils"
)

// TODO: Change Permission Mode for the directories and files
func InitHandler(args []string) error {

	if err := os.MkdirAll(utils.SasyPath, os.ModePerm); err != nil {
		return fmt.Errorf("error in creating .sasy: %v", err)
	}
	subdirs := []string{"objects", "refs"}
	for i := range subdirs {
		if err := os.MkdirAll(path.Join(utils.SasyPath, subdirs[i]), os.ModePerm); err != nil {
			return fmt.Errorf("error in creating %s; %v", subdirs, err)
		}
	}

	headPath := path.Join(utils.SasyPath, "HEAD")
	if err := os.WriteFile(headPath, []byte{}, os.ModePerm); err != nil {
		return fmt.Errorf("error in creating HEAD: %v", err)
	}

	fmt.Printf("Initialized empty Sasy repository in %s\n", utils.SasyPath)
	return nil
}
