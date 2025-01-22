package common

import "math"

func Round[T float64 | float32](num T, pos int) T {
	result := math.Round(float64(num)*math.Pow10(pos)) / math.Pow10(pos)
	return T(result)
}
