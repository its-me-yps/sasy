package model

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path"
	"sasy/utils"
)

type DatabaseObject interface {
	GetContent() string
}

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

func (d *Database) Save(obj DatabaseObject) error {
	oid := getOid(obj)
	subDir := path.Join(d.objectsDir, oid[:2])
	fileName := path.Join(subDir, oid[2:])

	//Return if fileName already exist(This happens if contents of file have not been changed so oid remains same as before)
	if _, err := os.Stat(fileName); err == nil {
		return nil
	}

	if err := os.Mkdir(subDir, os.ModePerm); err != nil {
		return err
	}
	if err := os.WriteFile(fileName, utils.Compress([]byte(obj.GetContent())), 0644); err != nil {
		return err
	}
	return nil
}

func getOid(d DatabaseObject) string {
	h := sha1.New()
	io.WriteString(h, d.GetContent())
	return hex.EncodeToString(h.Sum(nil))
}
