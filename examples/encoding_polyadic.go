package main

import (
	"fmt"
	"math/rand"
	"time"
)

func delay() {
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
}

type Triple struct {
	n1 int
	n2 int
	n3 int
}

func send_poly_wrong(c chan int, i int, j int, k int) {
	c <- i
	c <- j
	c <- k
}
func recv_poly_wrong(c chan int, i *int, j *int, k *int) {
	*i = <-c
	*j = <-c
	*k = <-c
}

func send_poly(cc chan chan int, i int, j int, k int) {
	c := make(chan int, 10)
	cc <- c
	c <- i
	c <- j
	c <- k
}

func recv_poly(cc chan chan int, i *int, j *int, k *int) {
	c := <-cc
	*i = <-c
	*j = <-c
	*k = <-c
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// In Go
	u := make(chan Triple, 1)
	u <- Triple{1, 2, 3}
	x := <-u
	fmt.Printf("Received: (%d, %d, %d)\n", x.n1, x.n2, x.n3)

	v := make(chan chan int, 10)
	// v := make(chan int, 10)
	for k := 0; k < 1000; k++ {
		go func() {
			// send_poly_wrong(v, 1, 2, 3)
			send_poly(v, 1, 2, 3)
		}()
		go func() {
			// send_poly_wrong(v, 4, 5, 6)
			send_poly(v, 1, 2, 3)
		}()
		var i, j, k int
		go func() {
			var i, j, k int
			delay()
			// recv_poly_wrong(v, &i, &j, &k)
			recv_poly(v, &i, &j, &k)
		}()
		// recv_poly_wrong(v, &i, &j, &k)
		recv_poly(v, &i, &j, &k)
		fmt.Printf("Received: (%d, %d, %d)\n", i, j, k)
		time.Sleep(200 * time.Millisecond)
	}
}
