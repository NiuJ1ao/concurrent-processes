package main

import "fmt"

func test_golang() {

}

func main() {
	a := make(chan int, 5)

	a <- 1
	a <- 2
	a <- 3
	a <- 4
	a <- 5
	// a <- 6

	fmt.Println(<-a)
}
