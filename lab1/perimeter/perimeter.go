package perimeter

import (
	// "fmt"
	"math"
)

type Shape interface {
	Perimeter() float64
}

type Circle struct {
	x, y, r float64
}

type Rectangle struct {
	x1, y1, x2, y2 float64
}

//calculating the distance between two point
func distance(x1, y1, x2, y2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(a*a + b*b)
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.r
}

func (rct Rectangle) Perimeter() float64 {
	length := distance(rct.x1, rct.y1, rct.x1, rct.y2)
	width := distance(rct.x1, rct.y1, rct.x2, rct.y1)
	return 2 * (length + width)

}

// func main() {

// 	c := Circle{x: 7, y: 8, r: 9}
// 	fmt.Println(c.Perimeter())

// 	rct := Rectangle{5, 10, 15, 20}
// 	fmt.Println(rct.Perimeter())
// }
