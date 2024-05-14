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
    return fmt.Errorf("Error in creating database: %v", err)
	}
  
	blobEntries := []*model.Blob{}
	for _, file := range files {
		blob := model.CreateBlob(wd, file)

		if err := database.Save(blob.Oid, []byte(blob.Content)); err != nil {
      return fmt.Errorf("Error saving blob %s in db: %v", blob.Name, err)
		}
		blobEntries = append(blobEntries, blob)
	}

	tree := model.CreateTree(&blobEntries)
	if err := database.Save(tree.Oid, []byte(tree.Content)); err != nil {
    return fmt.Errorf("Error saving tree %s in db: %v", tree.Oid, err)
	}

	// Author info from environment variable
	name := os.Getenv("SASY_AUTHOR_NAME")
	email := os.Getenv("SASY_AUTHOR_EMAIL")
	author := model.Author{Name: name, Email: email, T: time.Now()}

	// commit message from stdin
	stdin, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}
	m := string(stdin)

  sasyPath := path.Join(wd, ".sasy")
  refs := model.Refs{Path: sasyPath}
  // Reading commit id of parent from refs
  parent, err := refs.ReadHead()
  if err != nil {
    return fmt.Errorf("Error in Reading refs: %v", err)
  }

	commit := model.CreateCommit(parent, tree.Oid, author, m)
	database.Save(commit.Oid, []byte(commit.Content))
  // TODO: Implement LockFile to safely update HEAD in race condition
  if err := refs.UpdateHead(commit.Oid); err != nil {
    return fmt.Errorf("Error in updating HEAD: %v", err)
  }

  // Stdout msg on succesfull commit
  commitStdout := "[(root-commit)" // For the first commit
  // if parent != "" => non root commit
  if parent != "" {
    commitStdout = "[main"
  }

	fmt.Printf("%s %s] ",commitStdout, commit.Oid)
	fmt.Printf("%s\n", strings.Split(m, "\n")[0])
	return nil
}
