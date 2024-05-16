package sasy

import (
	"flag"
	"fmt"
	"os"
	"path"
	"sasy/pkg/model"
	"sasy/utils"
	"strings"
)

func CommitHandler(args []string) error {

	// Parsing commit message from command line args
	fs := flag.NewFlagSet("Commit", flag.ExitOnError)
	commitMessage := ""
	fs.StringVar(&commitMessage, "message", "", "commit message")
	fs.StringVar(&commitMessage, "m", "", "commit message")
	fs.Parse(args)

	if len(commitMessage) == 0 {
		return fmt.Errorf("error: empty commit message")
	}
	wd, _ := os.Getwd()
	database, err := model.CreateDatabase(wd)
	if err != nil {
		return fmt.Errorf("error in creating database: %v", err)
	}

	blobEntries := []*model.Blob{}
	if err := saveBlobsToDatabase(&blobEntries, database); err != nil {
		return err
	}

	tree := model.CreateTree(&blobEntries)
	if err := database.Save(tree.Oid, []byte(tree.Content)); err != nil {
		return fmt.Errorf("error saving tree %s in db: %v", tree.Oid, err)
	}

	author := model.NewAuthor()

	sasyPath := path.Join(wd, ".sasy")
	refs := model.Refs{Path: sasyPath}
	// Reading commit id of parent from refs
	parent, err := refs.ReadHead()
	if err != nil {
		return fmt.Errorf("error in Reading refs: %v", err)
	}

	commit := model.CreateCommit(parent, tree.Oid, *author, commitMessage)

	if err := database.Save(commit.Oid, []byte(commit.Content)); err != nil {
		return fmt.Errorf("error saving commit %v", err)
	}
	// TODO: Implement LockFile to safely update HEAD in race condition
	if err := refs.UpdateHead(commit.Oid); err != nil {
		return fmt.Errorf("error in updating HEAD: %v", err)
	}

	// Stdout msg on succesfull commit
	commitStdout := "[(root-commit)" // For the first commit
	// if parent != "" => non root commit
	if parent != "" {
		commitStdout = "[main"
	}

	fmt.Printf("%s %s] ", commitStdout, commit.Oid)
	fmt.Printf("%s\n", strings.Split(commitMessage, "\n")[0])
	return nil
}

func saveBlobsToDatabase(blobEntries *[]*model.Blob, database *model.Database) error {
	var files []string
	var wd string
	var err error
	if files, wd, err = utils.Ls(); err != nil {
		return fmt.Errorf("error reading files in db: %v", err)
	}
	for _, file := range files {
		blob := model.CreateBlob(wd, file)

		if err := database.Save(blob.Oid, []byte(blob.Content)); err != nil {
			return fmt.Errorf("error saving blob %s in db: %v", blob.Name, err)
		}
		*blobEntries = append(*blobEntries, blob)
	}
	return nil
}
