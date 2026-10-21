package distr

import (
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Uniform реализует равномерный закон ИЗ2 на [a, b]
func Uniform(a, b float64) float64 {
	return a + rand.Float64()*(b-a)
}

// Exponential реализует экспоненциальный закон ПЗ1
func Exponential(lambda float64) float64 {
	u := rand.Float64()
	for u == 0 {
		u = rand.Float64()
	}
	return -math.Log(1.0-u) / lambda
}
