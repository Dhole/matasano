package main

import (
	"./set1"
	"./set2"
	"fmt"
)

func main() {

	bs := 16
	before_str := "comment1=cooking%20MCs;userdata="
	after_str := ";comment2=%20like%20a%20pound%20of%20bacon"
	len_pre_pad := len(before_str) % bs
	fill := make([]byte, len_pre_pad+bs)
	goal := ";admin=true;"
	mod := make([]byte, bs)
	set1.XorBytes(mod, []byte(goal), []byte(after_str[:len(goal)]))

	ciphtxt := set2.AESEncryptionOracle4(fill)

	set1.XorBytes(ciphtxt[len(before_str)+len_pre_pad:],
		ciphtxt[len(before_str)+len_pre_pad:], mod)

	if set2.AESDecryptionOracle4(ciphtxt) {
		fmt.Println("Got Admin!")
	} else {
		fmt.Println("Failed")
	}
}
