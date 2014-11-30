package main

import (
	"./set1"
	"encoding/hex"
	"fmt"
)

func main() {
	input_1_str := "1c0111001f010100061a024b53535009181c"
	input_2_str := "686974207468652062756c6c277320657965"

	input_1, _ := hex.DecodeString(input_1_str)
	input_2, _ := hex.DecodeString(input_2_str)
	output := set1.Xor(input_1, input_2)
	output_str := hex.EncodeToString(output)
	fmt.Println(output_str)
}
