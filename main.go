package main

func main() {
	rc := make(chan Request)
	b := NewBalancer()

	go requester(rc)

	b.balance(rc)
}
