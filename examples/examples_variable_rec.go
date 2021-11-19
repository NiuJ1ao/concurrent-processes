/*
 * Var(a, x) =^df a(z). Var1<a,z,x>
 * Var1(a, s, x) =^df s |> [ read : -s<x>. Var1<a, s, x>
 *                        || write : s(y). Var1<a, s, y>
 *                        || stop: Var<a, x>
 *                         ]
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
	stop  chan struct{}
}

// Var(a, x) =^df !a(z). Var1<a,z,x>
func Var(a <-chan ReadOrWrite, i int) {
	var m sync.Mutex
	x := i
	// Var1(a, s, x) =^df s |> [ read : -s<x>. Var1<a, s, x>
	//                        || write : s(y). Var1<a, s, y>
	//                        || stop: 0
	//                         ]
	Var1 := func(s ReadOrWrite) {
		for {
			select {
			case c := <-s.read:
				func() {
					m.Lock()
					defer m.Unlock()
					c <- x
					fmt.Printf("Variable read: %d\n", x)
				}()
			case y := <-s.write:
				func() {
					m.Lock()
					defer m.Unlock()
					x = y
					fmt.Printf("Variable updated: %d\n", x)
				}()
			case <-s.stop:
				return
			}
		}
	}
	for {
		s := <-a
		go Var1(s)
	}
}

// Reader(a, c) =^df (nu s)(-a<s>. s <| read. s(y).-c<y>. s <| read. s(y). -c<y> ... s <| stop, 0
func Reader(a chan<- ReadOrWrite, i int) []int {
	rd := make([]int, i)
	s := ReadOrWrite{make(chan chan int), make(chan int), make(chan struct{})}
	a <- s
	r := make(chan int)
	for k := 0; k < i; k++ {
		time.Sleep(200 * time.Millisecond)
		s.read <- r
		rd[k] = <-r
	}
	s.stop <- struct{}{}
	return rd
}

// Writer(a, x) =^df (nu s)(-a<s>. s <| write. -s<x>.0
func Writer(a chan<- ReadOrWrite, i int) {
	s := ReadOrWrite{make(chan chan int), make(chan int), make(chan struct{})}
	a <- s
	s.write <- i
	s.stop <- struct{}{}
}

const K = 10

func main() {
	a := make(chan ReadOrWrite)
	go Var(a, 10)

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Printf("=======> %d\n", Reader(a, K))
	}()
	for k := 0; k <= K; k++ {
		Writer(a, rand.Intn(100))
		time.Sleep(200 * time.Millisecond)
		if k == K/2 {
			go func() {
				defer wg.Done()
				fmt.Printf("=======> %d\n", Reader(a, K))
			}()
		}
	}
	wg.Wait()

}
