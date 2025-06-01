package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type DirectoryScanner struct {
	rootPath string
}

func NewDirectoryScanner(rootPath string) *DirectoryScanner {
	return &DirectoryScanner{rootPath: rootPath}
}

func (ds *DirectoryScanner) Scan() error {
	err := filepath.WalkDir(ds.rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			fmt.Printf("Directory: %s\n", path)
		} else {
			fmt.Printf("File: %s\n", path)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	return nil
}
