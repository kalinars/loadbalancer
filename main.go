package main

func main() {
	rc := make(chan Request)
	b := Balancer{}
	b.init(rc)

	go requester(rc)

	b.balance(rc)
}
