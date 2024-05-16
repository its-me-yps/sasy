package model

import (
	"fmt"
	"os"
	"path"
	"sasy/utils"
)

type Blob struct {
	Path    string
	Oid     string
	Content string
	Name    string
}

func CreateBlob(wd string, name string) *Blob {
	o := &Blob{}
	o.Name = name
	o.Path = path.Join(wd, name)
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
