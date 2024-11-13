package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync/atomic"
	"testing"
	"time"
)

type ActionIn struct {
	ID    int
	Value string
	t     *testing.T
}

var poolSize int = 10
var active = atomic.Int64{}
var allCount = atomic.Int64{}
var pool = NewWorkerPool(poolSize, WorkerAction)
var end = make(chan struct{})

func TestWorkerPool(t *testing.T) {

	go func() {
		for i := 0; i < 100; i++ {
			pool.Run(ActionIn{
				ID:    i,
				Value: fmt.Sprintf("%d", i),
				t:     t,
			})
		}
	}()

	<-end
	pool.stop()
}

func WorkerAction(in ActionIn) {
	active.Add(1)
	defer active.Add(-1)
	time.Sleep(time.Millisecond * 10)
	total := active.Load()
	assert.LessOrEqual(in.t, total, int64(poolSize))
	allCount.Add(1)
	count := allCount.Load()
	if count == int64(100) {
		close(end)
	}
}
