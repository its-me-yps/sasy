package sasy

import (
	"flag"
	"fmt"
	"sasy/pkg/model"
	"sasy/utils"
	"strings"
)

func CommitHandler(args []string) error {

	// Parsing commit message from command line args
	fs := flag.NewFlagSet("commit", flag.ExitOnError)
	commitMessage := ""
	fs.StringVar(&commitMessage, "message", "", "commit message")
	fs.StringVar(&commitMessage, "m", "", "commit message")
	fs.Parse(args)

	if len(commitMessage) == 0 {
		return fmt.Errorf("error: empty commit message")
	}
	database, err := model.CreateDatabase(utils.WorkindDir)
	if err != nil {
		return fmt.Errorf("error in creating database: %v", err)
	}

	rootTree, err := model.CreateTree(utils.WorkindDir)
	if err != nil {
		return err
	}

	// saving the root tree
	if err := database.Save(rootTree.Oid, []byte(rootTree.Content)); err != nil {
		return fmt.Errorf("error saving the root-tree %v", err)
	}

	// saving all entries(subtrees and blobs) of root-tree
	if err := model.SaveObjectsToDatabase(rootTree.Entries, database); err != nil {
		return err
	}

	author := model.NewAuthor()

	refs := model.Refs{Path: utils.SasyPath}
	// Reading commit id of parent from refs
	parent, err := refs.ReadHead()
	if err != nil {
		return fmt.Errorf("error in Reading refs: %v", err)
	}

	commit := model.CreateCommit(parent, rootTree.Oid, *author, commitMessage)

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
