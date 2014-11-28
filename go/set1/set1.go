package set1

import (
	"bytes"
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

func BestByteXored(input []byte) (k byte) {

	min_score := 10.0
	min_key := make([]byte, 1)
	key := make([]byte, 1)

	for i := 0x00; i < 0x100; i++ {
		key[0] = byte(i)
		xored := XorKey(input, key)
		hist := ComputeHist(xored)
		score := ComputeScore(Eng_freq, hist)
		if score < min_score {
			min_score = score
			copy(min_key, key)
		}
	}
	return min_key[0]
}

func IsECB(input []byte, bs int) (res bool) {
	n_blocks := len(input) / bs
	for i := 0; i < n_blocks-1; i++ {
		blk_a := input[bs*i : bs*(i+1)]
		for j := i + 1; j < n_blocks; j++ {
			blk_b := input[bs*j : bs*(j+1)]
			if bytes.Equal(blk_a, blk_b) {
				return true
			}
		}
	}
	return false
}
