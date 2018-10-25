package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	//"time"
)

type Request struct {
	fn func() int // Operation to perform
	c  chan int   // The channel to return the result
}

func requester(work chan<- Request) {
	http.HandleFunc("/requester", handleRequest(work))
	http.ListenAndServe(":8888", nil)

	/*
		// Fake load
		for {
			time.Sleep(time.Duration(rand.Int63n(nWorkers * int64(time.Second))))
			log.Printf("requester sends Request")
			work <- Request{func() int { return rand.Intn(nWorkers*2) }, c} // send request
			result := <-c		   // wait for answer

			log.Printf("Result %d is further processed.", result)
		}
	*/
}

func handleRequest(work chan<- Request) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("client request received")

		c := make(chan int)
		work <- Request{func() int { return rand.Intn(nWorkers * 2) }, c} // send Request to Balancer
		result := <-c                                                     // wait for answer

		if result >= 0 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Result %d is processed.", result)
			log.Printf("Result %d is further processed.", result)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Error 503 - Service is overloaded.")
		}
	}
}
