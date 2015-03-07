package main

import (
	//"./hmac"
	//"./sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

/*
UNFINISHED

I didn't finish this challenge but here are some comments about it:
First, special care must be taken to the concurrent http petitions. If they
work with goroutines and there are too many, the results will be very
unreliable! A small number of petitions should be made at the same time, so that
the timing values are accurate (and the goroutines don't stay waiting for too
long once they have received the http response). This is tricky.

Secondly, regarding to finding the mac byte one at a time, the idea would be
to do a number of petitions for every byte and take the timing results as
gausian populations. The results from one byte will be a different population
than the rest, this will correspond to the byte in the mac. To find the
different population (same distribution but with different mean), a statistical
test must be run. For example Mann-Whitney U test. When there are not enough
samples to find a different populations, more timings must be obtained until
a different population pops out.

*/

func main() {
	msg := []byte(url.QueryEscape("The cake is a lie"))
	//msg := []byte("HELLO")
	hmac := make([]byte, 20)

	iter := 100

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
				for k := 0; k < iter; k++ {
					start := time.Now()
					resp, _ := http.Get(u)
					times[j] += float64(time.Since(start))
					//fmt.Println(times[j])
					resp.Body.Close()
				}
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
	u := fmt.Sprintf(url_base, msg, hex.EncodeToString(hmac))
	resp, _ := http.Get(u)
	contents, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(string(contents))
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
