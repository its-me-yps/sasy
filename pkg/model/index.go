package model

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"sasy/utils"
	"sort"
)

type Header struct {
	Signature       string
	version         string
	numberOfEntries uint32
}

type Entry struct {
	Metadata      utils.Metadata
	entryId       string
	entryPathSize uint16
}

type Index struct {
	Header  Header
	Entries map[string]Entry
}

func ParseIndex(file *os.File) (Index, error) {

	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		return Index{
			Header:  Header{Signature: "DIRC", version: string([]byte{0, 0, 0, 2}), numberOfEntries: 0},
			Entries: make(map[string]Entry),
		}, nil
	}

	var h Header
	hSlice := make([]byte, 12)
	file.Read(hSlice)

	h.Signature = string(hSlice[0:4])
	h.version = string(hSlice[4:8])
	// Reading number of entries (4 byte) and converting it to uint32
	hBuf := bytes.NewBuffer(hSlice[8:12])
	binary.Read(hBuf, binary.BigEndian, &h.numberOfEntries)

	index := Index{
		Header:  h,
		Entries: make(map[string]Entry),
	}

	for entryNumber := uint32(0); entryNumber < h.numberOfEntries; entryNumber++ {
		var e Entry
		metadataSlice := make([]byte, 28)
		entryIdSlice := make([]byte, 20)
		sizeSlice := make([]byte, 2)

		file.Read(metadataSlice)
		file.Read(entryIdSlice)
		file.Read(sizeSlice)

		meta, err := utils.BytesToMetadata(metadataSlice)
		if err != nil {
			return Index{}, fmt.Errorf("wrong formatting of metadata for entry %d: %v", entryNumber, err)
		}
		e.Metadata = meta
		e.entryId = string(entryIdSlice)

		sizeBuf := bytes.NewBuffer(sizeSlice)
		binary.Read(sizeBuf, binary.BigEndian, &e.entryPathSize)

		nameSlice := make([]byte, e.entryPathSize)
		file.Read(nameSlice)

		Name := string(nameSlice)

		index.Entries[Name] = e
	}

	return index, nil
}

func (index *Index) Modify(path string, d os.DirEntry, Oid string) error {
	meta, err := utils.GetFileMetadata(d)
	if err != nil {
		return fmt.Errorf("stat() sys call error: %v", err)
	}

	entryId, _ := hex.DecodeString(Oid)
	index.Entries[path] = Entry{Metadata: meta, entryId: string(entryId), entryPathSize: uint16(len(path))}

	// Updating Number of entries in header, incase a new entry was added
	index.Header.numberOfEntries = uint32(len(index.Entries))

	return nil
}

func (index *Index) Save() error {

	file, err := os.OpenFile(path.Join(utils.SasyPath, "index"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	// Writing header
	headerBytes := make([]byte, 12)
	copy(headerBytes[0:4], []byte(index.Header.Signature))
	copy(headerBytes[4:8], []byte(index.Header.version))
	binary.BigEndian.PutUint32(headerBytes[8:12], index.Header.numberOfEntries)
	n, err := file.Write(headerBytes)
	if err != nil {
		return err
	}

	var indexSize uint64
	indexSize += uint64(n)

	// Extract keys and sort them
	keys := make([]string, 0, len(index.Entries))
	for key := range index.Entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Writing Entries
	for _, key := range keys {
		entry := index.Entries[key]
		
		metadataBytes := utils.MetadataToBytes(entry.Metadata)
		n, err := file.Write(metadataBytes)
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		n, err = file.Write([]byte(entry.entryId))
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		sizeBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(sizeBytes, uint16(entry.entryPathSize))
		n, err = file.Write(sizeBytes)
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		n, err = file.Write([]byte(key))
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		// if indexSize%8 {
		// TODO: Pad null bytes at the end of each entry so that indexSize is multiple of 8
		// Doing so will help in parsing
		// }
	}

	// TODO: Write checksum for whole index file at the end
	return nil
}
