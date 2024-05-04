package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
)

type Tree struct {
	blobs   []*Blob
	Oid     string
	Content string
}

func CreateTree(blobs *([]*Blob)) *Tree {
	t := &Tree{}
	t.blobs = *blobs
	t.setContent()
	t.setOid()
	return t
}

func (t *Tree) sortBlobs() {
	sort.Slice(t.blobs, func(i, j int) bool {
		return t.blobs[i].Name < t.blobs[j].Name
	})
}

func (t *Tree) setContent() {
	t.sortBlobs()
	for _, b := range t.blobs {
		t.Content += fmt.Sprintf("%d %s %s\x00", 100644, b.Name, b.Oid)
	}
	t.Content = fmt.Sprintf("%s %d\x00%s", "tree", len(t.Content), t.Content)
}

func (t *Tree) setOid() {
	h := sha1.New()
	io.WriteString(h, t.Content)
	t.Oid = hex.EncodeToString(h.Sum(nil))
}
