// +build darwin

package utils

import (
	"os"
	"syscall"
	"time"
)

func getUnixMetadata(fileInfo os.FileInfo) (Metadata, error) {
	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		return Metadata{}, syscall.EINVAL
	}

	return Metadata{
		Ctime: time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec)),
		Mtime: time.Unix(int64(stat.Mtimespec.Sec), int64(stat.Mtimespec.Nsec)),
		Mode:  uint16(stat.Mode),
		Size:  stat.Size,
	}, nil
}
