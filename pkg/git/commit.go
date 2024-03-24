package git

import (
	"bytes"
	"compress/zlib"
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
	objectDir := path.Join(wd, ".sasy", "objects")
	for _, file := range files {
		filePath := path.Join(wd, file)
		byteContent, _ := os.ReadFile(filePath)
		content := fmt.Sprintf("%s %d\x00%s", "blob", len(byteContent), string(byteContent))
		h := sha1.New()
		io.WriteString(h, content)
		name := fmt.Sprintf("%x", h.Sum(nil))
		dirName := name[0:2]
		fileName := name[2:]
		os.Mkdir(path.Join(objectDir, dirName), os.ModePerm)
		var b bytes.Buffer
		w := zlib.NewWriter(&b)
		w.Write([]byte(content))
		w.Close()
		os.WriteFile(path.Join(objectDir, dirName, fileName), b.Bytes(), 0644)
	}
	return nil
}
