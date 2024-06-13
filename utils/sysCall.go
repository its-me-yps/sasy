package utils

import (
	"os"
	"runtime"
)

func GetFileMetadata(d os.DirEntry) (Metadata, error) {
	fileInfo, err := d.Info()
	if err != nil {
		return Metadata{}, err
	}

	var metadata Metadata

	switch runtime.GOOS {
	// TODO: Implement syscall for windows os to get file metadata
	
	// case "windows":
	// metadata, err = getWindowsMetadata(fileInfo)
	default:
		metadata, err = getUnixMetadata(fileInfo)
	}

	if err != nil {
		return Metadata{}, err
	}

	return metadata, nil
}

// func getWindowsMetadata(fileInfo os.FileInfo) (Metadata, error) {
//
// }
