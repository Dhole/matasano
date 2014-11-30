package main

import (
	"./set2"
	"fmt"
)

func main() {
	set2.SetOracle2Key()
	input := make([]byte, 16*3)
	output := set2.AESEncryptionOracle2(input)
	fmt.Println(output)
}
