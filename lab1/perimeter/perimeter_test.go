package perimeter

import (
	"fmt"
	"testing"
)

type testpair struct {
	input  []float64
	result float64
}

var tests = []testpair{
	//test for circle
	{[]float64{0, 0, 5}, 31.41592653589793},
	{[]float64{3, 3, 3}, 18.84955592153876},
	{[]float64{7, 8, 9}, 56.548667764616276},

	//test for rectangle
	{[]float64{0, 0, 10, 10}, 40},
	{[]float64{10, 10, 30, 30}, 80},
	{[]float64{5, 10, 15, 20}, 40},
}

func TestPerimeter(t *testing.T) {
	for _, pair := range tests {
		//obtain number of arguments to determine whether it's Circle or Rectangle
		argsNum := len(pair.input)

		//declare a shape
		var shape Shape
		switch argsNum {
		//3 arguments is for Circle
		case 3:
			shape = Circle{pair.input[0], pair.input[1], pair.input[2]}
		//4 arguments is for Rectangle
		case 4:
			shape = Rectangle{pair.input[0], pair.input[1], pair.input[2], pair.input[3]}
		//if not 3 or 4 arguments, then it's not a valid input
		default:
			fmt.Println("Wrong input")
		}
		v := shape.Perimeter()
		if v != pair.result {
			t.Error("For", pair.input,
				"expected", pair.result,
				"got", v,
			)
		}
	}
}
