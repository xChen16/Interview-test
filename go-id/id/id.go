package id

import (
	"errors"
	"sync"
	"time"
)

const (
	epoch            = int64(1618372800000)             // 常量时间戳2021/04/14 12:00:00 +0800 CST
	workerIDBits     = uint64(5)                        // 机器ID占5bit
	dataCenterIDBits = uint64(5)                        // 数据中心ID占5bit
	sequenceBits     = uint64(12)                       // 序列12bit
	maxWorkerID      = int64(-1 ^ (-1 << workerIDBits)) //节点ID的最大值 用于防止溢出
	maxDataCenterID  = int64(-1 ^ (-1 << dataCenterIDBits))
	maxSequence      = int64(-1 ^ (-1 << sequenceBits))

	workLeft = uint8(12) // workLeft = sequenceBits 节点IDx向左偏移量
	dataLeft = uint8(17) // dataLeft = dataCenterIDBits + sequenceBits
	timeLeft = uint8(22) // timeLeft = workerIDBits + sequenceBits 时间戳向左偏移量
)

type Worker struct {
	mu           sync.Mutex
	LastStamp    int64 // 记录时间戳
	WorkerID     int64
	DataCenterID int64
	Sequence     int64
}

//分布式情况下手动指定id
func NewWorker(workerID, dataCenterID int64) *Worker {
	return &Worker{
		WorkerID:     workerID,
		LastStamp:    0,
		Sequence:     0,
		DataCenterID: dataCenterID,
	}
}

// 获取毫秒
func (w *Worker) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

func (w *Worker) NextID() (uint64, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	return w.nextID()
}

func (w *Worker) nextID() (uint64, error) {
	timeStamp := w.getMilliSeconds()
	if timeStamp < w.LastStamp {
		return 0, errors.New("time is turning back")
	}

	if w.LastStamp == timeStamp {

		w.Sequence = (w.Sequence + 1) & maxSequence

		if w.Sequence == 0 {
			for timeStamp <= w.LastStamp {
				timeStamp = w.getMilliSeconds()
			}
		}
	} else {
		w.Sequence = 0
	}

	w.LastStamp = timeStamp
	id := ((timeStamp - epoch) << timeLeft) |
		(w.DataCenterID << dataLeft) |
		(w.WorkerID << workLeft) |
		w.Sequence

	return uint64(id), nil
}
