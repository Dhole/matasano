package main

import (
	//"./hmac"
	//"./sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func main() {
	msg := []byte(url.QueryEscape("The cake is a lie"))
	//msg := []byte("HELLO")
	hmac := make([]byte, 20)

	url_base := "http://localhost:3000/test?file=%s&signature=%s"
	for i := 0; i < len(hmac); i++ {
		var wg sync.WaitGroup
		//var times [0x100]float64
		times := make([]float64, 0x100)
		// One goroutine for each byte. This is pretty aggressive with
		// all the gouroutines at the same time. Delays can appear
		// leading to bad results.
		for j := 0; j < 256; j++ {
			hmac[i] = byte(j)
			u := fmt.Sprintf(url_base, msg, hex.EncodeToString(hmac))
			wg.Add(1)
			go func(times []float64, u string, j int) {
				start := time.Now()
				resp, _ := http.Get(u)
				times[j] += float64(time.Since(start))
				//fmt.Println(times[j])
				resp.Body.Close()
				wg.Done()
				return
			}(times, u, j)
		}
		wg.Wait()
		//fmt.Println(times)

		max_time := 0.0
		max_j := 0
		for j, t := range times {
			if t > max_time {
				max_time = t
				max_j = j
			}
		}
		//fmt.Println(times)
		hmac[i] = byte(max_j)
		fmt.Println(hex.EncodeToString(hmac))
	}
	fmt.Println(fmt.Sprintf(url_base, msg, hex.EncodeToString(hmac)))
	/*
		sha1 := sha1.New()

		key := []byte("")
		msg := []byte("")
		fmt.Println(hex.EncodeToString(hmac.Calc(sha1, msg, key)))

		key = []byte("key")
		msg = []byte("The quick brown fox jumps over the lazy dog")
		fmt.Println(hex.EncodeToString(hmac.Calc(sha1, msg, key)))

		key = []byte("trimarans")
		msg = []byte("foo")
		fmt.Println(hex.EncodeToString(hmac.Calc(sha1, msg, key)))
	*/
}
