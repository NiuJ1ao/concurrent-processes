package main

import (
	"fmt"
	"time"
)

// [[u<v1,v2,v3>.P]]=(νc)u<c>|c(y1).(y1<v1>|c(y2).(y2<v2>|c(y3).(y3<v3>|[[P]])))
func send_poly_sync(u chan chan chan int, v1 int, v2 int, v3 int) {
	c := make(chan chan int, 3)
	go func() { u <- c }()
	y1 := <-c
	go func() { y1 <- v1 }()
	y2 := <-c
	go func() { y2 <- v2 }()
	y3 := <-c
	go func() { y3 <- v3 }()
}

// [[u(x1,x2,x3).P]]=u(y).(νd1,d2,d3)(y<d1>|d1(x1).(y<d2>|d2(x2).(y<d3>|d3(x3).[[P]])))
func recv_poly_sync(u chan chan chan int, x1 *int, x2 *int, x3 *int) {
	y := <-u
	d1 := make(chan int, 1)
	d2 := make(chan int, 1)
	d3 := make(chan int, 1)
	go func() { y <- d1 }()
	*x1 = <-d1
	go func() { y <- d2 }()
	*x2 = <-d2
	go func() { y <- d3 }()
	*x3 = <-d3
}

func main() {
	var i, j, k int
	u := make(chan chan chan int, 3)
	go func() {
		recv_poly_sync(u, &i, &j, &k)
	}()
	send_poly_sync(u, 1, 2, 3)

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Received: (%d, %d, %d)\n", i, j, k)
}
