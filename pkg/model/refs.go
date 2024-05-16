package model

import (
	"os"
	"path"
)

type Refs struct {
	Path string
}

func (r *Refs) UpdateHead(Oid string) error {
	err := os.WriteFile(path.Join(r.Path, "HEAD"), []byte(Oid), os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func (r *Refs) ReadHead() (string, error) {
	headPath := path.Join(r.Path, "HEAD")
	parent := ""
	data, err := os.ReadFile(headPath)
	if err != nil {
		if err := os.WriteFile(headPath, []byte(parent), os.ModeDevice); err != nil {
			return parent, err
		}
	}
	return string(data), nil
}
