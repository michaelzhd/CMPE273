package sleep

import (
	"testing"
	"time"
)

type testpair struct {
	input  int
	result int
}

var tests = []testpair{
	//test for circle
	{1, 1},
	{2, 2},
	{3, 3},
	{10, 10},
	// {17, 17},
	// {21, 21},
}

func TestSleep(t *testing.T) {
	for _, pair := range tests {
		v := pair.input
		testStart := time.Now().Second()
		Sleep(v)
		testEnd := time.Now().Second()
		gap := testEnd - testStart
		if v != int(gap) {
			t.Error("For", pair.input,
				"expected", pair.result,
				"got", gap,
			)
		}
	}
}
