package main

import (
	"fmt"
	"math/rand"
	"time"
)

func delay() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

// [[u<v>.P]]= (νc)(u<c>|c(y).(y<v>|[[P]]))
func sync_send(u chan chan chan int, v int) {
	c := make(chan chan int, 1) // (nu c)
	fmt.Printf("send: (nu c = %d)\n", c)
	go func() { u <- c }() // -u<c>
	fmt.Printf("send: -u<%d>\n", c)
	y := <-c // c(y)
	fmt.Printf("send: c(%d)\n", y)
	go func() {
		y <- v
	}() // -y<v>
	fmt.Printf("send: -%d < %d >\n", y, v)
}

func sync_recv(u chan chan chan int) int {
	y := <-u // u(y)
	fmt.Printf("recv: u(%d)\n", y)
	d := make(chan int, 1) // (nu d)
	fmt.Printf("recv: (nu d=%d)\n", d)
	go func() { y <- d }() // -y<d>
	fmt.Printf("recv: -%d < %d > \n", y, d)
	x := <-d // d(x)
	fmt.Printf("recv: %d ( %d ) \n", d, x)
	return x
}

func main() {
	rand.Seed(time.Now().UnixNano())
	u := make(chan chan chan int, 1)
	// // Deadlocks, as expected, similar to using regular go channels
	// sync_send(u, 10)
	// sync_recv(u)

	var x int
	go func() {
		x = sync_recv(u)
	}()
	sync_send(u, rand.Intn(100))

	// sync_send cannot continue until sync_recv has received a value,
	// but it can continue before the function returns
	time.Sleep(1 * time.Microsecond)
	fmt.Printf("P{%d/x}\n", x)

	// // Regular go
	// c := make(chan int)
	// var y int
	// go func() {
	// 	y = <-c
	// }()
	// c <- rand.Intn(100)

	// fmt.Printf("Q{%d/x}\n", y)
}
