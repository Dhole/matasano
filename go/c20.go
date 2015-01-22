package main

import (
	"./modes"
	"./set1"
	"bufio"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("data/20.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var ciphtxts [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s, err := base64.StdEncoding.DecodeString(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		ciphtxts = append(ciphtxts, s)
	}

	key := []byte("YELLOW SUBMARINE")
	nonce := uint64(0)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cipher := modes.NewCTR(block, nonce)

	max_len := 0
	for i := range ciphtxts {
		cipher.XORKeyStream(ciphtxts[i], ciphtxts[i])
		if len(ciphtxts[i]) > max_len {
			max_len = len(ciphtxts[i])
		}
	}

	key_stream := make([]byte, max_len)

	col_len := len(ciphtxts)
	col := make([]byte, col_len)
	for i := 0; i < len(key_stream); i++ {
		ind := 0
		for j := 0; j < col_len; j++ {
			if len(ciphtxts[j]) > i {
				//fmt.Println("i", i, "j", j)
				//fmt.Println(ciphtxts[j])
				col[ind] = ciphtxts[j][i]
				ind++
			}
		}
		key_stream[i] = set1.BestByteXored(col[:ind])
	}
	fmt.Println("key stream:", key_stream)
	clrtxt := make([]byte, max_len)
	for i, c := range ciphtxts {
		fmt.Println("Cleartext", i, ":")
		set1.XorBytes(clrtxt, c, key_stream)
		fmt.Println(string(clrtxt[:len(c)]))
		fmt.Println()
	}
}
