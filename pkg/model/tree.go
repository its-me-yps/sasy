package model

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sasy/utils"
)

type Tree struct {
	Entries []Object
	Oid     string
	Content string
	Path    string
	Name    string
}

// Takes Current Directory and creates tree from it
func CreateTree(cd string) (*Tree, error) {
	t := &Tree{}
	t.Path = cd

	dirEntries, err := os.ReadDir(cd)
	if err != nil {
		return nil, fmt.Errorf("error reading dir %s: %v\n", cd, err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.Name() == ".sasy" {
			continue
		}

		if dirEntry.IsDir() {
			subtree, err := CreateTree(path.Join(cd, dirEntry.Name()))
			subtree.Name = dirEntry.Name()
			if err != nil {
				return nil, fmt.Errorf("error in creating tree object for dir %s: %v\n", path.Join(cd, dirEntry.Name()), err)
			}
			t.Entries = append(t.Entries, subtree)
		} else {
			blob := CreateBlob(path.Join(cd, dirEntry.Name()), dirEntry)
			t.Entries = append(t.Entries, blob)
		}
	}

	t.setContent()
	t.Oid = utils.CalculateSHA1(t.Content)
	return t, nil
}

//	func (t *Tree) sortBlobs() {
//		sort.Slice(t.blobs, func(i, j int) bool {
//			return t.blobs[i].Name < t.blobs[j].Name
//		})
//	}

func (t *Tree) setContent() {
	for _, entry := range t.Entries {
		var mode string
		var Oid string
		var Name string
		switch entry.(type) {
		case *Tree:
			subTree, _ := entry.(*Tree)
			mode = "040000"
			Oid = subTree.Oid
			Name = subTree.Name
		case *Blob:
			blob, _ := entry.(*Blob)
			mode = blob.Mode 
			Oid = blob.Oid
			Name = blob.Name
		}
		// Decoding 40 hex character long Oid to []byte of lenght 20
		byteSlice, _ := hex.DecodeString(Oid)

		t.Content += fmt.Sprintf("%s %s\x00%s", mode, Name, string(byteSlice))
	}
	t.Content = fmt.Sprintf("%s %d\x00%s", "tree", len(t.Content), t.Content)
}
