/* https://play.golang.org/p/9HKnOTSimwo
 * EQ(a, b) =^df !FW<a, b> | !FW<b, a>
 */
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// !FW(a, b) =^df !(a(z). -b<z>)
func repl_forwarder(a, b chan int) {
	z := <-a
	go repl_forwarder(a, b)
	b <- z
}

func equator(a, b chan int) {
	go repl_forwarder(a, b) // !FW<a, b>
	go repl_forwarder(b, a) // !FW<b, a>
}

const STEPS = 20

func main() {
	a := make(chan int, 100)
	b := make(chan int, 100)
	equator(a, b)

	for k := 0; k <= STEPS; k++ {
		time.Sleep(1 * time.Second)
		i := rand.Intn(10)
		j := rand.Intn(10)
		a <- i
		b <- j
		x := <-b
		y := <-a
		fmt.Printf("EQ<a, b> | -a<%d> | -b<%d> | b(x). a(y). P\t-->^*\tEQ<a,b> |", i, j)
		fmt.Printf(" P{%d/x}{%d/y}\n", x, y)

	}
}
