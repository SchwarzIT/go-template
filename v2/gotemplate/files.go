package gotemplate

import (
	"os"
	"path/filepath"
)

// FindFiles returns a slice of strings representing the file paths of all files in the
// directory tree rooted at the specified path. Directories are skipped.
// If any error occurs while traversing the directory tree, it is returned.
func FindFiles(path string) ([]string, error) {
	var files []string

	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		files = append(files, filePath)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}
