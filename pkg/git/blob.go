package git

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

type Object struct {
	path       string
	objectType string
	oid        string
	content    string
	compressed []byte
}

func CreateObject(path string, objectType string) *Object {
	o := &Object{}
	o.path = path
	o.objectType = objectType
	o.setContent()
	o.setOid()
	o.compress()
	return o
}

func (o *Object) setContent() {
	content, err := os.ReadFile(o.path)
	if err != nil {
		fmt.Print("Error Occured while reading", o.path, err)
	}
	o.content = fmt.Sprintf("%s %d\x00%s", o.objectType, len(content), string(content))
}

func (o *Object) setOid() {
	h := sha1.New()
	io.WriteString(h, o.content)
	o.oid = hex.EncodeToString(h.Sum(nil))
}

func (o *Object) compress() {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(o.content))
	w.Close()
	o.compressed = b.Bytes()
}

func (o *Object) getOid() string {
	return o.oid
}

func (o *Object) getCompressed() []byte {
	return o.compressed
}
