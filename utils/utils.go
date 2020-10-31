package utils

import (
	"math"
)

// StandardDeviation calcula el desvío estándar de una serie de datos
func StandardDeviation(data ...float64) float64 {
	var sum, mean, sd float64

	for _, d := range data {
		sum += d
	}

	mean = sum / float64(len(data))

	for _, d := range data {
		sd += math.Pow(d-mean, 2)
	}

	sd = math.Sqrt(sd / 10)

	return sd
}
