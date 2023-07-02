package xp

import "math"

func CalcExperience(length uint, diversity uint) uint {
	// f(x;y) = ((0.025 x^{1.25})/(y^{-0.5}))+1
	return uint(math.Floor(((0.025 * math.Pow(float64(length), 1.25)) / math.Pow(float64(diversity), -0.5)) + 1))
}

func CalcLevel(xp uint) uint {
	return uint(math.Floor(0.1 * math.Pow(float64(xp), 0.5)))
}
