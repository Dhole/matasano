package main

import (
	"./modes"
	"./pkcs7"
	"./set1"
	"bytes"
	"crypto/aes"
	crand "crypto/rand"
	"errors"
	"fmt"
)

var key []byte
var iv []byte

func AESEncryptionOracle(src []byte) (dst []byte) {
	// Clean input
	src = bytes.Replace(src, []byte(";"), []byte(""), -1)
	src = bytes.Replace(src, []byte("="), []byte(""), -1)

	key = make([]byte, 16)
	_, err := crand.Read(key)
	if err != nil {
		panic(err)
	}
	iv = key
	fmt.Println("Key = IV:")
	fmt.Println(key)

	before_str := "comment1=cooking%20MCs;userdata="
	after_str := ";comment2=%20like%20a%20pound%20of%20bacon"
	before := []byte(before_str)
	after := []byte(after_str)

	dst = make([]byte, len(before)+len(src)+len(after))
	copy(dst, before)
	copy(dst[len(before):], src)
	copy(dst[len(before)+len(src):], after)
	dst, err = pkcs7.Pad(dst, 16)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	encrypter := modes.NewCBCEncrypter(block, iv)
	encrypter.CryptBlocks(dst, dst)
	return dst
}

func AESDecryptionOracle(src []byte) (dst []byte, err error) {
	dst = make([]byte, len(src))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	encrypter := modes.NewCBCDecrypter(block, iv)
	encrypter.CryptBlocks(dst, src)
	dst, _ = pkcs7.Unpad(dst)
	fmt.Println("Decription:")
	fmt.Println(string(dst))
	for _, v := range dst {
		if v >= 128 {
			return dst, errors.New("Not ASCII")
		}
	}
	return dst, nil
}

func main() {
	bs := 16
	//before_str := "comment1=cooking%20MCs;userdata="
	//after_str := ";comment2=%20like%20a%20pound%20of%20bacon"
	// length of before and after messages already makes 3 blocks
	fill := make([]byte, 0)

	ciphtxt := AESEncryptionOracle(fill)
	// Modify the ciphertext to get the custom plaintext
	copy(ciphtxt[bs:], make([]byte, bs))
	copy(ciphtxt[2*bs:], ciphtxt[:bs])
	clrtxt, err := AESDecryptionOracle(ciphtxt)
	if err != nil {
		fmt.Println(err)
	}
	// Recover key
	my_key := make([]byte, bs)
	set1.XorBytes(my_key, clrtxt[:bs], clrtxt[2*bs:])

	fmt.Println("Key = IV found:")
	fmt.Println(my_key)
}
