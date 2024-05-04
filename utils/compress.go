package utils

import (
	"bytes"
	"compress/zlib"
)

func Compress(content []byte) []byte {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write(content)
	w.Close()
	return b.Bytes()
}
