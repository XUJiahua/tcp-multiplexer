package serial

import (
	"github.com/sirupsen/logrus"
	"sync"
)

type Op int

const (
	None Op = iota
	Write
	Read
)

// return err means cancal
type Handler func() error

// NOTE: for tcp r/w call serialization
// 1 write should between 2 reads, 1 read should between 2 writes, except first write/read
// for example:
// w,r,w,r,...
// r,w,r,w,...

func SerialReadWriteWrapper(readHandler, writeHandler Handler, done chan struct{}) {
	logrus.Debug("start read/write serialization")
	// lock make sure only one Operation is running
	var lock sync.Mutex
	// lastOp make sure read is after write, or write is after read
	lastOp := None

	go func() {
		for {
			select {
			case <-done:
				logrus.Debug("close read goroutine")
				return
			default:
				lock.Lock()
				if lastOp == Read {
					lock.Unlock()
					continue
				}
				err := readHandler()
				if err != nil {
					lock.Unlock()
					continue
				}
				lastOp = Read
				lock.Unlock()
			}
		}
	}()

	for {
		select {
		case <-done:
			logrus.Debug("close write goroutine")
			return
		default:
			lock.Lock()
			if lastOp == Write {
				lock.Unlock()
				continue
			}
			err := writeHandler()
			if err != nil {
				lock.Unlock()
				continue
			}
			lastOp = Write
			lock.Unlock()
		}
	}
}
