package core

import (
	"sync"
)

type WorkerPool[TIn any] struct {
	jobs         chan TIn
	action       func(TIn)
	size         int
	workerPoolWg sync.WaitGroup
	done         chan struct{}
}

func NewWorkerPool[TIn any](size int, action func(TIn)) *WorkerPool[TIn] {
	pool := &WorkerPool[TIn]{
		jobs:   make(chan TIn, size),
		action: action,
		size:   size,
		done:   make(chan struct{}),
	}
	pool.start()
	return pool
}

func (wp *WorkerPool[TIn]) worker(stop chan struct{}) {
	wp.workerPoolWg.Add(1)
	defer wp.workerPoolWg.Done()

	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			wp.action(job)
		case <-stop:
			return
		}
	}
}

func (wp *WorkerPool[TIn]) start() {
	for i := 0; i < wp.size; i++ {
		stopChan := make(chan struct{})
		go wp.worker(stopChan)
	}
}

func (wp *WorkerPool[TIn]) stop() {
	close(wp.jobs)
	close(wp.done)
	wp.workerPoolWg.Wait()
}

func (wp *WorkerPool[TIn]) Run(in TIn) {
	wp.jobs <- in
}

func (wp *WorkerPool[TIn]) Close() error {
	go wp.stop()
	return nil
}
