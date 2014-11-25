package main

import (
	"./set1"
	"encoding/hex"
	"fmt"
)

func main() {

	input_str := `Burning 'em, if you ain't quick and nimble
I go crazy when I hear a cymbal`
	key_str := "ICE"
	input := []byte(input_str)
	key := []byte(key_str)

	output := set1.XorKey(input, key)

	fmt.Println(hex.EncodeToString(output))
}
