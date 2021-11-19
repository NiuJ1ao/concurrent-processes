/* https://play.golang.org/p/r2h_uNxBnd5
 * NN(a) =^def !a(x).(nu b)-x<b>
 */
package main

import (
	"fmt"
	"time"
)

// NN(a) =^def !a(x).(nu b)-x<b>
func name_gen(a chan chan chan int) {
	for {
		x := <-a
		go func() {
			b := make(chan int)
			x <- b
		}()
	}
}

// NN1(a) =^def (nu b)!a(x).-x<b>
func name_gen1(a chan chan chan int) {
	b := make(chan int)
	for {
		x := <-a
		go func() { x <- b }()
	}
}

// NN2(a) =^def !(nu b)a(x).-x<b>
func name_gen2(a chan chan chan int) {
	for {
		b := make(chan int)
		x := <-a
		go func() { x <- b }()
	}
}

const STEPS = 5

func main() {
	a := make(chan chan chan int)
	// go name_gen(a)
	// go name_gen1(a)
	go name_gen2(a)

	for k := 0; k < STEPS; k++ {
		time.Sleep(1 * time.Second)
		// -a(c)|c(x).P
		c := make(chan chan int)
		a <- c
		b := <-c
		fmt.Printf("NN<a> | -a<c> | c(x). P \t-->^*\t NN<a> | ")
		fmt.Printf("P{%d/x}\n", b)
	}
}
