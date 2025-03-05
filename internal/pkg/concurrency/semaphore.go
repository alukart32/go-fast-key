package concurrency

import "sync"

type Semaphore struct {
	count int
	max   int
	cond  sync.Cond
}

func NewSemaphore(limit int) *Semaphore {
	if limit <= 0 {
		limit = 1
	}
	return &Semaphore{
		max:  limit,
		cond: *sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Semaphore) Acquire() {
	if s == nil {
		return
	}

	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	for s.count >= s.max {
		s.cond.Wait()
	}

	s.count++
}

func (s *Semaphore) Release() {
	if s == nil {
		return
	}

	s.cond.L.Lock()
	defer s.cond.L.Unlock()

	s.count--

	s.cond.Signal()
}
