/* ATM(a, y) =^df a(z).ATM1<a,y,z>
 * ATM1(a, y, s) =^df s |> [ balance : -s<y>.ATM1<a,y,s>
 *                        || deposit : s(z).-s<y+z>.ATM1<a, y+z, s>
 *                         ]
 * Deposit(a, y) =^df (nu s)-a<s>. s <| deposit. -s<y>.s(x).0
 * Balance(a, c) =^df (nu s)-a<s>. s <| balance. s(x).-c<x>.0
 */
package main

import (
	"fmt"
	"sync"
	"time"
)

// s
type AtmChan struct {
	balance chan chan int
	deposit chan int
}

// ATM(a, y) =^df a(z).ATM1<a,y,z>
// ATM1(a, y, s) =^df s |> [ balance : -s<y>.ATM1<a,y,s>
//                        || deposit : s(z).-s<y+z>.ATM1<a, y+z, s>
//                         ]
func ATM(a <-chan AtmChan, initial_balance int) {
	var m sync.Mutex
	current_balance := initial_balance
	for {
		z := <-a
		go func() {
			select {
			case c := <-z.balance:
				m.Lock()
				defer m.Unlock()
				c <- current_balance
				fmt.Printf("read: %d\n", current_balance)
			case y := <-z.deposit:
				m.Lock()
				defer m.Unlock()
				current_balance += y
				fmt.Printf("deposit: %d\n", y)
			}
		}()
	}
}

// Deposit(a, y) =^df (nu s)-a<s>. s <| deposit. -s<y>.s(x).0
func Deposit(a chan<- AtmChan, i int) {
	s := AtmChan{make(chan chan int), make(chan int)}
	a <- s
	s.deposit <- i
}

// Balance(a, c) =^df (nu s)-a<s>. s <| balance. s(x).-c<x>.0
func Balance(a chan<- AtmChan) int {
	s := AtmChan{make(chan chan int), make(chan int)}
	a <- s
	r := make(chan int)
	s.balance <- r
	return <-r
}

const K = 100

func main() {
	a := make(chan AtmChan)
	go ATM(a, 0)

	for k := 0; k < K; k++ {
		Deposit(a, k)
		fmt.Printf("=======> %d\n", Balance(a))
	}

	time.Sleep(200 * time.Millisecond)
	fmt.Printf("=======> %d\n", Balance(a))

}
