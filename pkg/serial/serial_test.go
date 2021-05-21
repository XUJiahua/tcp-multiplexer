package serial

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"math/rand"
	"testing"
	"time"
)

func TestSerialReadWriteWrapper(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	done := make(chan struct{})

	var arr []int
	go func() {
		SerialReadWriteWrapper(func() error {
			time.Sleep(time.Millisecond)
			arr = append(arr, 1)
			return nil
		}, func() error {
			time.Sleep(time.Millisecond)
			arr = append(arr, 0)
			return nil
		}, done)
	}()

	time.Sleep(time.Second)

	close(done)
	time.Sleep(time.Second)
	for i := 1; i < len(arr); i++ {
		if arr[i] == arr[i-1] {
			spew.Dump(arr)
			panic(arr)
		}
	}
}

func TestSerialReadWriteWrapper2(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	done := make(chan struct{})

	var arr []int
	go func() {
		SerialReadWriteWrapper(func() error {
			time.Sleep(time.Millisecond)
			if rand.Int()%2 == 0 {
				return errors.New("123")
			}
			arr = append(arr, 1)
			return nil
		}, func() error {
			time.Sleep(time.Millisecond)
			if rand.Int()%2 == 0 {
				return errors.New("123")
			}
			arr = append(arr, 0)
			return nil
		}, done)
	}()

	time.Sleep(time.Second)

	close(done)
	time.Sleep(time.Second)
	for i := 1; i < len(arr); i++ {
		if arr[i] == arr[i-1] {
			spew.Dump(arr)
			panic(arr)
		}
	}
}
