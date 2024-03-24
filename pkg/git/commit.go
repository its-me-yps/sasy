package git

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/sayymeer/sasy/utils"
)

func CommitHandler() error {
	files, _ := utils.Ls()
	wd, _ := os.Getwd()
	for _, file := range files {
		filePath := path.Join(wd, file)
		byteContent, _ := os.ReadFile(filePath)
		content := fmt.Sprintf("%s %d/0%s", "blob", len(byteContent), string(byteContent))
		h := sha1.New()
		io.WriteString(h, content)
		fmt.Printf("%x\n", h.Sum(nil))
	}
	return nil
}
