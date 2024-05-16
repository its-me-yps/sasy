package model

import (
	"fmt"
	"sasy/utils"
)

type Commit struct {
	Parent  string // Oid of parent commit
	TreeId  string // Oid of tree that commit points to
	Oid     string // Oid of the commit
	Author  Author
	Message string
	Content string // content to be stored
}

func CreateCommit(parent string, treeOid string, a Author, m string) *Commit {
	c := &Commit{}
	c.Parent = parent
	c.TreeId = treeOid
	c.Author = a
	c.Message = m
	c.setContent()
	c.Oid = utils.CalculateSHA1(c.Content)
	return c
}

// content of commit blob that is to be stored in database
func (c *Commit) setContent() {
	var s string
	if c.Parent != "" {
		s = fmt.Sprintf("parent %s\n", c.Parent)
	}

	c.Content = fmt.Sprintf("tree %s\n%sauthor %s\ncommitter %s\n\n%s", c.TreeId, s, c.Author.toStr(), c.Author.toStr(), c.Message)
}
