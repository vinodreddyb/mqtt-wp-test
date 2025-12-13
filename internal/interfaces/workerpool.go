package interfaces

import (
	"context"
	"log"
	"sync"
	"time"

	"mqtt-wp-test/internal/domain"
)

type Processor interface {
	Process(context.Context, domain.Message) error
}

type WorkerPool struct {
	name    string
	jobs    chan domain.Message
	workers int
	proc    Processor
}

func NewWorkerPool(name string, workers, queue int, p Processor) *WorkerPool {
	return &WorkerPool{
		name:    name,
		jobs:    make(chan domain.Message, queue),
		workers: workers,
		proc:    p,
	}
}

func (w *WorkerPool) Jobs() chan<- domain.Message {
	return w.jobs
}

func (w *WorkerPool) Start(ctx context.Context, wg *sync.WaitGroup) {
	for i := 0; i < w.workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			log.Printf("[%s] worker-%d started", w.name, id)

			for {
				select {
				case <-ctx.Done():
					log.Printf("[%s] worker-%d stopping", w.name, id)
					return

				case msg := <-w.jobs:
					jobCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
					_ = w.proc.Process(jobCtx, msg)
					cancel()
				}
			}
		}(i)
	}
}
