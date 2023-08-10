package utils

import "math"

func HoursOfUnix(unix int64) uint {
	return uint(math.Floor(float64(unix) / 1)) //(60 * 60)
}
