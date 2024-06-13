package utils

import (
	"os"
	"runtime"
	"syscall"
	"time"
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

func getUnixMetadata(fileInfo os.FileInfo) (Metadata, error) {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return Metadata{}, syscall.EINVAL
	}

	return Metadata{
		Ctime: time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec)),
		Mtime: time.Unix(int64(stat.Mtimespec.Sec), int64(stat.Mtimespec.Nsec)),
		Mode:  stat.Mode,
		Size:  stat.Size,
	}, nil
}

// func getWindowsMetadata(fileInfo os.FileInfo) (Metadata, error) {
//
// }
