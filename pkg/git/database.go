package git

import (
	"os"
	"path"
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

func (d *Database) SaveObject(o *Object) error {
	subDir := path.Join(d.objectsDir, o.Oid[:2])
	fileName := path.Join(subDir, o.Oid[2:])
	if err := os.Mkdir(subDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(fileName, o.Compressed, 0644); err != nil {
		return err
	}
	return nil
}
