package main

import (
	"./set1"
	"./set2"
	"bytes"
	"fmt"
)

type oracle func([]byte) []byte

func FindOracleBlockSize(ora oracle) (blockSize int) {
	blockSize = 0
	msg := bytes.Repeat([]byte{byte('A')}, 0)
	out := ora(msg)
	len_fst := len(out)
	for i := 1; i < 256; i++ {
		msg = bytes.Repeat([]byte{byte('A')}, i)
		out := ora(msg)
		if len(out) > len_fst {
			blockSize = len(out) - len_fst
			break
		}
	}
	return blockSize
}

func IsOracleECB(ora oracle, blockSize int) (is_ECB bool) {
	input := make([]byte, 16*3)
	output := ora(input)
	return set1.IsECB(output, blockSize)
}

func main() {
	set2.SetOracle2Key()
	//input := make([]byte, 16*3)
	//output := set2.AESEncryptionOracle2(input)
	//fmt.Println(output)
	blockSize := FindOracleBlockSize(set2.AESEncryptionOracle2)
	if IsOracleECB(set2.AESEncryptionOracle2, blockSize) {
		fmt.Println("It's ECB")
	} else {
		fmt.Println("It's not ECB")
	}

}
