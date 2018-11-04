package handlers

import (
	"log"
)

type ConnLimiter struct {
	max    int
	bucket chan int
}

func NewConnLimiter(cc int) *ConnLimiter {
	return &ConnLimiter{
		max:    cc,
		bucket: make(chan int, cc),
	}
}

func (cl *ConnLimiter) GetConn() bool {

	if len(cl.bucket) >= cl.max {
		log.Print("reached the rate limitation ")
		return false
	}
	cl.bucket <- 1
	return true
}

func (cl *ConnLimiter) ReleaseConn() {

	c := <-cl.bucket
	log.Printf("new connection can come %v", c)
}
