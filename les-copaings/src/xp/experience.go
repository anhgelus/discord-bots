package xp

import "math"

func CalcExperience(length uint, diversity uint) uint {
	// f(x;y) = ((0.025 x^{1.25})/(y^{-0.5}))+1
	return uint(math.Floor(((0.025 * math.Pow(float64(length), 1.25)) / math.Pow(float64(diversity), -0.5)) + 1))
}

func CalcExperienceFromVocal(length uint) uint {
	// f(x)=((0.25 x^{1.3})/(60000))+1
	return uint(math.Floor(((0.25 * math.Pow(float64(length), 1.3)) / 60000) + 1))
}

func CalcLevel(xp uint) uint {
	// f(x)=0.1 x^0.5
	return uint(math.Floor(0.1 * math.Pow(float64(xp), 0.5)))
}

func CalcXpForLevel(level uint) uint {
	// f(x)=0.1 x^0.5
	// f(x)/0.1 = x^0.5
	// (f(x)/0.1)^2 = x
	return uint(math.Floor(math.Pow(10*float64(level), 2)))
}
