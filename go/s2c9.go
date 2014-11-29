package main

import (
	"./set2"
	"fmt"
)

func main() {
	input_str := "YELLOR SUBMARINE"
	input := []byte(input_str)
	padded_input, _ := set2.PKCS7Padding(input, 20)
	fmt.Println(padded_input)
}
