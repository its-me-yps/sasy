package model

import (
	"fmt"
	"os"
	"sasy/utils"
)

type Blob struct {
	Path    string
	Oid     string
	Mode    string
	Content string
	Name    string
}

func CreateBlob(path string, file os.DirEntry) *Blob {
	o := &Blob{}
	o.Name = file.Name()
	o.Path = path
	o.Mode = utils.GetFileMode(file)
	o.setContent()
	o.Oid = utils.CalculateSHA1(o.Content)
	return o
}

func (o *Blob) setContent() {
	content, err := os.ReadFile(o.Path)
	if err != nil {
		fmt.Print("Error Occured while reading", o.Path, err)
	}
	o.Content = fmt.Sprintf("%s %d\x00%s", "blob", len(content), string(content))
}
