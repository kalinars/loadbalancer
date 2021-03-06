package main

import (
	"log"
	"math/rand"
	"time"
)

type Worker struct {
	id       int          // for debugging purposes
	requests chan Request // work to do (buffered channel)
	pending  int          // count of pending tasks
	index    int          // index in the heap
}

func (w *Worker) work(done chan *Worker) {
	log.Printf("Worker started...")
	for {
		req := <-w.requests                                                  // get Request from balancer
		log.Printf("[%d] Worker %d: Request received", req.trId, w.id)
		time.Sleep(time.Duration(rand.Intn(5) * 10 * int(time.Millisecond))) // simulate work
		req.c <- req.fn()                                                    // call fn() and send result
		log.Printf("[%d] Worker %d: Request completed", req.trId, w.id)
		done <- w                                                            // we've finished this request
		log.Printf("[%d] Worker %d: Balancer received done", req.trId, w.id)
	}
}

type Pool []*Worker

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

func (p Pool) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].index = i
	p[j].index = j
}

func (p *Pool) Push(x interface{}) {
	item := x.(*Worker)
	item.index = len(*p)
	*p = append(*p, item)
}

func (p *Pool) Pop() interface{} {
	newLen := len(*p) - 1
	res := (*p)[newLen]
	res.index = -1 // some say for safety
	*p = (*p)[0:newLen]

	return res
}
