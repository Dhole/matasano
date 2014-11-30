package main

import (
	"./pkcs7"
	"fmt"
)

func main() {
	input_str := "YELLOR SUBMARINE"
	input := []byte(input_str)
	input_padded, _ := pkcs7.Pad(input, 20)
	fmt.Println(input_padded)
	input_unpadded, _ := pkcs7.Unpad(input_padded)
	fmt.Println(input_unpadded)
}
