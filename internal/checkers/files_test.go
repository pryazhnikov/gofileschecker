package checkers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFilesCheckGroup(t *testing.T) {
	// Arrange
	expectedHash := "test-hash-123"
	expectedFile := "initial-file.txt"

	// Act
	fcg := newFilesCheckGroup(expectedHash, expectedFile)

	// Assert
	assert.NotNil(t, fcg, "newFilesCheckGroup should not return nil")
	assert.Equal(t, expectedHash, fcg.Hash(), "Hash should match the input value")

	files := fcg.Files()
	assert.Len(t, files, 1, "Should have exactly one file")
	assert.Equal(t, expectedFile, files[0], "File should match the input value")

	assert.False(t, fcg.HasMultipleFiles(), "HasMultipleFiles should be false for single file")
}

func TestFilesCheckGroup_AddFile(t *testing.T) {
	// Arrange
	initialHash := "test-hash"
	initialFile := "file1.txt"
	fcg := newFilesCheckGroup(initialHash, initialFile)
	assert.False(t, fcg.HasMultipleFiles(), "HasMultipleFiles should be false initially")

	// Act
	fcg.addFile("file2.txt")

	// Assert
	assert.Equal(t, initialHash, fcg.Hash(), "Hash should remain unchanged")

	files := fcg.Files()
	assert.Len(t, files, 2, "Should have exactly two files")
	assert.Equal(t, initialFile, files[0], "First file should match initial file")
	assert.Equal(t, "file2.txt", files[1], "Second file should match added file")

	assert.True(t, fcg.HasMultipleFiles(), "HasMultipleFiles should be true after adding file")
}

func TestFilesCheckGroup_HasFile(t *testing.T) {
	// Arrange
	initialHash := "test-hash"
	initialFile := "file1.txt"
	fcg := newFilesCheckGroup(initialHash, initialFile)

	// Act & Assert
	assert.True(t, fcg.HasFile(initialFile), "Should return true for existing file")
	assert.False(t, fcg.HasFile("non-existent.txt"), "Should return false for non-existent file")

	// Add another file and verify
	extraFile := "file2.txt"
	assert.NotEqual(t, initialFile, extraFile, "Precondition failed: file names should be different")

	fcg.addFile(extraFile)
	assert.True(t, fcg.HasFile(extraFile), "Should return true for newly added file")
}

func TestFilesCheckGroup_AddTheSameFile(t *testing.T) {
	// Arrange
	initialHash := "test-hash"
	initialFile := "file1.txt"
	fcg := newFilesCheckGroup(initialHash, initialFile)
	assert.False(t, fcg.HasMultipleFiles(), "HasMultipleFiles should be false initially")
	assert.Equal(t, 1, fcg.FilesCount(), "The only file is expected")

	// Act
	fcg.addFile(initialFile)

	// Assert
	assert.Equal(t, 1, fcg.FilesCount(), "The only file is expected after adding the same one again")
}

// The changes in the slice returned by Files() should not change the state of FilesCheckGroup
func TestFilesCheckGroup_FilesShouldBeImmutable(t *testing.T) {
	initialHash := "test-hash"
	initialFile := "file1.txt"         // Should be placed into check group
	modifedFile := "modified-file.txt" // Should NOT be placed into the check group

	assert.NotEqual(t, initialFile, modifedFile, "Precondition failed: files should be different")

	// Arrange
	fcg := newFilesCheckGroup(initialHash, initialFile)

	assert.True(t, fcg.HasFile(initialFile), "Group should contain the initial file")
	assert.False(t, fcg.HasFile(modifedFile), "Group should not contain the modified file")

	// Act
	files := fcg.Files()
	originalLen := len(files)

	// Modify the returned slice by changing existing value
	// We should modify the existing values here, append usage is not allowed
	// (because it can change the resurned slice pointer)
	files[0] = modifedFile

	// Assert
	assert.Equal(t, originalLen, fcg.FilesCount(), "FilesCount should remain unchanged")
	assert.True(t, fcg.HasFile(initialFile), "Group should contain the initial file")
	assert.False(t, fcg.HasFile(modifedFile), "Group should not contain the modified file")

	// Verify original files are intact
	assert.Equal(t, []string{initialFile}, fcg.Files(), "Internal files slice should be unchanged")
}

func TestFilesCheckGroup_CommonPathPrefix(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected string
	}{
		{
			name:     "single file without directory",
			files:    []string{"file.txt"},
			expected: "",
		},
		{
			name:     "single file with directory",
			files:    []string{"/path/to/file.txt"},
			expected: "/path/to/",
		},
		{
			name:     "two files in same directory",
			files:    []string{"/path/to/file1.txt", "/path/to/file2.txt"},
			expected: "/path/to/",
		},
		{
			name:     "two files in different directories",
			files:    []string{"/path/to/file1.txt", "/path/different/file2.txt"},
			expected: "/path/",
		},
		{
			name:     "three files with common prefix",
			files:    []string{"/path/to/dir1/file1.txt", "/path/to/dir2/file2.txt", "/path/to/dir3/file3.txt"},
			expected: "/path/to/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a FilesCheckGroup with a dummy hash
			fcg := newFilesCheckGroup("dummy-hash", "")
			fcg.files = tt.files // Override the files directly for testing

			result := fcg.CommonPathPrefix()
			assert.Equal(t, tt.expected, result, "CommonPathPrefix returned unexpected result")
		})
	}
}
