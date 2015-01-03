package main

import (
	"./modes"
	"crypto/aes"
	"encoding/base64"
	"fmt"
)

func main() {
	ciphtxt_b64 := "L77na/nrFsKvynd6HzOoG7GHTLXsTVu9qvY/2syLXzhPweyyMTJULu/6/kXX0KSvoOLSFQ=="
	ciphtxt, _ := base64.StdEncoding.DecodeString(ciphtxt_b64)

	key := []byte("YELLOW SUBMARINE")
	nonce := uint64(0)

	clrtxt := make([]byte, len(ciphtxt))

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	cipher := modes.NewCTR(block, nonce)
	cipher.XORKeyStream(clrtxt, ciphtxt)
	fmt.Println("Clear:", string(clrtxt))
}
