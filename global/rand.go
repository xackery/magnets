package global

import "math/rand"

func Rand(min, max int) int {
	if max < min {
		max = min
	}
	return rand.Intn(max-min+1) + min
}
