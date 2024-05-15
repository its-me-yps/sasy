package utils

import "os"

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
