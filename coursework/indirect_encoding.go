package main

import (
	"fmt"
	"time"
)

// [[u<v>.P]]= (νc)(u<c>|c(y).(y<v>|[[P]]))
func send_mona_sync(u chan chan chan int, v int) {
	c := make(chan chan int, 1)
	fmt.Printf("send_mona_sync: (nu c=%p)\n", c)
	go func() { u <- c }()
	fmt.Printf("send_mona_sync: u=%p <- c=%p\n", u, c)
	y := <-c
	fmt.Printf("send_mona_sync: y=%p <- c=%p\n", y, c)
	go func() { y <- v }()
	fmt.Printf("send_mona_sync: y=%p <- v=%d\n", y, v)
}

// [[u<v1,..,vn>.P]]= (νc)u<c>.c<v1>...c<vn>.[[P]]
func send_poly_sync(u chan chan chan chan int, v1 int, v2 int, v3 int) {
	c := make(chan chan chan int, 3)
	fmt.Printf("send_poly_sync: (nu c=%p)\n", c)
	u <- c
	fmt.Printf("send_poly_sync: u=%p <- c=%p\n", u, c)
	send_mona_sync(c, v1)
	send_mona_sync(c, v2)
	send_mona_sync(c, v3)
}

// [[u(x).P]]=u(y).(νd)(y<d>|d(x).[[P]])
func recv_mona_sync(u chan chan chan int) int {
	y := <-u
	fmt.Printf("recv_mona_sync: y=%d <- u=%p\n", y, u)
	d := make(chan int, 1)
	fmt.Printf("recv_mona_sync: (nu d=%p)\n", d)
	go func() { y <- d }()
	fmt.Printf("recv_mona_sync: y=%d <- d=%p\n", y, d)
	x := <-d
	fmt.Printf("recv_mona_sync: x=%d <- d=%p\n", x, d)
	return x
}

// [[u(x1,..,xn).P]]=u(z).z(x1)...z(xn).[[P]]
func recv_poly_sync(u chan chan chan chan int, x1 *int, x2 *int, x3 *int) {
	z := <-u
	fmt.Printf("recv_poly_sync: z=%p <- u=%p\n", z, u)
	*x1 = recv_mona_sync(z)
	*x2 = recv_mona_sync(z)
	*x3 = recv_mona_sync(z)
}

func main() {
	var i, j, k int
	u := make(chan chan chan chan int, 3)
	go func() {
		recv_poly_sync(u, &i, &j, &k)
	}()
	send_poly_sync(u, 1, 2, 3)

	time.Sleep(200 * time.Microsecond)
	fmt.Printf("Received: (%d, %d, %d)\n", i, j, k)
}
