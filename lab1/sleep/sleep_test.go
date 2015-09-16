package sleep

import (
	"math"
	"testing"
	"time"
)

var tests = []int{
	//test set
	1, 2, 3, 4, 5, 6, 7, 8, 9, 8, 7, 6, 5, 4, 3, 2, 1,
}

//as program itself run with several hundred milliseconds, there might be deviation
const ALLOWED_DEVIATION = 1

func TestSleep(t *testing.T) {
	for _, input := range tests {
		testStart := time.Now().Unix()
		Sleep(input)
		testEnd := time.Now().Unix()
		gap := testEnd - testStart
		diff := int(gap) - input
		if math.Abs(float64(diff)) > ALLOWED_DEVIATION {
			t.Error("For", input,
				"expected", input,
				"got", gap,
			)
		}
	}
}
