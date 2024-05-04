package sasy 

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
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
	o.setOid()
	return o
}

func (o *Blob) setContent() {
	content, err := os.ReadFile(o.Path)
	if err != nil {
		fmt.Print("Error Occured while reading", o.Path, err)
	}
	o.Content = fmt.Sprintf("%s %d\x00%s", "blob", len(content), string(content))
}

func (o *Blob) setOid() {
	h := sha1.New()
	io.WriteString(h, o.Content)
	o.Oid = hex.EncodeToString(h.Sum(nil))
}
