package epoch

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	offset := int64(1420070400000)
	provider := New(offset)
	
	if provider == nil {
		t.Error("Expected provider to be created, got nil")
	}
	
	if provider.epochOffset != offset {
		t.Errorf("Expected offset %d, got %d", offset, provider.epochOffset)
	}
}

func TestGetTimeStamp(t *testing.T) {
	offset := int64(1420070400000)
	provider := New(offset)
	
	timestamp1 := provider.GetTimeStamp()
	time.Sleep(time.Millisecond * 2) // Small delay to ensure different timestamps
	timestamp2 := provider.GetTimeStamp()
	
	if timestamp1 <= 0 {
		t.Errorf("Expected positive timestamp, got %d", timestamp1)
	}
	
	if timestamp2 <= timestamp1 {
		t.Errorf("Expected timestamp2 (%d) to be greater than timestamp1 (%d)", timestamp2, timestamp1)
	}
}

func TestGetTimeStamp_WithZeroOffset(t *testing.T) {
	provider := New(0)
	timestamp := provider.GetTimeStamp()
	
	// Should be close to current Unix timestamp in milliseconds
	currentUnixMilli := time.Now().UTC().UnixMilli()
	if timestamp < currentUnixMilli-1000 || timestamp > currentUnixMilli+1000 {
		t.Errorf("Expected timestamp to be close to current time, got %d, expected around %d", timestamp, currentUnixMilli)
	}
}

func TestGetTimeStamp_WithOffset(t *testing.T) {
	offset := int64(1420070400000) // Jan 1, 2015 00:00:00 UTC in milliseconds
	provider := New(offset)
	timestamp := provider.GetTimeStamp()
	
	// Should be current time minus offset
	expectedApprox := time.Now().UTC().UnixMilli() - offset
	if timestamp < expectedApprox-1000 || timestamp > expectedApprox+1000 {
		t.Errorf("Expected timestamp to be around %d, got %d", expectedApprox, timestamp)
	}
}
