package utils

import (
	"errors"
	"fmt"
	"math"
	"sort"
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

func EnhancedTempCalculation(listTemp []float64) (float64, float64, int) {
	var listTempProm []float64
	var subList1 []float64
	var subList2 []float64
	var tam, tamSubList, removed int
	var sum, media, desvEst, mediana1, mediana2, minAccept, maxAccept, sumAvg, avg float64
	tam = len(listTemp)
	tamSubList = tam / 2
	//saca la media inicial
	for i := 0; i < tam; i++ {
		sum += listTemp[i]
	}
	media = sum / float64(tam)
	fmt.Println("La media es : ", media)
	//saca la desviacion estandar
	for j := 0; j < tam; j++ {
		desvEst += math.Pow(listTemp[j]-media, 2)
	}
	desvEst = math.Sqrt(desvEst / float64(tam))
	fmt.Println("La desviacion estandar es : ", desvEst)
	sort.Float64s(listTemp)

	//dividir la lista en dos para sacar la mediana
	for i := 0; i < tamSubList; i++ {
		subList1 = append(subList1, listTemp[i])
		subList2 = append(subList2, listTemp[tamSubList+i])
	}
	//sacar mediana
	if tamSubList%2 == 0 {
		mediana1 = (subList1[(tamSubList/2)-1] + subList1[(tamSubList/2)]) / 2
		mediana2 = (subList2[(tamSubList/2)-1] + subList2[(tamSubList/2)]) / 2
	} else {
		//impar
		mediana1 = subList1[(tamSubList/2)-1]
		mediana2 = subList2[(tamSubList/2)-1]
	}
	// sacar los valores minimos y maximos aceptables con el Test de Tukey
	minAccept = mediana1 - 1.5*(mediana2-mediana1)
	maxAccept = mediana2 + 1.5*(mediana2-mediana1)
	fmt.Println("El minimo aceptable es : ", minAccept)
	fmt.Println("El maximo aceptable es : ", maxAccept)
	for i := 0; i < tam; i++ {
		if listTemp[i] > minAccept && listTemp[i] < maxAccept {
			listTempProm = append(listTempProm, listTemp[i])
			sumAvg += listTemp[i]
		} else {
			removed++
		}
	}
	avg = sumAvg / float64(len(listTempProm))
	fmt.Println("La temperatura es : ", avg)
	//saca la desviacion estandar despues de evaluar
	desvEst = 0
	for j := 0; j < len(listTempProm); j++ {
		desvEst += math.Pow(listTempProm[j]-avg, 2)
	}
	desvEst = math.Sqrt(desvEst / float64(len(listTempProm)))
	fmt.Println("La desviacion estandar es : ", desvEst)
	return desvEst, avg, removed
}
