package sleep

import (
	// "fmt"
	// "strconv"
	"time"
)

func Sleep(numOfSeconds int) {

	//open a channel and block the current program
	<-time.After(time.Duration(numOfSeconds) * time.Second)
	// fmt.Println("Slept " + strconv.Itoa(int(numOfSeconds)) + " seconds")

}

// func main() {

// 	fmt.Println("start at " + strconv.Itoa(time.Now().Second()))
// 	Sleep(5)
// 	fmt.Println("start at " + strconv.Itoa(time.Now().Second()))

// }
