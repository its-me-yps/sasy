package sasy 

import (
	"bytes"
	"compress/zlib"
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
	d.objectsDir = path.Join(d.workingDir, ".git", "objects")
	return d, nil
}

func (d *Database) Save(oid string, content []byte) error {
	subDir := path.Join(d.objectsDir, oid[:2])
	fileName := path.Join(subDir, oid[2:])
	if err := os.Mkdir(subDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(fileName, compress(content), 0644); err != nil {
		return err
	}
	return nil
}

func compress(content []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(content)
	w.Close()
	return b.Bytes()
}
