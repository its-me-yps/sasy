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
	d.objectsDir = path.Join(d.workingDir, ".git", "objects")
	return d, nil
}

func (d *Database) Save(oid string, content []byte) error {
	subDir := path.Join(d.objectsDir, oid[:2])
	fileName := path.Join(subDir, oid[2:])
	if err := os.Mkdir(subDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(fileName, utils.Compress(content), 0644); err != nil {
		return err
	}
	return nil
}
