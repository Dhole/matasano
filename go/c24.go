package main

import (
	"./mt19937"
	crand "crypto/rand"
	//"./set1"
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func oracle(src []byte) (dst []byte) {
	rndPrefix := make([]byte, rand.Intn(32))
	_, err := crand.Read(rndPrefix)
	if err != nil {
		panic(err)
	}
	dst = make([]byte, len(rndPrefix)+len(src))
	copy(dst, rndPrefix)
	copy(dst[len(rndPrefix):], src)
	//fmt.Println(">", string(dst))
	streamCipher(dst, dst, uint16(rand.Intn(2^16)))
	return dst
}

func streamCipher(dst []byte, src []byte, key uint16) {
	m := mt19937.New()
	m.Init(uint32(key))
	for i := 0; i < len(src); i++ {
		dst[i] = src[i] ^ byte(m.ExtractNumber())
	}
}

func bruteForce(ciphtxt []byte) uint16 {
	msg := []byte(strings.Repeat("A", 14))
	clrtxt := make([]byte, len(ciphtxt))
	for i := 0; i < 2^16; i++ {
		streamCipher(clrtxt, ciphtxt, uint16(i))
		if bytes.Equal(msg, clrtxt[len(ciphtxt)-14:]) {
			return uint16(i)
		}
	}
	panic("No key found")
}

func genPassResetToken() (token []byte) {
	m := mt19937.New()
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	m.Init(uint32(time.Now().Unix()))
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	token = make([]byte, 16)
	for i := 0; i < len(token); i++ {
		token[i] = byte(m.ExtractNumber())
	}
	return token
}

func checkToken(token []byte, start_time, end_time int64) (bool, uint32) {
	m := mt19937.New()
	tmp_token := make([]byte, 16)
	for i := end_time; i > start_time-10; i-- {
		m.Init(uint32(i))
		for j := 0; j < len(token); j++ {
			tmp_token[j] = byte(m.ExtractNumber())
		}
		if bytes.Equal(token, tmp_token) {
			return true, uint32(i)
		}
	}
	return false, 0
}

func main() {
	rand.Seed(time.Now().UnixNano())
	/* DEMO
	clrtxt := []byte("Hello world!")
	streamCipher(clrtxt, clrtxt, 1337)
	fmt.Println("Ciphered:", clrtxt)
	streamCipher(clrtxt, clrtxt, 1337)
	fmt.Println("Deciphered:", string(clrtxt))
	*/
	// a
	fmt.Println("Part A")
	clrtxt := []byte(strings.Repeat("A", 14))
	ciphtxt := oracle(clrtxt)
	fmt.Println("Ciphertext:", ciphtxt)
	key := bruteForce(ciphtxt)
	fmt.Println("key:", key)
	streamCipher(ciphtxt, ciphtxt, key)
	fmt.Println("Cleartext:", string(ciphtxt))

	// b
	fmt.Println("\nPart B")
	start_time := time.Now().Unix()
	token := genPassResetToken()
	end_time := time.Now().Unix()
	fmt.Println("Start time:", start_time)
	fmt.Println("End time:  ", end_time)
	isToken, time := checkToken(token, start_time, end_time)
	if isToken {
		fmt.Println("Token", token, "is valid generated at", time)
	} else {
		fmt.Println("Token", token, "is not valid")
	}

}
