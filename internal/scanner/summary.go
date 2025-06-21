package scanner

import "sync"

type ScanSummaryStats struct {
	files       int
	directories int
	errors      int
	skipped     int
}

func (s ScanSummaryStats) Files() int {
	return s.files
}

func (s ScanSummaryStats) Directories() int {
	return s.directories
}

func (s ScanSummaryStats) Errors() int {
	return s.errors
}

func (s ScanSummaryStats) Skipped() int {
	return s.skipped
}

type ScanSummary struct {
	data ScanSummaryStats
	mu   sync.RWMutex
}

func (s *ScanSummary) Files() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.files
}

func (s *ScanSummary) Directories() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.directories
}

func (s *ScanSummary) Errors() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.errors
}

func (s *ScanSummary) Skipped() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.skipped
}

func (s *ScanSummary) AddFile() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.files++
}

func (s *ScanSummary) AddDirectory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.directories++
}

func (s *ScanSummary) AddError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.errors++
}

func (s *ScanSummary) AddSkipped() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.skipped++
}

func (s *ScanSummary) Stats() ScanSummaryStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}
