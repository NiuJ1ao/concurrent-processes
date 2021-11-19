// -a<b> | -a<d> | a(x).-c<x>
package main

import (
	"fmt"
	"time"
)

func main() {
	a := make(chan int, 2)
	c := make(chan int, 2)

	// -a<b>
	go func() {
		a <- 1 // b = 1
	}()

	// -a<d>
	go func() {
		time.Sleep(10000 * time.Nanosecond)
		a <- 2 // d = 2
	}()

	x := <-a
	c <- x
	fmt.Println(<-c)

	x = <-a
	c <- x
	fmt.Println(<-c)
}
