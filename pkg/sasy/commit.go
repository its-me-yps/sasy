package sasy

import (
	"fmt"
	"os"

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
	fmt.Println("Commited Successfully")
	return nil
}
