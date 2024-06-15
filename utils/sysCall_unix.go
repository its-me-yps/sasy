// +build linux

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
		Ctime: time.Unix(int64(stat.Ctim.Sec), int64(stat.Ctim.Nsec)),
		Mtime: time.Unix(int64(stat.Mtim.Sec), int64(stat.Mtim.Nsec)),
		Mode:  uint32(stat.Mode),
		Size:  stat.Size,
	}, nil
}
