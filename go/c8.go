package main

import (
	"./set1"
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"os"
)

func main() {

	file, err := os.Open("data/8.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bs := 16
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input, _ := hex.DecodeString(scanner.Text())
		if set1.IsECB(input, bs) {
			fmt.Println("ECB detected at:")
			fmt.Println(hex.EncodeToString(input))
		}
	}
}
