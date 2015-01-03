package main

import (
	"./modes"
	"./pkcs7"
	"crypto/aes"
	crand "crypto/rand"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"
)

type pad_oracle func([]byte, []byte) bool

var RndStrings = []string{
	"MDAwMDAwTm93IHRoYXQgdGhlIHBhcnR5IGlzIGp1bXBpbmc=",
	"MDAwMDAxV2l0aCB0aGUgYmFzcyBraWNrZWQgaW4gYW5kIHRoZSBWZWdhJ3MgYXJlIHB1bXBpbic=",
	"MDAwMDAyUXVpY2sgdG8gdGhlIHBvaW50LCB0byB0aGUgcG9pbnQsIG5vIGZha2luZw==",
	"MDAwMDAzQ29va2luZyBNQydzIGxpa2UgYSBwb3VuZCBvZiBiYWNvbg==",
	"MDAwMDA0QnVybmluZyAnZW0sIGlmIHlvdSBhaW4ndCBxdWljayBhbmQgbmltYmxl",
	"MDAwMDA1SSBnbyBjcmF6eSB3aGVuIEkgaGVhciBhIGN5bWJhbA==",
	"MDAwMDA2QW5kIGEgaGlnaCBoYXQgd2l0aCBhIHNvdXBlZCB1cCB0ZW1wbw==",
	"MDAwMDA3SSdtIG9uIGEgcm9sbCwgaXQncyB0aW1lIHRvIGdvIHNvbG8=",
	"MDAwMDA4b2xsaW4nIGluIG15IGZpdmUgcG9pbnQgb2g=",
	"MDAwMDA5aXRoIG15IHJhZy10b3AgZG93biBzbyBteSBoYWlyIGNhbiBibG93",
}

var blockSize = 16

var key []byte

func GenKey() {
	key = make([]byte, blockSize)
	_, err := crand.Read(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("key set")
}

func EncCBC() (dst []byte, iv []byte) {
	iv = make([]byte, blockSize)
	_, err := crand.Read(iv)
	if err != nil {
		panic(err)
	}

	str := RndStrings[rand.Intn(len(RndStrings))]
	dst, _ = base64.StdEncoding.DecodeString(str)
	//dst = []byte(str)
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
	return dst, iv
}

func PaddingOracle(src []byte, iv []byte) (valid bool) {
	dst := make([]byte, len(src))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	encrypter := modes.NewCBCDecrypter(block, iv)
	encrypter.CryptBlocks(dst, src)
	dst, err = pkcs7.Unpad(dst)
	if err == nil {
		return true
	} else {
		return false
	}
}

func GetBlock(ora pad_oracle, blk []byte, prev_blk []byte) (clr_blk []byte) {
	clr_blk = make([]byte, blockSize)
	iv := make([]byte, blockSize)

	i := blockSize - 1
	j := 0
	// First byte is special
	// If the cleartext appears to be ... 02 XX,
	// The Oracle would return true with both padding of 01 and 02
	for j = 0; j < 256; j++ {
		pad := byte(1)
		iv[i] = byte(j)
		if ora(blk, iv) {
			iv[i-1] = 128
			if ora(blk, iv) {
				clr_blk[i] = prev_blk[i] ^ byte(j) ^ pad
				break
			}
		}
	}
	for i = blockSize - 2; i >= 0; i-- {
		pad := byte(blockSize - i)
		for j = i + 1; j < blockSize; j++ {
			iv[j] = prev_blk[j] ^ clr_blk[j] ^ pad
		}
		for j = 0; j < 256; j++ {
			iv[i] = byte(j)
			if ora(blk, iv) {
				clr_blk[i] = prev_blk[i] ^ byte(j) ^ pad
				break
			}
		}
		if j == 256 {
			panic("Failed getting byte")
		}
	}
	return clr_blk
}

func GetClear(ora pad_oracle, ciphtxt, iv []byte) (clrtxt []byte) {
	clrtxt = make([]byte, len(ciphtxt))
	copy(clrtxt, GetBlock(ora, ciphtxt[:blockSize], iv))
	for i := 1; i < len(ciphtxt)/blockSize; i++ {
		start := blockSize * i
		end := blockSize * (i + 1)
		copy(clrtxt[start:end], GetBlock(ora, ciphtxt[start:end],
			ciphtxt[start-blockSize:end-blockSize]))
	}
	clrtxt, _ = pkcs7.Unpad(clrtxt)
	return clrtxt
}

func main() {
	rand.Seed(time.Now().UnixNano())
	GenKey()
	ciphtxt, iv := EncCBC()
	//blk0 := GetBlock(PaddingOracle, ciphtxt[:blockSize], iv)
	//fmt.Println(string(blk0))
	clrtxt := GetClear(PaddingOracle, ciphtxt, iv)
	fmt.Println("Secret:", string(clrtxt))
}
