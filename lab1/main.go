package main

import (
	"fmt"
	// fib "lab1/fib"
	p "lab1/perimeter"
	// sleep "lab1/sleep"
)

func main() {
	// fmt.Println(fib.Fib(50))
	// sleep.Sleep(5)
	fmt.Println("done.")
	c := Circle(0, 0, 5)
	fmt.Println(c.Perimeter())
}
