package main

import (
	"fmt"
	"time"
)

func main() {
	Fib44 := make(chan int)
	Fib45 := make(chan int)
	go spinner(100 * time.Millisecond)
	go func(Fib44 chan int) {
		n := 44
		FibN := fib(n)
		fmt.Println(FibN)
		Fib44 <- FibN
	}(Fib44)
	go func(Fib45 chan int) {
		n1 := 45
		FibN1 := fib(n1)
		fmt.Println(FibN1)
		Fib45 <- FibN1
	}(Fib45)
	<-Fib44
	<-Fib45

}

func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
