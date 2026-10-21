package utils

import (
	"math"
	"math/rand"
)

// Uniform генерирует случайную величину по равномерному закону ИЗ2 на отрезке [a, b]
func Uniform(a, b float64) float64 {
	return a + rand.Float64()*(b-a)
}

// Exponential генерирует случайную величину по экспоненциальному закону ПЗ1
func Exponential(lambda float64) float64 {
	u := rand.Float64()
	for u == 0 {
		u = rand.Float64()
	}
	return -math.Log(1.0-u) / lambda
}

// MeanAndVariance рассчитывает среднее значение и дисперсию среза чисел
func MeanAndVariance(data []float64) (mean float64, variance float64) {
	n := len(data)
	if n == 0 {
		return 0, 0
	}
	var sum float64
	for _, v := range data {
		sum += v
	}
	mean = sum / float64(n)

	var sumSq float64
	for _, v := range data {
		sumSq += math.Pow(v-mean, 2)
	}
	variance = sumSq / float64(n)
	return mean, variance
}
