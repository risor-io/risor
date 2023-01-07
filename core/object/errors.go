package object

import (
	"math"
)

func NewArgsError(fn string, takes, given int) *Error {
	return Errorf("type error: %s() takes exactly %d arguments (%d given)",
		fn, takes, given)
}

func NewArgsRangeError(fn string, takesMin, takesMax, given int) *Error {
	if math.Abs(float64(takesMax-takesMin)) <= 0.0001 {
		return Errorf("type error: %s() takes %d or %d arguments (%d given)",
			fn, takesMin, takesMax, given)
	}
	return Errorf("type error: %s() takes between %d and %d arguments (%d given)",
		fn, takesMin, takesMax, given)
}
