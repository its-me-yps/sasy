package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
)

type Commit struct {
	// Oid of tree that commit points to
	TreeId string
	// Oid of the commit
	Oid     string
	Author  Author
	Message string
	// content to be stored
	Content string
}

func CreateCommit(treeOid string, a Author, m string) *Commit {
	c := &Commit{}
	c.TreeId = treeOid
	c.Author = a
	c.Message = m
	c.setContent()
	c.setOid()
	return c
}

// content of commit blob that is to be stored in database
func (c *Commit) setContent() {
	c.Content = fmt.Sprintf("tree %s\nauthor %s\ncommitter %s\n\n%s", c.TreeId, c.Author.toStr(), c.Author.toStr(), c.Message)
}

func (c *Commit) setOid() {
	h := sha1.New()
	io.WriteString(h, c.Content)
	c.Oid = hex.EncodeToString(h.Sum(nil))
}
