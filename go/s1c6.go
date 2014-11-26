package main

import (
	"./hamming"
	"./set1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {

	input_b64, err := ioutil.ReadFile("6.txt")
	if err != nil {
		log.Fatal(err)
	}

	input_b64_str := string(input_b64)
	input, err := base64.StdEncoding.DecodeString(input_b64_str)
	if err != nil {
		log.Fatal(err)
	}

	//key := make([]byte, 1)

	min_dist := 8.0
	min_key_size := 2
	for i := 2; i < 40; i++ {
		key_size := i
		dist := 0.0
		rep := 100 / key_size
		for j := 0; j < rep; j++ {
			fst := input[key_size*j : key_size*(j+1)]
			snd := input[key_size*(j+1) : key_size*(j+2)]
			dist += float64(hamming.DistByteSlice(fst, snd))
		}
		dist /= float64(key_size * rep)
		if dist < min_dist {
			min_dist = dist
			min_key_size = key_size
		}
	}

	key := make([]byte, min_key_size)

	col_len := len(input) / min_key_size
	col := make([]byte, col_len)
	for i := 0; i < min_key_size; i++ {
		for j := 0; j < col_len; j++ {
			col[j] = input[min_key_size*j+i]
		}
		key[i] = set1.BestByteXored(col)
	}
	fmt.Println("key:", string(key))
	fmt.Println("text:\n", string(set1.XorKey(input, key)))
	//fmt.Println(min_key_size)
	/*
		for key_len := 2; key_len < 40; key_len++ {




		min_min_score := 10.0
		min_min_key := make([]byte, 1)
		var min_min_input []byte

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			input, _ := hex.DecodeString(scanner.Text())

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
			if min_score < min_min_score {
				min_min_score = min_score
				min_min_key = min_key
				min_min_input = input
			}
		}

		fmt.Println(string(set1.XorKey(min_min_input, min_min_key)))
	*/
}
