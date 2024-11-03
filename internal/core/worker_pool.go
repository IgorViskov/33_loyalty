package core

import (
	"github.com/labstack/gommon/log"
	"sync"
)

type WorkerPool[TIn any, TOut any] struct {
	jobs         chan TIn
	results      chan Result[TOut]
	action       func(TIn) Result[TOut]
	handler      func(Result[TOut]) error
	size         int
	workerPoolWg sync.WaitGroup
	done         chan struct{}
}

func NewWorkerPool[TIn any, TOut any](size int, action func(TIn) Result[TOut], handler func(Result[TOut]) error) *WorkerPool[TIn, TOut] {
	pool := &WorkerPool[TIn, TOut]{
		jobs:    make(chan TIn, size),
		results: make(chan Result[TOut], size),
		action:  action,
		handler: handler,
		size:    size,
		done:    make(chan struct{}),
	}
	pool.start()
	return pool
}

func (wp *WorkerPool[TIn, TOut]) worker(stop chan struct{}) {
	wp.workerPoolWg.Add(1)
	defer wp.workerPoolWg.Done()

	for {
		select {
		case job, ok := <-wp.jobs:
			if !ok {
				return
			}
			wp.results <- wp.action(job)
		case <-stop:
			return
		}
	}
}

func (wp *WorkerPool[TIn, TOut]) handle() {
	for {
		select {
		case result, ok := <-wp.results:
			if !ok {
				return
			}
			err := wp.handler(result)
			if err != nil {
				log.Error(err)
			}
		case <-wp.done:
			return
		}
	}
}

func (wp *WorkerPool[TIn, TOut]) start() {
	go wp.handle()

	for i := 0; i < wp.size; i++ {
		stopChan := make(chan struct{})
		go wp.worker(stopChan)
	}
}

func (wp *WorkerPool[TIn, TOut]) stop() {
	close(wp.jobs)
	close(wp.done)
	wp.workerPoolWg.Wait()
	close(wp.results)
}

func (wp *WorkerPool[TIn, TOut]) Run(in TIn) {
	wp.jobs <- in
}

func (wp *WorkerPool[TIn, TOut]) Close() error {
	go wp.stop()
	return nil
}
