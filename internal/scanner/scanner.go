package scanner

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog"
)

const summaryPeriod = 100

type DirectoryScanner struct {
	logger       zerolog.Logger
	checker      FileChecker
	scannedPaths map[string]bool
	summary      *ScanSummary
	mu           sync.RWMutex
}

type FileChecker interface {
	Check(path string) (string, error)
}

func NewDirectoryScanner(logger zerolog.Logger, checker FileChecker) *DirectoryScanner {
	return &DirectoryScanner{
		logger:       logger,
		checker:      checker,
		scannedPaths: make(map[string]bool),
		summary:      &ScanSummary{},
		mu:           sync.RWMutex{},
	}
}

func (ds *DirectoryScanner) Scan(rootPath string) error {
	// Get absolute path to handle different path formats pointing to same directory
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	if ds.scannedPaths == nil {
		return fmt.Errorf("scanned paths field is not initialized")
	}

	if ds.isPathScanned(absPath) {
		ds.logger.Info().Msgf("Directory already scanned, skipping: %s", absPath)
		return nil
	}

	ds.logger.Info().Msgf("Starting directory scan: %s", absPath)

	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return ds.processDirectory(path)
		}

		// todo: to implement processing retry later
		return ds.processFile(path)
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	ds.markPathAsScanned(absPath)
	return nil
}

func (ds *DirectoryScanner) Summary() ScanSummary {
	// Return a copy of the summary to prevent external modifications
	return *ds.summary
}

func (ds *DirectoryScanner) isPathScanned(absPath string) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.scannedPaths[absPath]
}

func (ds *DirectoryScanner) markPathAsScanned(absPath string) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.scannedPaths[absPath] = true
}

func (ds *DirectoryScanner) processDirectory(path string) error {
	ds.logger.Debug().
		Str("path", path).
		Msg("Directory found, nothing to do here.")
	ds.summary.AddDirectory()
	return nil
}

func (ds *DirectoryScanner) processFile(path string) error {
	ds.logger.Debug().
		Str("path", path).
		Msg("File found, the check is expected")

	ds.summary.AddFile()
	checkRes, err := ds.checker.Check(path)
	if err != nil {
		ds.logger.Warn().
			Str("path", path).
			Msgf("Cannot check file: %v", err)
		ds.summary.AddError()
		return err
	}

	ds.logger.Debug().
		Str("path", path).
		Str("hash", checkRes).
		Msg("File was checked")

	if ds.summary.Files()%summaryPeriod == 0 {
		ds.logger.Info().Msgf(
			"%d files processed, errors: %d...",
			ds.summary.Files(),
			ds.summary.Errors(),
		)
	}

	return nil
}
