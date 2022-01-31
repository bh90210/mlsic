package util

import (
	"strconv"

	"github.com/mb-14/gomarkov"
)

func Intervallic(s []int) []int {
	var i []int
loop:
	for k := range s {
		switch {
		case k+1 >= len(s):
			break loop

		case s[k] < s[k+1]:
			i = append(i, s[k+1]-s[k])

		case s[k] > s[k+1]:
			i = append(i, s[k]-s[k+1])
		}
	}
	return i
}

func Normilazion(val, min, max float64) float64 {
	return (val - min) * (max - min)
}

func Markov(input []float64, min, max float64) []float64 {
	chain := gomarkov.NewChain(1)

	var sitoa []string
	for _, v := range input {
		sitoa = append(sitoa, strconv.FormatFloat(v, 'E', 0, 32))
	}

	chain.Add(sitoa)

	var ret []float64
	for k, v := range sitoa {
		if k == 0 {
			continue
		}
		i, _ := chain.TransitionProbability(v, []string{sitoa[k-1]})
		ret = append(ret, Normilazion(i, min, max))
	}

	return ret
}

func Int2Float(input []int) []float64 {
	var floats []float64
	for _, v := range input {
		floats = append(floats, float64(v))
	}
	return floats
}

func Float2Int(input []float64) []int {
	var ints []int
	for _, v := range input {
		ints = append(ints, int(v))
	}
	return ints
}
