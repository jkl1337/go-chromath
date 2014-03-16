package testnum

import (
	"math"
	"testing"
)

func Round(n float64, digits int) float64 {
	f := 1.0
	i := digits
	if i < 0 { i = -i }
	for i-- ; i >= 0 ; i-- {
		f = f*10.0
	}
	if digits < 0 {
		n /= f
	} else  {
		n *= f
	}

	if n >= 0.0 {
		n = math.Floor(n + 0.5)
	} else {
		n = math.Ceil(n - 0.5)
	}

	if digits < 0 {
		n *= f
	} else {
		n /= f
	}
	return n
}

func AlmostEqual(t *testing.T, a float64, b float64, places int, msg string) {
	if Round(math.Abs(a - b), places) != 0 {
		t.Error(msg, a, b)
	}
}
