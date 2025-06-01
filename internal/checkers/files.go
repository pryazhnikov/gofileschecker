package checkers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pryazhnikov/gofileschecker/internal/scanner"
)

type FileChecker struct {
	fileHashes map[string]string
	mu         sync.RWMutex
}

// Ensure FileChecker implements scanner.FileChecker interface
var _ scanner.FileChecker = (*FileChecker)(nil)

func NewFileChecker() *FileChecker {
	return &FileChecker{
		fileHashes: make(map[string]string),
	}
}

func (fc *FileChecker) Check(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", fmt.Errorf("failed to calculate hash: %w", err)
	}

	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.fileHashes[path] = hash

	return hash, nil
}
