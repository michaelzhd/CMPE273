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

	//test for rectangular
	{[]float64{0, 0, 10, 10}, 40},
	{[]float64{10, 10, 30, 30}, 80},
	{[]float64{5, 10, 15, 20}, 40},
}

func TestPerimeter(t *testing.T) {
	for _, pair := range tests {
		argsNum := len(pair.input)
		switch argsNum {
		case 3:
			shape = new Circle{pair.input[0], pair.input[1], pair.input[2]};
		case 4:
			shape = new Rectangular{pair.input[0], pair.input[1], pair.input[2], pair.input[3]};
		default:
			fmt.Println("Wrong input");
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
