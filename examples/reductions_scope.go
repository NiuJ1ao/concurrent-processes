// (nu b)(-a<b>) | a(x). -c<x>
package main

import (
	"fmt"
)

// (nu b)(-a<b>)
func new_b_out_a(a chan chan int) {
	// b is created here
	b := make(chan int, 1)
	fmt.Println("(nu b):", b)
	a <- b
	fmt.Println("-a<b>: output", b)
}

// a(x). -c<x>
func in_a_out_c(a, c chan chan int) {
	// b is received here, and sent along c
	b := <-a
	fmt.Println("a(x): input", b)
	c <- b
	fmt.Println("-c<x>: output", b)
}

func main() {
	a := make(chan chan int, 1)
	c := make(chan chan int, 1)

	go new_b_out_a(a)
	go in_a_out_c(a, c)

	// b is contained in channel c
	b := <-c
	fmt.Println("Reduced to: -c<b>: output", b)
}
