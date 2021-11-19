// -a<b> | !a(x). -a<x>
package main

import (
	"fmt"
)

// !a(x). -a<x>
func repl_in_a_out_a(a chan int) {
	x := <-a
	go repl_in_a_out_a(a) // Think of a server handling multiple requests in parallel
	fmt.Println("a(x):", x)
	a <- x + 1
	fmt.Println("-a<", x+1, ">")
}

// !a(x). -a<x>
func repl_in_a_out_a_loop(a chan int) {
	for {
		// TBD
	}
}

// -a<b>
func out_a_b(a chan int) {
	a <- 1
	fmt.Println("-a<", 1, ">")
}

func main() {
	a := make(chan int, 100)
	go out_a_b(a)
	go repl_in_a_out_a(a)

	var s string
	fmt.Scanf("%s", &s)
}
