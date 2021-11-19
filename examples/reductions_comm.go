// a(x). -x<c> | -a<b>
package main

import "fmt"

// a(x). -x<c>
func in_a_out_x(a chan chan int) {
	x := <-a
	fmt.Println("a(x): input", x)
	var i int
	fmt.Scanf("%d", &i)
	x <- i
	fmt.Println("x<c>: output", i)
}

// -a<b>
func out_a(a chan chan int, b chan int) {
	a <- b
	fmt.Println("a<b>: output", b)
}

func main() {
	a := make(chan chan int, 1)
	b := make(chan int, 1)

	/* a(x).-x<c> | -a<b> */
	go in_a_out_x(a)
	go out_a(a, b)

	/* value at b */
	fmt.Println("Channel b contains:", <-b)
}
