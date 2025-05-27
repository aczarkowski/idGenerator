package generator

import (
	"testing"
	"uidGenerator/timeprovider/epoch"
)

func TestGenerateID_SingleID(t *testing.T) {
	provider := epoch.New(1420070400000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}

	ids, err := worker.GenerateID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ids) != 1 {
		t.Errorf("Expected 1 ID, got %d", len(ids))
	}

	if ids[0] <= 0 {
		t.Errorf("Expected positive ID, got %d", ids[0])
	}
}

func TestGenerateID_MultipleIDs(t *testing.T) {
	provider := epoch.New(1420070400000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}

	numberOfIds := 10
	ids, err := worker.GenerateID(numberOfIds)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ids) != numberOfIds {
		t.Errorf("Expected %d IDs, got %d", numberOfIds, len(ids))
	}

	// Check that all IDs are unique
	idMap := make(map[int64]bool)
	for _, id := range ids {
		if idMap[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		idMap[id] = true
	}
}

func TestGenerateID_ZeroOrNegativeNumber(t *testing.T) {
	provider := epoch.New(1420070400000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}

	// Test with 0
	ids, err := worker.GenerateID(0)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(ids) != 1 {
		t.Errorf("Expected 1 ID when requesting 0, got %d", len(ids))
	}

	// Test with negative number
	ids, err = worker.GenerateID(-5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(ids) != 1 {
		t.Errorf("Expected 1 ID when requesting negative, got %d", len(ids))
	}
}

func TestGenerateID_IDStructure(t *testing.T) {
	provider := epoch.New(1420070400000)
	workerId := int64(5)
	threadId := int64(3)
	worker := &WorkerVariant{
		WorkerID:     workerId,
		ThreadId:     threadId,
		TimeProvider: provider,
	}

	ids, err := worker.GenerateID(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	id := ids[0]

	// Extract worker ID from the generated ID
	extractedWorkerId := (id >> (ThreadBits + CounterBitSize)) & MaxNodeId
	if extractedWorkerId != workerId {
		t.Errorf("Expected worker ID %d, got %d", workerId, extractedWorkerId)
	}

	// Extract thread ID from the generated ID
	extractedThreadId := (id >> CounterBitSize) & ThreadCap
	if extractedThreadId != threadId {
		t.Errorf("Expected thread ID %d, got %d", threadId, extractedThreadId)
	}
}

func TestGenerateID_LargeNumberOfIDs(t *testing.T) {
	provider := epoch.New(1420070400000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}

	numberOfIds := 1000
	ids, err := worker.GenerateID(numberOfIds)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(ids) != numberOfIds {
		t.Errorf("Expected %d IDs, got %d", numberOfIds, len(ids))
	}

	// Verify all IDs are unique
	idMap := make(map[int64]bool)
	for _, id := range ids {
		if idMap[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		idMap[id] = true
	}
}

func TestConstants(t *testing.T) {
	// Test that bit sizes add up to a reasonable total (60 bits used, 4 bits unused at top)
	totalBits := UnusedBits + EpochBits + NodeIdBits + ThreadBits + CounterBitSize
	expectedTotalBits := int64(60) // Based on the current configuration
	if totalBits != expectedTotalBits {
		t.Errorf("Expected total bits to be %d, got %d", expectedTotalBits, totalBits)
	}

	// Test maximum values are correctly calculated
	expectedThreadCap := int64((1 << ThreadBits) - 1)
	if ThreadCap != expectedThreadCap {
		t.Errorf("Expected ThreadCap to be %d, got %d", expectedThreadCap, ThreadCap)
	}

	expectedMaxNodeId := int64((1 << NodeIdBits) - 1)
	if MaxNodeId != expectedMaxNodeId {
		t.Errorf("Expected MaxNodeId to be %d, got %d", expectedMaxNodeId, MaxNodeId)
	}

	expectedMaxCounter := int64((1 << CounterBitSize) - 1)
	if MaxCounter != expectedMaxCounter {
		t.Errorf("Expected MaxCounter to be %d, got %d", expectedMaxCounter, MaxCounter)
	}
	
	// Test that the bit layout makes sense
	if EpochBits <= 0 || NodeIdBits <= 0 || ThreadBits <= 0 || CounterBitSize <= 0 {
		t.Error("All bit sizes should be positive")
	}
	
	// Test that the maximum values are within reasonable ranges
	if ThreadCap <= 0 || ThreadCap > 256 {
		t.Errorf("ThreadCap should be reasonable, got %d", ThreadCap)
	}
	
	if MaxNodeId <= 0 || MaxNodeId > 256 {
		t.Errorf("MaxNodeId should be reasonable, got %d", MaxNodeId)
	}
	
	if MaxCounter <= 0 || MaxCounter > 2048 {
		t.Errorf("MaxCounter should be reasonable, got %d", MaxCounter)
	}
}
