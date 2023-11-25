package utils

import (
	"math"
	"sort"
)

type Distribution []float64

func (d *Distribution) GetMin() float64 {
	if len(*d) == 0 {
		return 0.
	}
	var min = math.Inf(1)
	for _, val := range *d {
		if val < min {
			min = val
		}
	}
	return min
}
func (d *Distribution) GetMax() float64 {
	if len(*d) == 0 {
		return 0.
	}
	var max = math.Inf(-1)
	for _, val := range *d {
		if val > max {
			max = val
		}
	}
	return max
}

func (d *Distribution) GetQuantileValue(quantileNumber uint) float64 {
	if !(0 <= quantileNumber && quantileNumber <= 4) {
		return 0.
	}
	min := d.GetMin()
	quartile := (d.GetMax() - min) / 4
	return min + (quartile * float64(quantileNumber))
}

// TrimMax sorts and trims max values if difference is too big, difference 0.25 means 25% between max values
func (d *Distribution) TrimMax(difference float64) {
	var scaling = 1. + difference
	sort.Slice(*d, func(i, j int) bool {
		return (*d)[i] > (*d)[j]
	})
	for len(*d) >= 2 && (*d)[1]*scaling < (*d)[0] {
		*d = (*d)[1:]
	}
}
