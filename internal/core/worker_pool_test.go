package core

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type ActionIn struct {
	ID    int
	Value string
	ctx   context.Context
}

type ActionOut struct {
	ID  int
	err error
}

var poolSize int = 10
var active = atomic.Int64{}
var allCount = atomic.Int64{}
var pool = NewWorkerPool((int)(poolSize), WorkerAction, WorkerHandler)
var activeIds = sync.Map{}
var end = make(chan struct{})

func TestWorkerPool(t *testing.T) {

	go func() {
		for i := 0; i < 100; i++ {
			pool.Run(ActionIn{
				ID:    i,
				Value: fmt.Sprintf("%d", i),
			})
		}
	}()

	<-end
	pool.stop()
}

func WorkerAction(in ActionIn) Result[ActionOut] {
	active.Add(1)
	defer active.Add(-1)
	activeIds.Store(in.ID, struct{}{})
	time.Sleep(time.Millisecond * 1000)
	var err error
	total := active.Load()
	if total > int64(poolSize) {
		err = errors.New("pool size exceeded")
	}
	fmt.Printf("Work ID [%d], Active count: [%d] \r\n", in.ID, total)
	activeIds.Delete(in.ID)
	return Done(&ActionOut{
		ID:  in.ID,
		err: err,
	})
}

func WorkerHandler(r Result[ActionOut]) error {
	allCount.Add(1)
	count := allCount.Load()
	fmt.Printf("Responce Id [%d] received by [%d] \r\n", r.data.ID, allCount.Load())
	if count == int64(100) {
		close(end)
	}
	return nil
}
