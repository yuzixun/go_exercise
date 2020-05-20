package main

import (
	"fmt"
	"sync"
	"time"
)

const beginTime = 1420041600000000000

// const reservreBits = 1
// const timestampBits = 41
const dataCenterBits = 5
const workerBits = 5
const sequenceBits = 12

const maxWorkerID = -1 ^ (-1 << workerBits)
const maxDataCenterID = -1 ^ (-1 << dataCenterBits)

const workerIDShift = sequenceBits
const dataCenterIDShift = sequenceBits + workerBits
const timestampShift = sequenceBits + workerBits + dataCenterBits

const sequenceMask = -1 ^ (-1 << sequenceBits)

type generator struct {
	workerID      int64
	dataCenterID  int64
	sequence      int64
	lastTimestamp int64
	mu            sync.Mutex
}

func Constructor(workerID, dataCenterID int64) *generator {
	if workerID > maxWorkerID {
		panic(fmt.Sprintf("worker id cannot be greater than %d", maxWorkerID))
	}

	if workerID < 0 {
		panic("worker id cannot be less than 0")
	}

	if dataCenterID > maxDataCenterID {
		panic(fmt.Sprintf("data center id cannot be greater than %d", maxDataCenterID))
	}

	if dataCenterID < 0 {
		panic("data center id cannot be less than 0")
	}

	return &generator{
		workerID:      workerID,
		dataCenterID:  dataCenterID,
		sequence:      0,
		lastTimestamp: -1,
	}
}

func (this *generator) nextId() int64 {
	this.mu.Lock()
	defer this.mu.Unlock()

	timestamp := this.curTime()
	if timestamp < this.lastTimestamp {
		panic("Clock moved backwards.")
	}

	if timestamp == this.lastTimestamp {
		this.sequence = (this.sequence + 1) & sequenceMask

		if this.sequence == 0 {
			timestamp = this.blockToNextMillis(this.lastTimestamp)
		}
	} else {
		this.sequence = 0
	}

	this.lastTimestamp = timestamp

	return (timestamp-beginTime)<<timestampShift |
		(this.dataCenterID << dataCenterIDShift) |
		(this.workerID << workerIDShift) | this.sequence
}

func (this *generator) blockToNextMillis(lastTimestamp int64) int64 {
	timestamp := this.curTime()

	for timestamp <= this.lastTimestamp {
		fmt.Println(timestamp, this.lastTimestamp)
		timestamp = this.curTime()
	}

	return timestamp
}

func (this *generator) curTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
