package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var initialTime = uint32(time.Now().Unix())
var trCounter uint32

type Request struct {
	trId uint32     // Transaction ID
	fn   func() int // Operation to perform
	c    chan int   // The channel to return the result
}

func requester(work chan<- Request) {
	http.HandleFunc("/requester", handleRequest(work))
	serverPort := "8888"
	log.Printf("Requester: Server listening on port %s", serverPort)
	http.ListenAndServe(":" + serverPort, nil)

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
		trId := getTransactionId()
		log.Printf("[%d] Requester: new client request", trId)

		c := make(chan int)
		defer close(c)
		work <- Request{trId, func() int { return rand.Intn(nWorkers * 2) }, c} // send Request to Balancer
		result := <-c                                                           // wait for answer

		if result >= 0 {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Result %d is processed.", result)
			log.Printf("[%d] Requester: Result %d is further processed.", trId, result)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			fmt.Fprintf(w, "Error 503 - Service is overloaded.")
		}
	}
}

func getTransactionId() uint32 {
	epoch := uint32(time.Date(2016, 1, 1, 0, 0, 0, 0, time.UTC).Unix())
	id := initialTime - epoch
	id = id << (32 - 10)          // 10 bits countain arpoximately 50 years in seconds
	id |= trCounter & (1<<11 - 1) // get the last 10 bits from the counter

	trCounter++
	return id
}
