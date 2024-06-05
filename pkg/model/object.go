package model

import (
	"fmt"
)

// blob, tree, commit implement this interface by defining setContent() method
type Object interface {
	setContent()
}

// recursively saves objects to database 
func SaveObjectsToDatabase(entries []Object, database *Database) error {
	for _, entry := range entries {
		switch entry.(type) {
		case *Blob:
			blob, _ := entry.(*Blob)
			if err := database.Save(blob.Oid, []byte(blob.Content)); err != nil {
				return fmt.Errorf("error saving object %s of type \"blob\" in db: %v", blob.Name, err)
			}
		case *Tree:
			tree, _ := entry.(*Tree)
			if err := database.Save(tree.Oid, []byte(tree.Content)); err != nil {
				return fmt.Errorf("error saving object %s of type \"tree\" in db: %v", tree.Name, err)
			}

			if err := SaveObjectsToDatabase(tree.Entries, database); err != nil {
				return err
			}
		}
	}
	return nil
}
