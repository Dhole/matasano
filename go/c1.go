package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	//"io/ioutil"
	//"os"
)

func main() {
	input_str := "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	//input, _ := ioutil.ReadAll(os.Stdin)
	//input_str := string(input)

	raw, _ := hex.DecodeString(input_str)
	output_str := base64.StdEncoding.EncodeToString(raw)
	fmt.Println(output_str)
}
