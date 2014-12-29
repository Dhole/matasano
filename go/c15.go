package main

import (
	"./pkcs7"
	"fmt"
)

func main() {
	pad_msgs := make([]string, 3)
	pad_msgs[0] = "ICE ICE BABY\x04\x04\x04\x04"
	pad_msgs[1] = "ICE ICE BABY\x05\x05\x05\x05"
	pad_msgs[2] = "ICE ICE BABY\x01\x02\x03\x04"

	var msg []byte
	var err error
	for i, pad_msg := range pad_msgs {
		msg, err = pkcs7.Unpad([]byte(pad_msg))
		if err == nil {
			fmt.Println("msg", i, "valid:")
			fmt.Println(string(msg))
		} else {
			fmt.Println("msg", i, "invalid!")
			fmt.Println(err)
		}

	}
}
