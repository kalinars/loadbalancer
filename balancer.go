package main

import (
	"container/heap"
	"log"
)

const nWorkers = 5
const workerBufferSize = 2

type Balancer struct {
	pool Pool
	done chan *Worker
}

func (b *Balancer) balance(work chan Request) {
	for {
		select {
		case req := <-work: // received a Request
			b.dispatch(req) // ... so send to a Worker
		case w := <-b.done: // a Worker says it's done
			b.completed(w) // ... so update data that the work is done
		}
	}
}

func (b *Balancer) dispatch(req Request) {
	// Grab the least loaded worker
	w := heap.Pop(&b.pool).(*Worker)
	if w.pending > workerBufferSize {
		log.Printf("[%d] Balancer: Queue full. Apply back pressure.", req.trId)
		req.c <- -1
		heap.Push(&b.pool, w)
		return
	}
	// ... send it the task
	w.requests <- req
	// Add one to its work queue
	w.pending++
	// Put it in its place on the heap
	heap.Push(&b.pool, w)
	log.Printf("[%d] Balancer: Request received and sent to Worker %d, pending %d tasks", req.trId, w.id, w.pending)
}

func (b *Balancer) completed(w *Worker) {
	// Remove one from its queue
	w.pending--
	// Remove it from heap
	heap.Remove(&b.pool, w.index)
	// Put it in its new place on the heap
	heap.Push(&b.pool, w)
}

func (b *Balancer) init(work chan Request) {
	// Init channel
	b.done = make(chan *Worker)

	// Init Pool
	log.Printf("balancer: starting %d workers", nWorkers)
	for i := 0; i < nWorkers; i++ {
		w := &Worker{id: i}
		w.requests = make(chan Request, workerBufferSize)
		heap.Push(&b.pool, w)
		go w.work(b.done)
	}
}
