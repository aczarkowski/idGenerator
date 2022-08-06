package epoch

import "time"

// TimeProvider Implements the TimeProvider interface
type timeProvider struct {
	epochOffset int64
}

// New TimeProvider based on provided epoch
func New(epochOffset int64) *timeProvider {
	return &timeProvider{
		epochOffset: epochOffset,
	}
}

// GetTimeStamp Returns the current time in milliseconds
func (t *timeProvider) GetTimeStamp() int64 {
	return time.Now().UTC().UnixMilli() - t.epochOffset
}
