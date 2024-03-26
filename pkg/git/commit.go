package git

import (
	"os"
	"path"

	"github.com/sayymeer/sasy/utils"
)

func CommitHandler() error {
	files, _ := utils.Ls()
	wd, _ := os.Getwd()
	database, err := CreateDatabase(wd)
	if err != nil {
		return err
	}
	for _, file := range files {
		filePath := path.Join(wd, file)
		blob := CreateObject(filePath, "blob")
		if err := database.SaveObject(blob); err != nil {
			return err
		}
	}
	return nil
}
