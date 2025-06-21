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

type ScanSummaryCollector struct {
	data ScanSummaryStats
	mu   sync.RWMutex
}

func (s *ScanSummaryCollector) Files() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.files
}

func (s *ScanSummaryCollector) Directories() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.directories
}

func (s *ScanSummaryCollector) Errors() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.errors
}

func (s *ScanSummaryCollector) Skipped() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data.skipped
}

func (s *ScanSummaryCollector) AddFile() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.files++
}

func (s *ScanSummaryCollector) AddDirectory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.directories++
}

func (s *ScanSummaryCollector) AddError() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.errors++
}

func (s *ScanSummaryCollector) AddSkipped() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data.skipped++
}

func (s *ScanSummaryCollector) Stats() ScanSummaryStats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.data
}
