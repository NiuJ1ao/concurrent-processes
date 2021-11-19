/* https://play.golang.org/p/nQ4ZG-LuhKQ
 * D(a, b, c) =^def a(z).(-b<z> | -c<z>)
 */
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// a(z).(-b<z> | -c<z>)
func duplicator(a, b, c chan int) {
	z := <-a
	go func() {
		b <- z
	}()
	c <- z
}

const STEPS = 20

func main() {
	a := make(chan int, 1)
	b := make(chan int, 1)
	c := make(chan int, 1)

	for k := 0; k < STEPS; k++ {
		time.Sleep(500 * time.Millisecond)
		go duplicator(a, b, c)
		i := rand.Intn(10)
		a <- i
		fmt.Printf("D<a, b, c> | -a<%d> | b(x). c(y). P\t--->^*\t", i)
		fmt.Printf("P{%d/x}{%d/y}\n", <-b, <-c)
	}
}
