package main

import (
	"./set1"
	"encoding/hex"
	"fmt"
)

func main() {

	input_str := "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	input, _ := hex.DecodeString(input_str)

	key := make([]byte, 1)
	min_score := 10.0
	min_key := make([]byte, 1)

	for i := 0x00; i < 0x100; i++ {
		key[0] = byte(i)
		xored := set1.XorKey(input, key)
		hist := set1.ComputeHist(xored)
		score := set1.ComputeScore(set1.Eng_freq, hist)
		if score < min_score {
			min_score = score
			copy(min_key, key)
		}
	}
	fmt.Println(string(set1.XorKey(input, min_key)))
}
