package model

import (
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

	//Return if fileName already exist(This happens if contents of file have not been changed so oid remains same as before)
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
