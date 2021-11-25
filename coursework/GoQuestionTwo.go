/*
 */

package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var chan_counter, msg_counter uint32

/*************************************************
 ****************Direct Encoding******************
 *************************************************/
// [[u<v1,v2,v3>.P]]=(νc)u<c>|c(y1).(y1<v1>|c(y2).(y2<v2>|c(y3).(y3<v3>|[[P]])))
func directSendPolySync(u chan chan chan int, vs []int) {
	c := make(chan chan int, 1) // (νc)
	atomic.AddUint32(&chan_counter, 1)

	// u<c>
	go func() {
		u <- c
		atomic.AddUint32(&msg_counter, 1)
	}()

	// c(y1).(y1<v1>|c(y2).(y2<v2>|c(y3).(y3<v3>|[[P]])))
	for _, v := range vs {
		y := <-c
		go func(v int) {
			y <- v
			atomic.AddUint32(&msg_counter, 1)
		}(v)
	}
}

// [[u(x1,x2,x3).P]]=u(y).(νd1,d2,d3)(y<d1>|d1(x1).(y<d2>|d2(x2).(y<d3>|d3(x3).[[P]])))
func directRecvPolySync(u chan chan chan int, xs []int) {
	y := <-u // u(y)
	n := len(xs)

	// (νd1,d2,d3)
	privateChans := make([]chan int, n)
	for i := 0; i < n; i++ {
		privateChans[i] = make(chan int, 1)
		atomic.AddUint32(&chan_counter, 1)
	}

	// (y<d1>|d1(x1).(y<d2>|d2(x2).(y<d3>|d3(x3).[[P]])))
	for i, channel := range privateChans {
		go func() {
			y <- channel
			atomic.AddUint32(&msg_counter, 1)
		}()
		xs[i] = <-channel
	}
}

func directEncoding(vs []int) {
	fmt.Println("Direct encoding")
	n := len(vs)
	xs := make([]int, n)

	u := make(chan chan chan int, n)
	atomic.AddUint32(&chan_counter, 1)

	go func() {
		directRecvPolySync(u, xs)
	}()
	directSendPolySync(u, vs)

	time.Sleep(200 * time.Millisecond)
	for _, x := range xs {
		fmt.Printf("Direct received: %d\n", x)
	}
}

/*************************************************
 ****************Indirect Encoding****************
 *************************************************/
// [[u<v>.P]]= (νc)(u<c>|c(y).(y<v>|[[P]]))
func sendMonaSync(u chan chan chan int, v int) {
	c := make(chan chan int, 1) // (νc)
	atomic.AddUint32(&chan_counter, 1)
	fmt.Printf("sendMonaSync: (nu c=%p)\n", c)

	// u<c>
	go func() {
		u <- c
		atomic.AddUint32(&msg_counter, 1)
	}()
	fmt.Printf("sendMonaSync: u=%p <- c=%p\n", u, c)

	y := <-c // c(y)
	fmt.Printf("sendMonaSync: y=%p <- c=%p\n", y, c)

	// (y<v>|[[P]])
	go func() {
		y <- v
		atomic.AddUint32(&msg_counter, 1)
	}()
	fmt.Printf("sendMonaSync: y=%p <- v=%d\n", y, v)
}

// [[u<v1,..,vn>.P]]= (νc)u<c>.c<v1>...c<vn>.[[P]]
func sendPolySync(u chan chan chan chan int, vs []int) {
	c := make(chan chan chan int, 1) // (νc)
	atomic.AddUint32(&chan_counter, 1)
	fmt.Printf("sendPolySync: (nu c=%p)\n", c)

	u <- c // u<c>
	atomic.AddUint32(&msg_counter, 1)
	fmt.Printf("sendPolySync: u=%p <- c=%p\n", u, c)

	// c<v1>...c<vn>.[[P]]
	for _, v := range vs {
		sendMonaSync(c, v)
	}
}

// [[u(x).P]]=u(y).(νd)(y<d>|d(x).[[P]])
func recvMonaSync(u chan chan chan int) int {
	y := <-u // u(y)
	fmt.Printf("recvMonaSync: y=%d <- u=%p\n", y, u)

	d := make(chan int, 1) // (νd)
	atomic.AddUint32(&chan_counter, 1)
	fmt.Printf("recvMonaSync: (nu d=%p)\n", d)

	// y<d>
	go func() {
		y <- d
		atomic.AddUint32(&msg_counter, 1)
	}()
	fmt.Printf("recvMonaSync: y=%d <- d=%p\n", y, d)

	x := <-d // d(x).[[P]]
	fmt.Printf("recvMonaSync: x=%d <- d=%p\n", x, d)

	return x
}

// [[u(x1,..,xn).P]]=u(z).z(x1)...z(xn).[[P]]
func recvPolySync(u chan chan chan chan int, xs []int) {
	z := <-u // u(z)
	fmt.Printf("recvPolySync: z=%p <- u=%p\n", z, u)

	// z(x1)...z(xn).[[P]]
	for i := 0; i < len(xs); i++ {
		xs[i] = recvMonaSync(z)
	}
}

func indirectEncoding(vs []int) {
	fmt.Println("Indirect encoding")
	n := len(vs)
	xs := make([]int, n)

	u := make(chan chan chan chan int, n)
	atomic.AddUint32(&chan_counter, 1)

	go func() {
		recvPolySync(u, xs)
	}()
	sendPolySync(u, vs)

	time.Sleep(200 * time.Microsecond)

	for _, x := range xs {
		fmt.Printf("Indirect received: %d\n", x)
	}
}

func main() {
	vs := []int{1, 2, 3}

	chan_counter, msg_counter = 0, 0
	indirectEncoding(vs)
	fmt.Printf("channel counter = %d, message counter = %d\n",
		chan_counter, msg_counter)

	fmt.Println()
	chan_counter, msg_counter = 0, 0
	directEncoding(vs)
	fmt.Printf("channel counter = %d, message counter = %d\n",
		chan_counter, msg_counter)
}
