package scanner

type ScanSummary struct {
	files       int
	directories int
	errors      int
	skipped     int
}

func (s ScanSummary) Files() int {
	return s.files
}

func (s ScanSummary) Directories() int {
	return s.directories
}

func (s ScanSummary) Errors() int {
	return s.errors
}

func (s ScanSummary) Skipped() int {
	return s.skipped
}

func (s *ScanSummary) AddFile() {
	s.files++
}

func (s *ScanSummary) AddDirectory() {
	s.directories++
}

func (s *ScanSummary) AddError() {
	s.errors++
}

func (s *ScanSummary) AddSkipped() {
	s.skipped++
}
