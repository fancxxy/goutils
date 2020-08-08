package semaphore

import (
	"context"
	"sync"
)

// Semaphore 信号量
type Semaphore struct {
	ch chan struct{}
	wg *sync.WaitGroup
}

// New 创建size的信号量
func New(size int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, size),
		wg: new(sync.WaitGroup),
	}
}

// Acquire 获取
func (s *Semaphore) Acquire(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	case s.ch <- struct{}{}:
		break
	}
	s.wg.Add(1)
	return nil
}

// Release 释放
func (s *Semaphore) Release() {
	<-s.ch
	s.wg.Done()
}

// Wait 同步
func (s *Semaphore) Wait() {
	s.wg.Wait()
}
