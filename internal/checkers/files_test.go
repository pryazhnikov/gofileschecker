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
