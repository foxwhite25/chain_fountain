package main

import "math"

func MinFromSlice(slice []float64) float64 {
	m := math.Inf(1)
	for _, value := range slice {
		if value < m {
			m = value
		}
	}
	return m
}
