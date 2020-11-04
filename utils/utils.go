package utils

import (
	"errors"
	"math"
)

// StandardDeviation calcula el desvío estándar de una serie de datos
func StandardDeviation(data []float64) float64 {
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

func Min(values []float64) (min float64, e error) {
	if len(values) == 0 {
		return 0, errors.New("cannot detect a minimum value in an empty slice")
	}

	min = values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}

	return min, nil
}

func Max(values []float64) (max float64, e error) {
	if len(values) == 0 {
		return 0, errors.New("cannot detect a maximum value in an empty slice")
	}

	max = values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}

	return max, nil
}