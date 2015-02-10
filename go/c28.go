package main

import (
	"./sha1"
	"fmt"
	"reflect"
)

var key []byte

func genMac(msg []byte) (mac []byte) {
	h := sha1.New()
	h.Write(key)
	h.Write(msg)
	mac = h.Sum(nil)
	return mac
}

func checkMac(msg, mac []byte) bool {
	h := sha1.New()
	h.Write(key)
	h.Write(msg)
	new_mac := h.Sum(nil)
	return reflect.DeepEqual(mac, new_mac)
}

func main() {
	key = []byte("top secret")

	msg := []byte("The cake is a lie")
	mac := genMac(msg)
	fmt.Println("MAC:", mac)
	if checkMac(msg, mac) {
		fmt.Println("Untouched message, right MAC")
	} else {
		fmt.Println("Untouched message, wrong MAC")
	}

}
