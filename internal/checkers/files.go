package checkers

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"

	"github.com/pryazhnikov/gofileschecker/internal/scanner"
)

type FilesCheckGroup struct {
	mu    sync.RWMutex // Protects access to files slice
	hash  string       // common hash for all files in the group
	files []string
}

func (fcg *FilesCheckGroup) HasFile(file string) bool {
	fcg.mu.RLock()
	defer fcg.mu.RUnlock()
	return slices.Contains(fcg.files, file)
}

func (fcg *FilesCheckGroup) addFile(file string) {
	if fcg.HasFile(file) {
		return
	}

	fcg.mu.Lock()
	defer fcg.mu.Unlock()

	// Using slices is not the best option from time complexity perspective,
	// but we do not expected adding multiple files with the same hash many times.
	// So keeping a bit simpler slice-based solution instead of a map-based alternative.
	fcg.files = append(fcg.files, file)
}

func (fcg *FilesCheckGroup) Files() []string {
	fcg.mu.RLock()
	defer fcg.mu.RUnlock()
	return fcg.files
}

func (fcg *FilesCheckGroup) FilesCount() int {
	fcg.mu.RLock()
	defer fcg.mu.RUnlock()
	return len(fcg.files)
}

func (fcg *FilesCheckGroup) HasMultipleFiles() bool {
	return fcg.FilesCount() > 1
}

func (fcg *FilesCheckGroup) CommonPathPrefix() string {
	if len(fcg.files) == 0 {
		return ""
	}

	// For single file return its directory
	if len(fcg.files) == 1 {
		lastSlash := strings.LastIndex(fcg.files[0], "/")
		if lastSlash == -1 {
			return ""
		}

		return fcg.files[0][:lastSlash+1]
	}

	// Start with the first file's path
	prefix := fcg.files[0]
	for _, file := range fcg.files[1:] {
		// Find the common prefix between current prefix and the file
		for i := 0; i < len(prefix) && i < len(file); i++ {
			if prefix[i] != file[i] {
				prefix = prefix[:i]
				break
			}
		}

		// Ensure we don't cut in the middle of a directory name
		lastSlash := strings.LastIndex(prefix, "/")
		if lastSlash == -1 {
			return ""
		}
		prefix = prefix[:lastSlash+1]
	}

	return prefix
}

func (fcg *FilesCheckGroup) Hash() string {
	return fcg.hash
}

func newFilesCheckGroup(hash string, file string) *FilesCheckGroup {
	return &FilesCheckGroup{
		hash:  hash,
		files: []string{file},
		mu:    sync.RWMutex{},
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
