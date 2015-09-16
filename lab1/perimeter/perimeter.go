package perimeter

import (
	// "fmt"
	"math"
)

//declare Shape interface
type Shape interface {
	Perimeter() float64
}

//declare Circle struct
type Circle struct {
	//x, y: the coordinate of Circle center
	// r: the radius of the Circle
	x, y, r float64
}

//declare Rectangle struct
type Rectangle struct {
	// x1,y1: the coordinate of up left vertex
	// x2,y2: the coordinate of down right vertex
	x1, y1, x2, y2 float64
}

//calculating the distance between two points
func distance(x1, y1, x2, y2 float64) float64 {
	a := x2 - x1
	b := y2 - y1
	return math.Sqrt(a*a + b*b)
}

//implement Shape interface for Circle struct
func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.r
}

//implement Shape interface for Rectangle struct
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
