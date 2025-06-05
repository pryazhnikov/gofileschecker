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
		mu:           sync.RWMutex{},
	}
}

func (ds *DirectoryScanner) Scan(rootPath string) error {
	// Get absolute path to handle different path formats pointing to same directory
	absPath, err := filepath.Abs(rootPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Initialize scanned paths map if not exists
	if ds.scannedPaths == nil {
		return fmt.Errorf("scanned paths is not initialized")
	}

	if ds.isPathScanned(absPath) {
		ds.logger.Info().Msgf("Directory already scanned, skipping: %s", absPath)
		return nil
	}

	ds.logger.Info().Msgf("Starting directory scan: %s", absPath)

	filesCnt := 0
	errCnt := 0
	err = filepath.WalkDir(absPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return ds.processDirectory(path)
		}

		// todo: to implement processing retry later
		err = ds.processFile(path)
		if err != nil {
			errCnt++
		}

		filesCnt++
		if filesCnt%summaryPeriod == 0 {
			ds.logger.Info().Msgf("%d files processed, errors: %d...", filesCnt, errCnt)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan directory: %w", err)
	}

	ds.markPathAsScanned(absPath)
	return nil
}

func (ds *DirectoryScanner) isPathScanned(absPath string) bool {
	ds.mu.RLock()
	defer ds.mu.RUnlock()
	return ds.scannedPaths[absPath]
}

func (ds *DirectoryScanner) markPathAsScanned(absPath string) {
	ds.mu.Lock()
	ds.scannedPaths[absPath] = true
	ds.mu.Unlock()
}

func (ds *DirectoryScanner) processDirectory(path string) error {
	ds.logger.Debug().
		Str("path", path).
		Msg("Directory found, nothing to do here.")
	return nil
}

func (ds *DirectoryScanner) processFile(path string) error {
	ds.logger.Debug().
		Str("path", path).
		Msg("File found, the check is expected")

	checkRes, err := ds.checker.Check(path)
	if err != nil {
		ds.logger.Warn().
			Str("path", path).
			Msgf("Cannot check file: %v", err)
		return err
	}

	ds.logger.Debug().
		Str("path", path).
		Str("hash", checkRes).
		Msg("File was checked")

	return nil
}
