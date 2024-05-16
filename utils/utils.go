package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// Returns the list of files and working directory if no error occured
func Ls() ([]string, string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, wd, err
	}
	dir, _ := os.Open(wd)
	files, _ := dir.ReadDir(0)
	list := []string{}
	for _, file := range files {
		if !file.IsDir() {
			list = append(list, file.Name())
		}
	}
	return list, wd, nil
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
