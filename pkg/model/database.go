package model

import (
	"fmt"
	"os"
	"path"
	"sasy/utils"
)

type Database struct {
	workingDir string
	objectsDir string
}

func CreateDatabase(workdir string) (*Database, error) {
	d := &Database{}
	d.workingDir = workdir
	d.objectsDir = path.Join(d.workingDir, ".sasy", "objects")
	return d, nil
}

func (d *Database) Save(oid string, content []byte) error {
	subDir := path.Join(d.objectsDir, oid[:2])
	fileName := path.Join(subDir, oid[2:])

	// Return if fileName already exist(This happens if contents of file have not been changed so oid remains same as before)
	if _, err := os.Stat(fileName); err == nil {
		return nil
	}

	if err := os.MkdirAll(subDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(fileName, utils.Compress(content), 0444); err != nil {
		return err
	}
	return nil
}

func (d *Database) Read(oid string) ([]byte, error) {
	if len(oid) != 40 {
		return nil, fmt.Errorf("%s is not a valid object ID", oid)
	}
	filePath := path.Join(d.objectsDir, oid[:2], oid[2:])
	if !utils.CheckFileExists(filePath) {
		return nil, fmt.Errorf("no such object")
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}
