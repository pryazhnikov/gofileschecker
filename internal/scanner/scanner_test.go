package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockFileChecker struct {
	mu            sync.Mutex
	checkCount    int
	checkDuration time.Duration
	shouldError   bool
}

func (m *mockFileChecker) Check(path string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	
	m.checkCount++
	
	if m.checkDuration > 0 {
		time.Sleep(m.checkDuration)
	}
	
	if m.shouldError {
		return "", fmt.Errorf("mock error for %s", path)
	}
	
	return fmt.Sprintf("hash-%s", filepath.Base(path)), nil
}

func (m *mockFileChecker) getCheckCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.checkCount
}

func TestNewDirectoryScanner(t *testing.T) {
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)
	checker := &mockFileChecker{}

	t.Run("creates scanner successfully", func(t *testing.T) {
		scanner := NewDirectoryScanner(logger, checker)
		assert.NotNil(t, scanner)
		assert.NotNil(t, scanner.summary)
		assert.NotNil(t, scanner.scannedPaths)
	})
}

func TestDirectoryScanner_BasicFunctionality(t *testing.T) {
	tempDir := t.TempDir()
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt", "file5.txt"}
	for _, filename := range testFiles {
		err := os.WriteFile(filepath.Join(tempDir, filename), []byte("test content"), 0644)
		require.NoError(t, err)
	}

	t.Run("processes all files successfully", func(t *testing.T) {
		checker := &mockFileChecker{}
		scanner := NewDirectoryScanner(logger, checker)

		err := scanner.Scan(tempDir)

		require.NoError(t, err)
		assert.Equal(t, len(testFiles), checker.getCheckCount())
		summary := scanner.Summary()
		assert.Equal(t, len(testFiles), summary.Files())
		assert.Equal(t, 1, summary.Directories())
		assert.Equal(t, 0, summary.Errors())
	})

	t.Run("handles file check errors", func(t *testing.T) {
		checker := &mockFileChecker{shouldError: true}
		scanner := NewDirectoryScanner(logger, checker)

		err := scanner.Scan(tempDir)
		require.Error(t, err) // Scan should error on first file check failure

		// Only one file should be processed before error stops the scan
		assert.Equal(t, 1, checker.getCheckCount())
		summary := scanner.Summary()
		assert.Equal(t, 1, summary.Files())
		assert.Equal(t, 1, summary.Errors())
	})
}

func TestDirectoryScanner_ScanSummary(t *testing.T) {
	tempDir := t.TempDir()
	logger := zerolog.New(os.Stdout).Level(zerolog.Disabled)

	// Create test files in subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	require.NoError(t, err)

	numFiles := 10
	for i := 0; i < numFiles; i++ {
		filename := fmt.Sprintf("file%d.txt", i)
		err := os.WriteFile(filepath.Join(tempDir, filename), []byte("test content"), 0644)
		require.NoError(t, err)
	}

	// Create one file in subdirectory
	err = os.WriteFile(filepath.Join(subDir, "subfile.txt"), []byte("sub content"), 0644)
	require.NoError(t, err)

	checker := &mockFileChecker{}
	scanner := NewDirectoryScanner(logger, checker)

	err = scanner.Scan(tempDir)
	require.NoError(t, err)

	// Verify scan summary
	assert.Equal(t, numFiles+1, checker.getCheckCount()) // +1 for subfile
	summary := scanner.Summary()
	assert.Equal(t, numFiles+1, summary.Files())
	assert.Equal(t, 2, summary.Directories()) // tempDir + subDir
	assert.Equal(t, 0, summary.Errors())
}