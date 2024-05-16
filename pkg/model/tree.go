package model

import (
	"encoding/hex"
	"fmt"
	"os"
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
		info, err := os.Stat(b.Path)
		if err != nil {
			fmt.Printf("Error in reading Stats of file %s: %v\n", b.Path, err)
		}
		var mode string
		if info.Mode()&0100 != 0 {
			mode = "100755"
		} else {
			mode = "100644"
		}
		// Decoding 40 hex character long b.Oid to []byte of lenght 20
		byteSlice, _ := hex.DecodeString(b.Oid)

		t.Content += fmt.Sprintf("%s %s\x00%s", mode, b.Name, string(byteSlice))
	}
	t.Content = fmt.Sprintf("%s %d\x00%s", "tree", len(t.Content), t.Content)
}

func (t *Tree) GetContent() string {
	return t.Content
}
