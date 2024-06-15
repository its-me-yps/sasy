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
	fileSize := fileInfo.Size()
	if fileSize == 0 {
		return Index{
			Header:  Header{Signature: "DIRC", version: string([]byte{0, 0, 0, 2}), numberOfEntries: 0},
			Entries: make(map[string]Entry),
		}, nil
	}

	// checking integrity of index file
	// seek to the file 20 bytes before end
	// last 20 bytes of index is checksum for rest of the index
	file.Seek(-20, 2)
	checksum := make([]byte, 20)
	file.Read(checksum)
	
	bufferSize := fileSize - 20
	if bufferSize < 0 {
		return Index{}, fmt.Errorf("index file is corrupted")
	}
	// WARNING: buffer can be very large
	// TODO: Add chucks of file data to hasher (while parsing) and do integrity check at last
	buffer := make([]byte, bufferSize)
	file.ReadAt(buffer, 0)
	if hex.EncodeToString(checksum) != utils.CalculateSHA1(string(buffer)) {
		return Index{}, fmt.Errorf("index file is corrupted")
	}

	// Reset the file offset
	file.Seek(0, 0)

	// This variable will keep track of number of bytes read from file
	var bytesRead uint64

	var h Header
	hSlice := make([]byte, 12)
	file.Read(hSlice)
	bytesRead += 12

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
		bytesRead += 50

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
		bytesRead += uint64(e.entryPathSize)

		// TODO: Implement the case when entryPathSize overflows
		// In that case chunks of 8 bytes should be read until a null byte is encountered

		// Skip Reading the padded null bytes at the end of each Entry
		offset := 8 - bytesRead%8
		file.Seek(int64(offset), 1)
		bytesRead += offset

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

	file, _ := os.OpenFile(path.Join(utils.SasyPath, "index"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	defer file.Close()

	// This buffer will store the whole content and will be used to hash the whole file and write checksum for the file at the end
	// WARNING: buffer can become large, a better way to do this would be to write chunks of data in buffer and update hasher with that chunk. Then clear the buffer for next chunk of data.
	var buffer bytes.Buffer

	// Writing header
	headerBytes := make([]byte, 12)
	copy(headerBytes[0:4], []byte(index.Header.Signature))
	copy(headerBytes[4:8], []byte(index.Header.version))
	binary.BigEndian.PutUint32(headerBytes[8:12], index.Header.numberOfEntries)
	buffer.Write(headerBytes)
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
		buffer.Write(metadataBytes)
		n, err := file.Write(metadataBytes)
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		buffer.Write([]byte(entry.entryId))
		n, err = file.Write([]byte(entry.entryId))
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		sizeBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(sizeBytes, uint16(entry.entryPathSize))
		buffer.Write(sizeBytes)
		n, err = file.Write(sizeBytes)
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		buffer.Write([]byte(key))
		n, err = file.Write([]byte(key))
		if err != nil {
			return err
		}
		indexSize += uint64(n)

		// padding null bytes to ensure that index size remains multiple of 8 at the end of each entry (usecase in parsing)
		// if index size is already multiple of 8 at the end of an entry then we pad 8 null bytes
		padding := 8 - indexSize%8
		padBytes := make([]byte, padding)
		buffer.Write(padBytes)
		n, err = file.Write(padBytes)
		if err != nil {
			return err
		}
		indexSize += uint64(n)

	}

	checksum := utils.CalculateSHA1(string(buffer.Bytes()))
	checksumBytes, _ := hex.DecodeString(checksum)

	_, err = file.Write(checksumBytes)
	if err != nil {
		return err
	}

	return nil
}
