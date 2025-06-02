package checkers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/pryazhnikov/gofileschecker/internal/scanner"
)

type FilesCheckGroup struct {
	hash  string // common hash for all files in the group
	files []string
}

func (fcg *FilesCheckGroup) addFile(file string) {
	fcg.files = append(fcg.files, file)
}

func (fcg *FilesCheckGroup) HasMultipleFiles() bool {
	return len(fcg.files) > 1
}

func (fcg *FilesCheckGroup) Hash() string {
	return fcg.hash
}

func (fcg *FilesCheckGroup) Files() []string {
	return fcg.files
}

func newFilesCheckGroup(hash string, file string) *FilesCheckGroup {
	return &FilesCheckGroup{
		hash:  hash,
		files: []string{file},
	}
}

type FileChecker struct {
	fileGroups map[string]*FilesCheckGroup
	mu         sync.RWMutex
}

// Ensure FileChecker implements scanner.FileChecker interface
var _ scanner.FileChecker = (*FileChecker)(nil)

func NewFileChecker() *FileChecker {
	return &FileChecker{
		fileGroups: make(map[string]*FilesCheckGroup),
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

	hfr, ok := fc.fileGroups[hash]
	if ok {
		hfr.addFile(path)
	} else {
		fc.fileGroups[hash] = newFilesCheckGroup(hash, path)
	}

	return hash, nil
}

func (fc *FileChecker) GetDuplicatedFileGroups() []*FilesCheckGroup {
	fc.mu.RLock()
	defer fc.mu.RUnlock()

	var result []*FilesCheckGroup
	for _, hfr := range fc.fileGroups {
		if hfr.HasMultipleFiles() {
			result = append(result, hfr)
		}
	}

	return result
}
