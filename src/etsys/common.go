package etsys

import (
	"math"
	"math/rand"
)

func RndGauss(mu float64, sigma float64) float64 {
	return mu + (rand.Float64()-0.5)*sigma*2
}

func RndBool() bool {
	return rand.Float64() > 0.5
}

func Transpose(array [][]float64) [][]float64 {
	if len(array) == 0 {
		return array
	}
	transposed := make([][]float64, len(array[0]))
	for i, _ := range transposed {
		transposed[i] = make([]float64, len(array))
	}
	for i, _ := range array {
		for j, _ := range array[i] {
			transposed[j][i] = array[i][j]
		}
	}
	return transposed
}

func ArrayMin(array []float64) (min float64) {
	min = array[0]
	for _, v := range array[1:] {
		if v < min {
			min = v
		}
	}
	return
}

func ArrayMax(array []float64) (max float64) {
	max = array[0]
	for _, v := range array[1:] {
		if v > max {
			max = v
		}
	}
	return
}

func ArrayAverage(array []float64) (average float64) {
	average = 0
	for _, v := range array {
		average += v
	}
	average /= float64(len(array))
	return
}

func InArray(element float64, array []float64) bool {
	for _, el := range array {
		if el == element {
			return true
		}
	}
	return false
}

func EuclideanDistance(v1, v2 []float64) (dist float64) {
	if len(v1) != len(v2) {
		panic("dim")
	}
	dist = 0.0
	for i, _ := range v1 {
		dist += math.Pow(v1[i]-v2[i], 2)
	}
	dist = math.Sqrt(dist)
	return
}
