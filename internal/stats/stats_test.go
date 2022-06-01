package stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStats(t *testing.T) {
	s := SyncStats{Stats: Stats{Miss: 100}}
	s.Reset()
	assert.Equal(t, s.Miss, 0)

	s.IncRead(true, 100)
	s.IncRead(false, 100)
	s.IncWrite(1000)
	s.IncWrite(100)
	s.IncDelete()
	s.IncClear()
	s.ErrRead()
	s.ErrWrite()
	s.ErrDelete()
	s.ErrClear()

	assert.Equal(t, s.Miss, 1)
	assert.Equal(t, s.Hits, 1)
	assert.Equal(t, s.ReadBytes, 100)
	assert.Equal(t, s.WriteBytes, 1100)
	assert.Equal(t, s.ReadCount, 2)
	assert.Equal(t, s.WriteCount, 2)
	assert.Equal(t, s.DeleteCount, 1)
	assert.Equal(t, s.ErrReadCount, 1)
	assert.Equal(t, s.ErrWriteCount, 1)
	assert.Equal(t, s.ErrDeleteCount, 1)
	assert.Equal(t, s.ErrDeleteCount, 1)
}
