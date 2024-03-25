package git

import (
	"os"
	"path"

	"github.com/sayymeer/sasy/utils"
)

func CommitHandler() error {
	files, _ := utils.Ls()
	wd, _ := os.Getwd()
	objectDir := path.Join(wd, ".sasy", "objects")
	for _, file := range files {
		filePath := path.Join(wd, file)
		blob := CreateObject(filePath, "blob")
		name := blob.getOid()
		dirName := name[0:2]
		fileName := name[2:]
		os.Mkdir(path.Join(objectDir, dirName), os.ModePerm)
		os.WriteFile(path.Join(objectDir, dirName, fileName), blob.getCompressed(), 0644)
	}
	return nil
}
