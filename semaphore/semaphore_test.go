package semaphore

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestSemaphore(t *testing.T) {
	sem := New(5)
	var count int32
	for i := 0; i < 100; i++ {
		sem.Acquire(nil)
		go func(count *int32) {
			defer sem.Release()
			atomic.AddInt32(count, 1)
			time.Sleep(time.Second)
		}(&count)
	}
	sem.Wait()

	if count != 100 {
		t.Errorf("expect is %d, actual is %d", 100, count)
	}
}
