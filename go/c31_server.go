package main

import (
	"./hmac"
	"./sha1"
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/go-martini/martini"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var key []byte

func genKey() {
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
}

func insecure_compare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
		time.Sleep(50 * time.Millisecond)
	}
	return true
}

func test(r *http.Request) (status int, msg string) {
	qs := r.URL.Query()
	file := qs.Get("file")
	signature := qs.Get("signature")
	get_signature, err := hex.DecodeString(qs.Get("signature"))
	if err != nil {
		return 200, "Signature is not in hex\n"
	}
	//msg = r.FormValue("file")
	msg = fmt.Sprintf("File: %s\nSignature: %s", file, signature)

	sha1 := sha1.New()
	good_signature := hmac.Calc(sha1, []byte(file), key)
	if insecure_compare([]byte(get_signature), good_signature) {
		status = 200
		msg = fmt.Sprintf("%s\nValid Signature\n", msg)
	} else {
		status = 500
		msg = fmt.Sprintf("%s\nInvalid Signature\n", msg)
	}
	return status, msg
}

func main() {
	rand.Seed(time.Now().UnixNano())
	genKey()
	m := martini.Classic()
	m.Get("/test", test)
	m.Run()
}
