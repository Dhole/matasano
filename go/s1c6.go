package main

import (
	//"./set1"
	"./hamming"
	//"encoding/base64"
	"fmt"
	//"io/ioutil"
	//"log"
)


func main() {

	/*

	input_b64, err := ioutil.ReadFile("6.txt")
	if err != nil {
		log.Fatal(err)
	}

	input_b64_str := string(input_b64)
	input, err := base64.StdEncoding.DecodeString(input_b64_str)
	if err != nil {
		log.Fatal(err)
	}

	key := make([]byte, 1)

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
