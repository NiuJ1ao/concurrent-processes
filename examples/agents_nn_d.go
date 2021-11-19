/* https://play.golang.org/p/l_RCwLYgFS7
 * NN(a) =^def !a(x).(nu b)-x<b>
 * D(a, b, c) =^def a(z).(-b<z> | -c<z>)
 * Example: NN<a> | D<b, c, d> | -a<b> | c(x).P | d(x).Q
 */
package main

import (
	"fmt"
	"math/rand"
	"time"
)

// NN(a) =^df !a(x).(nu b)-x<b>
func name_gen(a chan chan chan int) {
	for {
		x := <-a
		go func() {
			b := make(chan int, 1)
			x <- b
		}()
	}
}

// DD(a, b, c) =^df a(z).(-b<z> | -c<z>)
func duplicator(a, b, c chan chan int) {
	z := <-a
	go func() {
		b <- z
		c <- z
	}()
}

const STEPS = 5

func main() {
	a := make(chan chan chan int, 1)
	b := make(chan chan int, 1)
	c := make(chan chan int, 1)
	d := make(chan chan int, 1)
	go name_gen(a) // NN<a>

	for k := 0; k < STEPS; k++ {
		time.Sleep(1 * time.Second)
		go duplicator(b, c, d)
		a <- b
		i := rand.Intn(10)
		fmt.Printf("NN<a> | D<b, c, d> | -a<b> | c(x).-x<%d>.P | d(x).x(y).Q\t--->^*\t NN<a> | ", i)
		go func() {
			x := <-c
			x <- i // Think some complex computation
			fmt.Printf("P{%d/x} | ", x)
		}()
		x := <-d
		fmt.Printf("Q{%d/x}{%d/y}\n\n", x, <-x)
	}
}
