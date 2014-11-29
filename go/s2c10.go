package main

import (
	"./set1"
	"crypto/aes"
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
	last_cblk := iv
	for len(ctxt) > 0 {
		block.Decrypt(ptxt[:bs], ctxt[:bs])
		copy(ptxt[:bs], set1.XorKey(ptxt[:bs], last_cblk))
		last_cblk = ctxt[:bs]
		ptxt = ptxt[bs:]
		ctxt = ctxt[bs:]
	}
	fmt.Println(string(plaintext))
}
