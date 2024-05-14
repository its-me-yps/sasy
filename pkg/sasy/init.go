package sasy 

import (
	"fmt"
	"os"
	"path"
)

// TODO: Change Permission Mode for the directories and files
func InitHandler() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	sasyPath := path.Join(wd, ".sasy")

	if err := os.Mkdir(sasyPath, os.ModePerm); err != nil {
    return fmt.Errorf("Error in creating .sasy: %v", err)
	}
	subdirs := []string{"objects", "refs"}
	for i := range subdirs {
		if err = os.Mkdir(path.Join(sasyPath, subdirs[i]), os.ModePerm); err != nil {
			return fmt.Errorf("Error in creating %s; %v", subdirs, err)
		}
	}

  headPath := path.Join(sasyPath, "HEAD")
  if err := os.WriteFile(headPath, []byte{}, os.ModePerm); err != nil {
    return fmt.Errorf("Error in creating HEAD: %v", err)
  }

	fmt.Printf("Initialized empty sasy repository in %s\n", sasyPath)
	return nil
}

