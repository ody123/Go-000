package main

import (
	"fmt"
	"time"
)

func sum(i int, c chan int) {
	res := 0
	for ; i > 0; i-- {
		res += 1
	}
	c <- res
}

func main() {
	for {
		a := time.After(0)
		b := time.After(0)
		time.Sleep(1)
		select {
		case <-b:
			fmt.Println("0")
		case <-a:
			fmt.Println("1")
		default:
			fmt.Println("2")
		}
		time.Sleep(1000)
	}
}

type Context interface {
	// Deadline returns the time when work done on behalf of this context
	// should be canceled. Deadline returns ok==false when no deadline is
	// set. Successive calls to Deadline return the same results.
	Deadline() (deadline time.Time, ok bool)
}
