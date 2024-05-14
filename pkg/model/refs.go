package model

import (
	"io"
	"os"
	"path"
)

type Refs struct {
	Path string
}

func (r *Refs) UpdateHead(Oid string) error {
	file, err := os.OpenFile(path.Join(r.Path, "HEAD"), os.O_WRONLY | os.O_TRUNC, os.ModePerm)
  if err != nil {
    return err
  }

  if _, err := io.WriteString(file, Oid); err != nil {
    return err
  }

  return nil
}

func (r *Refs) ReadHead() (string, error) {
	file, err := os.Open(path.Join(r.Path, "HEAD"))
	if err != nil {
		return "", err
	}
  defer file.Close()

	head, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	return string(head), nil
}
