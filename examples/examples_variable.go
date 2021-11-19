/*
 * Var(a, x) =^df a(z). z |> [read : -z<x>. Var<a, x> || write : z(y). Var<a, y>]
 * Reader(a, c) =^df (nu s)(-a<s>. s <| read. s(y).-c<y>
 * Writer(a, x) =^df (nu s)(-a<s>. s <| write. -s<x>.0
 */
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ReadOrWrite struct {
	read  chan chan int
	write chan int
}

/*
 * Var(a, x) =^df a(z). z |> [read : -z<x>. Var<a, x> || write : z(y). Var<a, y>]
 */
func Var(a <-chan ReadOrWrite, i int) {
	var m sync.Mutex
	x := i
	for {
		z := <-a
		go func() {
			select {
			case c := <-z.read:
				m.Lock()
				defer m.Unlock()
				c <- x
				fmt.Printf("Variable read: %d\n", x)
			case y := <-z.write:
				m.Lock()
				defer m.Unlock()
				x = y
				fmt.Printf("Variable updated: %d\n", x)
			}
		}()
	}
}

// Reader(a, c) =^df (nu s)(-a<s>. s <| read. s(y).-c<y>
func Reader(a chan<- ReadOrWrite) int {
	s := ReadOrWrite{make(chan chan int), make(chan int)}
	a <- s
	r := make(chan int)
	s.read <- r
	return <-r
}

// Writer(a, x) =^df (nu s)(-a<s>. s <| write. -s<x>.0
func Writer(a chan<- ReadOrWrite, i int) {
	s := ReadOrWrite{make(chan chan int), make(chan int)}
	a <- s
	s.write <- i
}

const K = 100

func main() {
	a := make(chan ReadOrWrite)
	go Var(a, 10)

	for k := 0; k <= K; k++ {
		Writer(a, rand.Intn(100))
		fmt.Printf("=======> %d\n", Reader(a))
		time.Sleep(200 * time.Millisecond)
	}

}
