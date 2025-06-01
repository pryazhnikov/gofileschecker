package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type DirectoryScanner struct {
	rootPath string
	checker  FileChecker
}

type FileChecker interface {
	Check(path string) (string, error)
}

func NewDirectoryScanner(rootPath string, checker FileChecker) *DirectoryScanner {
	return &DirectoryScanner{rootPath: rootPath, checker: checker}
}

func (ds *DirectoryScanner) Scan() error {
	err := filepath.WalkDir(ds.rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Nothing to do here
			fmt.Printf("Directory: %s\n", path)
		} else {
			checkRes, err := ds.checker.Check(path)
			if err != nil {
				fmt.Printf("File: %s - cannot check\n", path)
			} else {
				fmt.Printf("File: %s (result: %s)\n", path, checkRes)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	return nil
}
