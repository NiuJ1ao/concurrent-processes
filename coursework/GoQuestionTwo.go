/*
 * In this question, we input a list of three integers as example (it is also compatible with different size).
 * Since there always is a receiver if there is a sender, we only count the number of senders for messages.
 * In this case, the indirect encoding requires 10 channels and 13 messages, where the direct encoding only
 * requires 5 channels and 7 messages.
 * Channels:
 * 		For indirect sender, it requires a private channel of every variable and a extra channel
 * 		for the first sender created in the polyadic synchronization, in the monatic synchronization to
 * 		monatic asynchronization encoding (n + 1), whereas the direct sender only needs one private channel.
 * 		For indrect receiver, it also require a extra channle for the first receiver in the polyadic synchronization (1).
 *		In this case, delta = 3 + 1 + 1 = 5 is matched.
 * Messages:
 *		For indirect sender, it requires sends for channels of every variable. For the first send of private channel
 * 		in the polyadic synchronization, two more sends are required than direct encoding (n + 2).
 * 		For direct receiver, send a extra message for first receiver (1).
 *		In this case, delta = 3 + 2 + 1 = 6 is matched.
 */

package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

var chan_counter, msg_counter uint32

/***************************************************
 **************** Direct Encoding ******************
 ***************************************************/
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

/***************************************************
 **************** Indirect Encoding ****************
 ***************************************************/
// [[u<v>.P]]= (νc)(u<c>|c(y).(y<v>|[[P]]))
func sendMonaValSync(u chan chan chan int, v int) {
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

// [[u<v>.P]]= (νc)(u<c>|c(y).(y<v>|[[P]]))
func sendMonaChanSync(u chan chan chan chan chan chan int, v chan chan chan int) {
	c := make(chan chan chan chan chan int, 1) // (νc)
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
func sendPolySync(u chan chan chan chan chan chan int, vs []int) {
	c := make(chan chan chan int, 1) // (νc)
	atomic.AddUint32(&chan_counter, 1)
	fmt.Printf("sendPolySync: (nu c=%p)\n", c)

	sendMonaChanSync(u, c) // u<c>
	atomic.AddUint32(&msg_counter, 1)
	fmt.Printf("sendPolySync: u=%p <- c=%p\n", u, c)

	// c<v1>...c<vn>.[[P]]
	for _, v := range vs {
		sendMonaValSync(c, v)
	}
}

// [[u(x).P]]=u(y).(νd)(y<d>|d(x).[[P]])
func recvMonaValSync(u chan chan chan int) int {
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

// [[u(x).P]]=u(y).(νd)(y<d>|d(x).[[P]])
func recvMonaChanSync(u chan chan chan chan chan chan int) chan chan chan int {
	y := <-u // u(y)
	fmt.Printf("recvMonaSync: y=%d <- u=%p\n", y, u)

	d := make(chan chan chan chan int, 1) // (νd)
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
func recvPolySync(u chan chan chan chan chan chan int, xs []int) {
	z := recvMonaChanSync(u) // u(z)
	fmt.Printf("recvPolySync: z=%p <- u=%p\n", z, u)

	// z(x1)...z(xn).[[P]]
	for i := 0; i < len(xs); i++ {
		xs[i] = recvMonaValSync(z)
	}
}

func indirectEncoding(vs []int) {
	fmt.Println("Indirect encoding")
	n := len(vs)
	xs := make([]int, n)

	u := make(chan chan chan chan chan chan int, n)
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
