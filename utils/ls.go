package utils

import "os"

func Ls() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir, _ := os.Open(wd)
	files, _ := dir.ReadDir(0)
	list := []string{}
	for _, file := range files {
		if !file.IsDir() {
			list = append(list, file.Name())
		}
	}
	return list, nil
}
