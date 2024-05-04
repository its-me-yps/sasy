package sasy

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"time"

	"sasy/pkg/model"
	"sasy/utils"
)

func CommitHandler() error {
	files, _ := utils.Ls()
	wd, _ := os.Getwd()
	database, err := model.CreateDatabase(wd)
	if err != nil {
		return err
	}

	blobEntries := []*model.Blob{}
	for _, file := range files {
		blob := model.CreateBlob(wd, file)

		if err := database.Save(blob.Oid, []byte(blob.Content)); err != nil {
			return err
		}
		blobEntries = append(blobEntries, blob)
	}

	tree := model.CreateTree(&blobEntries)
	if err := database.Save(tree.Oid, []byte(tree.Content)); err != nil {
		return err
	}

	// Author info from environment variable
	name := os.Getenv("GIT_AUTHOR_NAME")
	email := os.Getenv("GIT_AUTHOR_EMAIL")
	author := model.Author{Name: name, Email: email, T: time.Now()}

	// commit message from stdin
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	m := string(stdin)

	commit := model.CreateCommit(tree.Oid, author, m)
	database.Save(commit.Oid, []byte(commit.Content))

	if err := os.WriteFile(path.Join(wd, ".git", "HEAD"), []byte(commit.Oid), 0644); err != nil {
		return err
	}
	fmt.Printf("[(root-commit) %s] ", commit.Oid)
	fmt.Printf("%s", strings.Split(m, "\n")[0])
	return nil
}
