package fib

import "testing"

type testpair struct {
	input  int
	result int
}

var tests = []testpair{
	{0, 0},
	{1, 1},
	{2, 1},
	{3, 2},
	{4, 3},
	{5, 5},
	{6, 8},
	{7, 13},
	{10, 55},
	{15, 610},
	{20, 6765},
	{30, 832040},
	{40, 102334105},
	{50, 12586269025},
}

func TestFib(t *testing.T) {
	for _, pair := range tests {
		v := Fib(pair.input)
		if v != pair.result {
			t.Error("For", pair.input,
				"expected", pair.result,
				"got", v,
			)
		}
	}
}
