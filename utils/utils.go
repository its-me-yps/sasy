package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/fs"
	"path/filepath"
)

// Returns the list of files and working directory if no error occured
func Ls() ([]string, error) {
	list := []string{}
	if err := filepath.WalkDir(WorkindDir, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return err
		}
		relpath, _ := filepath.Rel(WorkindDir, path)
		for _, entry := range IgnoreDirs {
			if len(relpath) >= len(entry) && relpath[0:len(entry)] == entry {
				return err
			}
		}
		list = append(list, path)
		return err
	}); err != nil {
		return nil, err
	}
	return list, nil
}

// To compress using zlib compression technique
func Compress(content []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(content)
	w.Close()
	return b.Bytes()
}

// Function to calculate OID using SHA1
func CalculateSHA1(content string) string {
	h := sha1.New()
	io.WriteString(h, content)
	return hex.EncodeToString(h.Sum(nil))
}
