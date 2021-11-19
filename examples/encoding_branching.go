package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Tag int

// labels
const (
	A Tag = iota
	B
	C
)

type Val struct {
	tag Tag
	val int
}

type Triple struct {
	p1 chan struct{}
	p2 chan struct{}
	p3 chan struct{}
}

func branch_pi(u chan chan Triple) Tag {
	t := <-u
	cs := Triple{make(chan struct{}), make(chan struct{}), make(chan struct{})}
	t <- cs
	var wg sync.WaitGroup
	wg.Add(1)
	var i Tag
	go func() {
		defer wg.Done()
		<-cs.p1
		i = A
	}()
	go func() {
		defer wg.Done()
		<-cs.p2
		i = B
	}()
	go func() {
		defer wg.Done()
		<-cs.p3
		i = C
	}()
	wg.Wait()
	return i
}

func sel_pi(u chan chan Triple, t Tag) {
	c := make(chan Triple)
	u <- c
	tr := <-c
	switch t {
	case A:
		tr.p1 <- struct{}{}
	case B:
		tr.p2 <- struct{}{}
	case C:
		tr.p3 <- struct{}{}
	default:
		fmt.Printf("Should never end here")
	}

}

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Printf("\nGo Option 1\n")
	// go option 1
	c1 := make(chan Val, 1)
	c1 <- Val{Tag(rand.Intn(3)), rand.Intn(100)}
	switch x1 := <-c1; x1.tag {
	case A:
		fmt.Printf("A: %d\n", x1.val)
	case B:
		fmt.Printf("B: %d\n", x1.val)
	case C:
		fmt.Printf("C: %d\n", x1.val)
	default:
		fmt.Printf("Should never end here")
	}

	fmt.Printf("\nGo Option 2\n")
	// go option 2
	var alts = []chan int{
		make(chan int),
		make(chan int),
		make(chan int),
	}
	go func() { alts[rand.Intn(3)] <- rand.Intn(100) }()
	select {
	case x := <-alts[0]:
		fmt.Printf("A: %d\n", x)
	case x := <-alts[1]:
		fmt.Printf("B: %d\n", x)
	case x := <-alts[2]:
		fmt.Printf("C: %d\n", x)
	}

	// go option 3: using reflect.Select, etc
	// ...

	fmt.Printf("\nPi Calculus encoding\n")
	// pi calculus encoding
	u := make(chan chan Triple, 1)
	go func() { sel_pi(u, Tag(rand.Intn(3))) }()
	x := branch_pi(u)
	switch x {
	case A:
		fmt.Printf("A\n")
	case B:
		fmt.Printf("B\n")
	case C:
		fmt.Printf("C\n")
	}

}
