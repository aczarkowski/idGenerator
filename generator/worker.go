package generator

import (
	"errors"
	"sync"
	"time"
	"uidGenerator/timeprovider"
)

var UnusedBits int64 = 1
var EpochBits int64 = 41      // 42 bits
var NodeIdBits int64 = 3      // 49 bits
var ThreadBits int64 = 5      // 57 bits
var CounterBitSize int64 = 10 // 64 bits UID
var ThreadCap int64 = (1 << ThreadBits) - 1
var MaxNodeId int64 = (1 << NodeIdBits) - 1
var MaxCounter int64 = (1 << CounterBitSize) - 1

type WorkerVariant struct {
	WorkerID      int64                     // It is the Node ID
	ThreadId      int64                     // Will be assigned during startup
	lastTimeStamp int64                     //Used to remember the last time stamp
	lastCounter   int64                     //Used to remember the last counter value
	TimeProvider  timeprovider.TimeProvider // Used to get the current time either as epoch or Julian
	mutex         sync.Mutex               // Ensures thread-safe access to worker state
}

// 64 bits UID
func (w *WorkerVariant) GenerateID(numberOfIds int) ([]int64, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	
	var ids []int64
	currentTime := w.TimeProvider.GetTimeStamp()
	if currentTime < w.lastTimeStamp {
		return nil, errors.New("invalid previous time stamp")
	}
	if numberOfIds <= 0 {
		numberOfIds = 1
	}
	
	var counter int64
	
	// If we're in the same timestamp as last generation, continue from last counter
	if currentTime == w.lastTimeStamp {
		counter = w.lastCounter + 1
	} else {
		// New timestamp, reset counter
		counter = 0
	}
	
	for {
		// Check if we've exhausted the counter for this timestamp
		if counter > MaxCounter {
			// Wait for next timestamp
			for {
				nextTime := w.TimeProvider.GetTimeStamp()
				if nextTime > currentTime {
					currentTime = nextTime
					counter = 0
					break
				}
				// If timestamp hasn't changed, wait briefly before checking again
				// This handles high-frequency scenarios without returning an error
				time.Sleep(time.Nanosecond)
			}
		}
		
		id := currentTime << (NodeIdBits + ThreadBits + CounterBitSize)
		id |= w.WorkerID << (ThreadBits + CounterBitSize)
		id |= w.ThreadId << CounterBitSize
		id |= counter

		ids = append(ids, id)
		counter++
		
		if len(ids) == numberOfIds {
			w.lastTimeStamp = currentTime
			w.lastCounter = counter - 1 // Store the last used counter
			break
		}
	}
	return ids, nil
}
