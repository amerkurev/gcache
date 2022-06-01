package stats

import "sync"

// Stats collects significant metrics about cache operations.
type Stats struct {
	Hits           int
	Miss           int
	ReadBytes      int
	WriteBytes     int
	ReadCount      int
	WriteCount     int
	DeleteCount    int
	ClearCount     int
	ErrReadCount   int
	ErrWriteCount  int
	ErrDeleteCount int
	ErrClearCount  int
}

// SyncStats implements concurrency-safe methods above Stats data fields.
type SyncStats struct {
	mx sync.Mutex
	Stats
}

// Snapshot returns copy of collected metrics.
func (s *SyncStats) Snapshot() Stats {
	s.mx.Lock()
	defer s.mx.Unlock()

	return s.Stats
}

// Reset sets all metric fields to zero-value.
func (s *SyncStats) Reset() {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.Hits = 0
	s.Miss = 0
	s.ReadBytes = 0
	s.WriteBytes = 0
	s.ReadCount = 0
	s.WriteCount = 0
	s.DeleteCount = 0
	s.ClearCount = 0
	s.ErrReadCount = 0
	s.ErrWriteCount = 0
	s.ErrDeleteCount = 0
	s.ErrClearCount = 0
}

// IncRead increments metrics of read operation.
func (s *SyncStats) IncRead(hits bool, n int) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.ReadCount++

	if hits {
		s.Hits++
		s.ReadBytes += n
	} else {
		s.Miss++
	}
}

// IncWrite increments metrics of write operation.
func (s *SyncStats) IncWrite(n int) {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.WriteCount++
	s.WriteBytes += n
}

// IncDelete increments metrics of delete operation.
func (s *SyncStats) IncDelete() {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.DeleteCount++
}

// IncClear increments metrics of clear operation.
func (s *SyncStats) IncClear() {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.ClearCount++
}

// ErrRead increments the read error counter.
func (s *SyncStats) ErrRead() {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.ErrReadCount++
}

// ErrWrite increments the write error counter.
func (s *SyncStats) ErrWrite() {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.ErrWriteCount++
}

// ErrDelete increments the delete error counter.
func (s *SyncStats) ErrDelete() {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.ErrDeleteCount++
}

// ErrClear increments the clear error counter.
func (s *SyncStats) ErrClear() {
	s.mx.Lock()
	defer s.mx.Unlock()

	s.ErrClearCount++
}
