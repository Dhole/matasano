package set1

import (
	"math"
)

var Eng_freq = map[byte]float64{
	'a': 0.08167,
	'b': 0.01492,
	'c': 0.02782,
	'd': 0.04253,
	'e': 0.12702,
	'f': 0.02228,
	'g': 0.02015,
	'h': 0.06094,
	'i': 0.06966,
	'j': 0.00153,
	'k': 0.00772,
	'l': 0.04025,
	'm': 0.02406,
	'n': 0.06749,
	'o': 0.07507,
	'p': 0.01929,
	'q': 0.00095,
	'r': 0.05987,
	's': 0.06327,
	't': 0.09056,
	'u': 0.02758,
	'v': 0.00978,
	'w': 0.02360,
	'x': 0.00150,
	'y': 0.01974,
	'z': 0.00074,
}

func Xor(a, b []byte) (x []byte) {
	x = make([]byte, len(a))
	for i := range x {
		x[i] = a[i] ^ b[i]
	}
	return x
}

func XorKey(a, b []byte) (x []byte) {
	x = make([]byte, len(a))
	len_b := len(b)
	for i := range x {
		x[i] = a[i] ^ b[i%len_b]
	}
	return x
}

func ComputeHist(data []byte) (hist []float64) {
	hist = make([]float64, 0x100)
	len_data := len(data)
	for _, v := range data {
		hist[v]++
	}
	for i := range hist {
		hist[i] /= float64(len_data)
	}
	return hist
}

func ComputeScore(eng_freq map[byte]float64, hist []float64) (res float64) {
	for k, v := range eng_freq {
		res += math.Pow((hist[k]-v), 2.0) / (hist[k] + v)
	}
	return res / 2
}
