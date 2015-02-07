package main

import (
	"./sha1"
	"bufio"
	"fmt"
	"hash"
	"math/rand"
	"os"
	"reflect"
	"time"
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

func MDPad(len_msg int) (pad []byte) {
	var len_pad int
	if len_msg%64 > 64-9 {
		len_pad = 64 + (len_msg % 64)
	} else {
		len_pad = 64 - (len_msg % 64)
	}
	pad = make([]byte, len_pad)
	pad[0] = 0x80
	len_msg *= 8 // length in bits
	for i := uint(0); i < 8; i++ {
		pad[i+uint(len_pad)-8] = byte(len_msg >> (64 - 8 - 8*i))
	}
	//fmt.Println(pad_msg)
	return pad
}

func newShaFromHash(sha_sum []byte, len uint64) hash.Hash {
	var new_h [5]uint32
	for i := 0; i < 5; i++ {
		for j := uint32(0); j < 4; j++ {
			new_h[i] |= uint32(sha_sum[i*4+int(j)]) << (8 * (3 - j))
		}
	}
	//fmt.Println("New H:", new_h)
	h := sha1.NewSpecialInit(new_h, len)
	return h
}

func genCookie() (cookie []byte, mac []byte) {
	// Choose random keyword
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}
	key = []byte(words[rand.Intn(len(words))])
	fmt.Println("Secret key:", string(key))

	cookie = []byte("comment1=cooking%20MCs;userdata=foo;comment2=%20like%20a%20pound%20of%20bacon")

	mac = genMac(cookie)
	return cookie, mac
}

func main() {
	rand.Seed(time.Now().UnixNano())
	cookie, mac := genCookie()
	fmt.Println("cookie:", string(cookie))
	fmt.Println("MAC:", mac)

	custom := []byte(";admin=true;")

	var cookie2 []byte
	var mac2 []byte
	len_key := 0
	max_len_key := 100
	for len_key = 0; len_key < max_len_key; len_key++ {
		pad_cookie := append(cookie, MDPad(len(cookie)+len_key)...)
		cookie2 = append(pad_cookie, custom...)
		h := newShaFromHash(mac, uint64(len(pad_cookie)+len_key))
		h.Write(custom)
		mac2 = h.Sum(nil)
		if checkMac(cookie2, mac2) {
			fmt.Println("Key length:", len_key)
			break
		}
	}
	fmt.Println("New cookie:", string(cookie2))
	fmt.Println("New MAC:", mac2)
	if checkMac(cookie2, mac2) {
		fmt.Println("Forged message has valid MAC!")
	} else {
		fmt.Println("Forged message has invalid MAC")
	}
}
