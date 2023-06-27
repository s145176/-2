package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Pool struct {
	workers []*worker
	m       sync.Mutex
	cap     int
	length  int32
}

func NewPool(cap int) *Pool {
	return &Pool{
		workers: make([]*worker, cap, cap),
		cap:     cap,
		length:  int32(cap),
	}
}

func (p *Pool) Run() {

	for i := 0; i < p.cap; i++ {
		p.workers[i] = newWorker(p)
		go p.workers[i].run()
		fmt.Println("worker", i, "running")
	}
}

func (p *Pool) SubMit(task Task) {
	for {
		if atomic.AddInt32(&p.length, -1) >= 0 {
			p.m.Lock()
			worker := p.workers[0]
			p.workers = p.workers[1:]
			p.m.Unlock()
			worker.taskChan <- task
			break
		}
		atomic.AddInt32(&p.length, 1)
	}
}

type Task func()

type worker struct {
	taskChan chan Task
	pool     *Pool
}

func newWorker(pool *Pool) *worker {
	return &worker{
		taskChan: make(chan Task, 1),
		pool:     pool,
	}
}

func (w *worker) run() {

	for v := range w.taskChan {
		v()
		atomic.AddInt32(&w.pool.length, 1)
		w.pool.m.Lock()
		w.pool.workers = append(w.pool.workers, w)
		w.pool.m.Unlock()
	}

}
