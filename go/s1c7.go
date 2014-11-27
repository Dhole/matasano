package main

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {

	input_b64, err := ioutil.ReadFile("7.txt")
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

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		log.Fatal("Needs to be multiple of the blocksize! ",
			len(ciphertext))
	}

	plaintext := make([]byte, len(ciphertext))
	ptxt := plaintext
	ctxt := ciphertext
	for len(ctxt) > 0 {
		block.Decrypt(ptxt[:bs], ctxt[:bs])
		ptxt = ptxt[bs:]
		ctxt = ctxt[bs:]
	}
	fmt.Println(string(plaintext))
}
