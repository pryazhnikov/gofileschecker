package scanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScanSummary_InitialState(t *testing.T) {
	// Arrange & Act
	summary := &ScanSummary{}

	// Assert
	assert.Equal(t, 0, summary.Files(), "Files should be 0 initially")
	assert.Equal(t, 0, summary.Directories(), "Directories should be 0 initially")
	assert.Equal(t, 0, summary.Errors(), "Errors should be 0 initially")
	assert.Equal(t, 0, summary.Skipped(), "Skipped should be 0 initially")
}

func TestScanSummary_AddFile(t *testing.T) {
	// Arrange
	summary := &ScanSummary{}

	// Act
	summary.AddFile()

	// Assert
	assert.Equal(t, 1, summary.Files(), "Files should be 1 after AddFile()")
	assert.Equal(t, 0, summary.Directories(), "Directories should remain 0")
	assert.Equal(t, 0, summary.Errors(), "Errors should remain 0")
	assert.Equal(t, 0, summary.Skipped(), "Skipped should remain 0")
}

func TestScanSummary_AddDirectory(t *testing.T) {
	// Arrange
	summary := &ScanSummary{}

	// Act
	summary.AddDirectory()

	// Assert
	assert.Equal(t, 0, summary.Files(), "Files should remain 0")
	assert.Equal(t, 1, summary.Directories(), "Directories should be 1 after AddDirectory()")
	assert.Equal(t, 0, summary.Errors(), "Errors should remain 0")
	assert.Equal(t, 0, summary.Skipped(), "Skipped should remain 0")
}

func TestScanSummary_AddError(t *testing.T) {
	// Arrange
	summary := &ScanSummary{}

	// Act
	summary.AddError()

	// Assert
	assert.Equal(t, 0, summary.Files(), "Files should remain 0")
	assert.Equal(t, 0, summary.Directories(), "Directories should remain 0")
	assert.Equal(t, 1, summary.Errors(), "Errors should be 1 after AddError()")
	assert.Equal(t, 0, summary.Skipped(), "Skipped should remain 0")
}

func TestScanSummary_AddSkipped(t *testing.T) {
	// Arrange
	summary := &ScanSummary{}

	// Act
	summary.AddSkipped()

	// Assert
	assert.Equal(t, 0, summary.Files(), "Files should remain 0")
	assert.Equal(t, 0, summary.Directories(), "Directories should remain 0")
	assert.Equal(t, 0, summary.Errors(), "Errors should remain 0")
	assert.Equal(t, 1, summary.Skipped(), "Skipped should be 1 after AddSkipped()")
}

func TestScanSummary_MultipleIncrements(t *testing.T) {
	// Arrange
	summary := &ScanSummary{}

	// Act
	summary.AddFile()
	summary.AddFile()
	summary.AddFile()
	summary.AddDirectory()
	summary.AddDirectory()
	summary.AddError()
	summary.AddSkipped()
	summary.AddSkipped()
	summary.AddSkipped()
	summary.AddSkipped()

	// Assert
	assert.Equal(t, 3, summary.Files(), "Files should be 3 after 3 AddFile() calls")
	assert.Equal(t, 2, summary.Directories(), "Directories should be 2 after 2 AddDirectory() calls")
	assert.Equal(t, 1, summary.Errors(), "Errors should be 1 after 1 AddError() call")
	assert.Equal(t, 4, summary.Skipped(), "Skipped should be 4 after 4 AddSkipped() calls")
}