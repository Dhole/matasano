package main

import (
	"./modes"
	//"./set2"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	input_b64, err := ioutil.ReadFile("data/10.txt")
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

	//plaintext := set2.AESCBCDecrypt(ciphertext, iv, key)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	plaintext := make([]byte, len(ciphertext))
	cbc_dec := modes.NewCBCDecrypter(block, iv)
	cbc_dec.CryptBlocks(plaintext, ciphertext)

	fmt.Print(string(plaintext))
}
