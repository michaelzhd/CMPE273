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
	{5, 5},
	{7, 7},
	{3, 3},
	{7, 7},
	{3, 3},
	{7, 7},
	{3, 3},
	{7, 7},
	{3, 3},
}

func TestSleep(t *testing.T) {
	for _, pair := range tests {
		v := pair.input
		testStart := time.Now().Unix()
		Sleep(v)
		testEnd := time.Now().Unix()
		gap := testEnd - testStart
		diff := gap - int64(v)
		if diff > 1 && diff < -1 {
			t.Error("For", pair.input,
				"expected", pair.result,
				"got", gap,
			)
		}
	}
}
