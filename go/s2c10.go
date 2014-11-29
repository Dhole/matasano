package main

import (
	"./set2"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	input_b64, err := ioutil.ReadFile("10.txt")
	if err != nil {
		log.Fatal(err)
	}

	input_b64_str := string(input_b64)
	input, err := base64.StdEncoding.DecodeString(input_b64_str)
	if err != nil {
		log.Fatal(err)
	}

	ciphertext := input

	key_str := "YELLOW SUBMARINE"
	key := []byte(key_str)
	iv := make([]byte, 16)

	plaintext := set2.AESCBCDecrypt(ciphertext, iv, key)

	fmt.Print(string(plaintext))
}
