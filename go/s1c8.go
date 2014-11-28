package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func IsECB(input []byte, bs int) (res bool) {
	n_blocks := len(input) / bs
	for i := 0; i < n_blocks-1; i++ {
		blk_a := input[bs*i : bs*(i+1)]
		for j := i + 1; j < n_blocks; j++ {
			blk_b := input[bs*j : bs*(j+1)]
			if bytes.Equal(blk_a, blk_b) {
				return true
			}
		}
	}
	return false
}

func main() {

	file, err := os.Open("8.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bs := 16
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input, _ := hex.DecodeString(scanner.Text())
		if IsECB(input, bs) {
			fmt.Println("ECB detected at:")
			fmt.Println(hex.EncodeToString(input))
		}
	}
}
