package fib

import (
	"fmt"
	// "strconv"
)

//prepare a vector to store calculated fibonnacci results
var fibStorage = make([]int, 100)

func Fib(n int) int {

	//deal with illegal input
	if n < 0 {
		fmt.Println("Wrong input!")
		return 0
	}

	//recursion escape condition
	if n == 0 || n == 1 {
		fibStorage[n] = n
		return n
	} else if fibStorage[n-1] != 0 {
		fibStorage[n] = fibStorage[n-1] + fibStorage[n-2]
		return fibStorage[n]
	} else {
		return Fib(n-1) + Fib(n-2)
	}
}

// func main() {

// 	fmt.Println(Fib(50))

// 	// for i := 1; i < 45; i++ {
// 	// 	fmt.Println(strconv.Itoa(i) + " --" + strconv.Itoa(Fib(i)))
// 	// }

// }
