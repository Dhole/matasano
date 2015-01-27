package main

import (
	"./modes"
	"./set1"
	"crypto/aes"
	crand "crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
)

var key []byte
var nonce uint64

func setKey() {
	key = make([]byte, 16)
	_, err := crand.Read(key)
	if err != nil {
		panic(err)
	}
}

func encrypt() (ciphtxt []byte) {
	input, err := ioutil.ReadFile("data/25_dec.txt")
	if err != nil {
		log.Fatal(err)
	}

	clrtxt := input
	ciphtxt = make([]byte, len(clrtxt))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	nonce = 1337
	cipher := modes.NewCTR(block, nonce)
	cipher.XORKeyStream(ciphtxt, clrtxt)

	return ciphtxt
}

func edit(ciphtxt []byte, newtxt []byte, offset int) {
	if len(ciphtxt)-offset < len(newtxt) {
		panic("New text too big")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	bs := uint64(block.BlockSize())
	ini_block := uint64(offset) / bs

	cipher := modes.NewCTR2(block, nonce, ini_block)
	cipher.XORKeyStream(ciphtxt[ini_block*bs:offset+len(newtxt)],
		ciphtxt[ini_block*bs:offset+len(newtxt)])

	copy(ciphtxt[offset:], newtxt)

	cipher = modes.NewCTR2(block, nonce, ini_block)
	cipher.XORKeyStream(ciphtxt[ini_block*bs:offset+len(newtxt)],
		ciphtxt[ini_block*bs:offset+len(newtxt)])
}

func main() {
	setKey()
	ciphtxt := encrypt()
	ciphtxt_copy := make([]byte, len(ciphtxt))
	copy(ciphtxt_copy, ciphtxt)
	newtxt := make([]byte, len(ciphtxt))
	edit(ciphtxt, newtxt, 0)
	set1.XorBytes(ciphtxt, ciphtxt_copy, ciphtxt)
	fmt.Println("Decription:")
	fmt.Println(string(ciphtxt))
}
