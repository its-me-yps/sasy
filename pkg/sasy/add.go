package sasy

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sasy/pkg/model"
	"sasy/utils"
	"sort"
)

func AddHandler(arg []string) error {
	if len(arg) == 0 {
		return fmt.Errorf("no file path passed as argument to add")
	}

	// Create index File if it does not exist
	file, err := os.OpenFile(path.Join(utils.SasyPath, "index"), os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("error opening .sasy/index: %v", err)
	}
	defer file.Close()

	index, err := model.ParseIndex(file)
	if err != nil {
		return fmt.Errorf("error parsing .sasy/index: %v", err)
	}

	database, err := model.CreateDatabase(utils.WorkindDir)
	if err != nil {
		return fmt.Errorf("error in creating database: %v", err)
	}

	// for lexical ordering
	sort.Strings(arg)

	for _, entry := range arg {
		filepath.WalkDir(entry, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				fmt.Printf("error traversing path %s: %v\n", path, err)
			}

			if d.IsDir() {
				return nil
			}

			blob := model.CreateBlob(path, d)
			if err := database.Save(blob.Oid, []byte(blob.Content)); err != nil {
				fmt.Printf("error saving blob for file %s: %v\n", path, err)
			}
			if err := index.Modify(path, d, blob.Oid); err != nil {
				fmt.Printf("error inserting entry for file %s in index: %v\n", path, err)
			}

			return nil
		})
	}

	if err := index.Save(); err != nil {
		return fmt.Errorf("error in updating index: %v", err)
	}

	return nil
}
