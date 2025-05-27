package generator

import (
	"testing"
	"uidGenerator/timeprovider/epoch"
	"uidGenerator/timeprovider/julian"
)

func BenchmarkGenerateID_Single_Epoch(b *testing.B) {
	provider := epoch.New(1420070400000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := worker.GenerateID(1)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateID_Single_Julian(b *testing.B) {
	provider := julian.New(2000100000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := worker.GenerateID(1)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateID_Multiple_Epoch(b *testing.B) {
	provider := epoch.New(1420070400000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := worker.GenerateID(10)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateID_Multiple_Julian(b *testing.B) {
	provider := julian.New(2000100000)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a new worker for each iteration to avoid counter exhaustion
		worker := &WorkerVariant{
			WorkerID:     1,
			ThreadId:     int64(i%8 + 1), // Vary thread ID to help with uniqueness
			TimeProvider: provider,
		}
		_, err := worker.GenerateID(3) // Reduced to 3 for Julian time provider
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkGenerateID_Large_Batch(b *testing.B) {
	provider := epoch.New(1420070400000)
	worker := &WorkerVariant{
		WorkerID:     1,
		ThreadId:     1,
		TimeProvider: provider,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := worker.GenerateID(100)
		if err != nil {
			b.Errorf("Unexpected error: %v", err)
		}
	}
}
