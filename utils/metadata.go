package utils

import (
	"bytes"
	"encoding/binary"
	"time"
)

type Metadata struct {
	Ctime time.Time
	Mtime time.Time
	Mode  uint32
	Size  int64
}

func MetadataToBytes(metadata Metadata) []byte {
	var buf bytes.Buffer

	ctimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(ctimeBytes, uint64(metadata.Ctime.UnixNano()))
	buf.Write(ctimeBytes)

	mtimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(mtimeBytes, uint64(metadata.Mtime.UnixNano()))
	buf.Write(mtimeBytes)

	modeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(modeBytes, metadata.Mode)
	buf.Write(modeBytes)

	sizeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBytes, uint64(metadata.Size))
	buf.Write(sizeBytes)

	return buf.Bytes()
}

func BytesToMetadata(data []byte) (Metadata, error) {
	var ctime int64
	var mtime int64
	var mode uint32
	var size int64

	buf := bytes.NewReader(data)

	if err := binary.Read(buf, binary.LittleEndian, &ctime); err != nil {
		return Metadata{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &mtime); err != nil {
		return Metadata{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &mode); err != nil {
		return Metadata{}, err
	}

	if err := binary.Read(buf, binary.LittleEndian, &size); err != nil {
		return Metadata{}, err
	}

	return Metadata{
		Ctime: time.Unix(0, ctime),
		Mtime: time.Unix(0, mtime),
		Mode:  mode,
		Size:  size,
	}, nil
}
