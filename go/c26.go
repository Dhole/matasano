package main

import (
	"./modes"
	"./set1"
	"bytes"
	"crypto/aes"
	crand "crypto/rand"
	"fmt"
	"strings"
)

var key []byte

func CTREncryptionOracle(src []byte) (dst []byte) {
	// Clean input
	src = bytes.Replace(src, []byte(";"), []byte(""), -1)
	src = bytes.Replace(src, []byte("="), []byte(""), -1)

	key = make([]byte, 16)
	_, err := crand.Read(key)
	if err != nil {
		panic(err)
	}
	nonce := uint64(0)

	before_str := "comment1=cooking%20MCs;userdata="
	after_str := ";comment2=%20like%20a%20pound%20of%20bacon"
	before := []byte(before_str)
	after := []byte(after_str)

	dst = make([]byte, len(before)+len(src)+len(after))
	copy(dst, before)
	copy(dst[len(before):], src)
	copy(dst[len(before)+len(src):], after)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipher := modes.NewCTR(block, nonce)
	cipher.XORKeyStream(dst, dst)
	return dst
}

func CTRDecryptionOracle(src []byte) (is_admin bool) {
	dst := make([]byte, len(src))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	nonce := uint64(0)

	cipher := modes.NewCTR(block, nonce)
	cipher.XORKeyStream(dst, src)

	fmt.Println("Decription:")
	fmt.Println(string(dst))
	return strings.Contains(string(dst), ";admin=true;")
}

func main() {
	before_str := "comment1=cooking%20MCs;userdata="
	//after_str := ";comment2=%20like%20a%20pound%20of%20bacon"
	goal := ";admin=true"
	//fill := make([]byte, len(goal))
	fill := bytes.Repeat([]byte{byte('A')}, len(goal))
	fmt.Println("Sent:")
	fmt.Println(string(fill))

	ciphtxt := CTREncryptionOracle(fill)

	set1.XorBytes(ciphtxt[len(before_str):],
		ciphtxt[len(before_str):], fill)
	set1.XorBytes(ciphtxt[len(before_str):],
		ciphtxt[len(before_str):], []byte(goal))

	if CTRDecryptionOracle(ciphtxt) {
		fmt.Println("Got Admin!")
	} else {
		fmt.Println("Failed")
	}
}
