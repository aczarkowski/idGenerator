package generator

import (
	"errors"
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
	TimeProvider  timeprovider.TimeProvider // Used to get the current time either as epoch or Julian
}

// 64 bits UID
func (w *WorkerVariant) GenerateID(numberOfIds int) ([]int64, error) {
	var ids []int64
	currentTime := w.TimeProvider.GetTimeStamp()
	prevTime := w.lastTimeStamp
	if currentTime < prevTime {
		errors.New("invalid previous time stamp")
	}
	if numberOfIds <= 0 {
		numberOfIds = 1
	}
	var counter int64 = 0
	for {
		if counter == MaxCounter || prevTime == currentTime {
			nextTime := w.TimeProvider.GetTimeStamp()
			if nextTime == currentTime || prevTime == currentTime {
				continue
			}
			currentTime = nextTime
			counter = 0
		} else {
			id := currentTime << (NodeIdBits + ThreadBits + CounterBitSize)
			id |= w.WorkerID << (ThreadBits + CounterBitSize)
			id |= w.ThreadId << CounterBitSize
			id |= counter

			ids = append(ids, id)
			counter++
			if len(ids) == numberOfIds {
				w.lastTimeStamp = currentTime
				break
			}
		}
	}
	return ids, nil
}
