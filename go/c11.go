package main

import (
	"./set2"
	"fmt"
	"math/rand"
	"reflect"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cbc_count := 0
	ecb_count := 0
	input := make([]byte, 16*3)

	for i := 0; i < 100; i++ {
		output := set2.AESEncryptionOracle(input)
		if reflect.DeepEqual(output[16*1:16*2], output[16*2:16*3]) {
			ecb_count++
		} else {
			cbc_count++
		}
	}

	fmt.Println("ECB:", ecb_count, "\nCBC:", cbc_count)
}
