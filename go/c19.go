package main

import (
	"./modes"
	"./set1"
	"crypto/aes"
	"encoding/base64"
	"fmt"
)

var clrtxts_b64 = []string{
	"SSBoYXZlIG1ldCB0aGVtIGF0IGNsb3NlIG9mIGRheQ==",
	"Q29taW5nIHdpdGggdml2aWQgZmFjZXM=",
	"RnJvbSBjb3VudGVyIG9yIGRlc2sgYW1vbmcgZ3JleQ==",
	"RWlnaHRlZW50aC1jZW50dXJ5IGhvdXNlcy4=",
	"SSBoYXZlIHBhc3NlZCB3aXRoIGEgbm9kIG9mIHRoZSBoZWFk",
	"T3IgcG9saXRlIG1lYW5pbmdsZXNzIHdvcmRzLA==",
	"T3IgaGF2ZSBsaW5nZXJlZCBhd2hpbGUgYW5kIHNhaWQ=",
	"UG9saXRlIG1lYW5pbmdsZXNzIHdvcmRzLA==",
	"QW5kIHRob3VnaHQgYmVmb3JlIEkgaGFkIGRvbmU=",
	"T2YgYSBtb2NraW5nIHRhbGUgb3IgYSBnaWJl",
	"VG8gcGxlYXNlIGEgY29tcGFuaW9u",
	"QXJvdW5kIHRoZSBmaXJlIGF0IHRoZSBjbHViLA==",
	"QmVpbmcgY2VydGFpbiB0aGF0IHRoZXkgYW5kIEk=",
	"QnV0IGxpdmVkIHdoZXJlIG1vdGxleSBpcyB3b3JuOg==",
	"QWxsIGNoYW5nZWQsIGNoYW5nZWQgdXR0ZXJseTo=",
	"QSB0ZXJyaWJsZSBiZWF1dHkgaXMgYm9ybi4=",
	"VGhhdCB3b21hbidzIGRheXMgd2VyZSBzcGVudA==",
	"SW4gaWdub3JhbnQgZ29vZCB3aWxsLA==",
	"SGVyIG5pZ2h0cyBpbiBhcmd1bWVudA==",
	"VW50aWwgaGVyIHZvaWNlIGdyZXcgc2hyaWxsLg==",
	"V2hhdCB2b2ljZSBtb3JlIHN3ZWV0IHRoYW4gaGVycw==",
	"V2hlbiB5b3VuZyBhbmQgYmVhdXRpZnVsLA==",
	"U2hlIHJvZGUgdG8gaGFycmllcnM/",
	"VGhpcyBtYW4gaGFkIGtlcHQgYSBzY2hvb2w=",
	"QW5kIHJvZGUgb3VyIHdpbmdlZCBob3JzZS4=",
	"VGhpcyBvdGhlciBoaXMgaGVscGVyIGFuZCBmcmllbmQ=",
	"V2FzIGNvbWluZyBpbnRvIGhpcyBmb3JjZTs=",
	"SGUgbWlnaHQgaGF2ZSB3b24gZmFtZSBpbiB0aGUgZW5kLA==",
	"U28gc2Vuc2l0aXZlIGhpcyBuYXR1cmUgc2VlbWVkLA==",
	"U28gZGFyaW5nIGFuZCBzd2VldCBoaXMgdGhvdWdodC4=",
	"VGhpcyBvdGhlciBtYW4gSSBoYWQgZHJlYW1lZA==",
	"QSBkcnVua2VuLCB2YWluLWdsb3Jpb3VzIGxvdXQu",
	"SGUgaGFkIGRvbmUgbW9zdCBiaXR0ZXIgd3Jvbmc=",
	"VG8gc29tZSB3aG8gYXJlIG5lYXIgbXkgaGVhcnQs",
	"WWV0IEkgbnVtYmVyIGhpbSBpbiB0aGUgc29uZzs=",
	"SGUsIHRvbywgaGFzIHJlc2lnbmVkIGhpcyBwYXJ0",
	"SW4gdGhlIGNhc3VhbCBjb21lZHk7",
	"SGUsIHRvbywgaGFzIGJlZW4gY2hhbmdlZCBpbiBoaXMgdHVybiw=",
	"VHJhbnNmb3JtZWQgdXR0ZXJseTo=",
	"QSB0ZXJyaWJsZSBiZWF1dHkgaXMgYm9ybi4=",
}

func main() {
	key := []byte("YELLOW SUBMARINE")
	nonce := uint64(0)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cipher := modes.NewCTR(block, nonce)

	ciphtxts := make([][]byte, len(clrtxts_b64))
	max_len := 0
	for i := range ciphtxts {
		ciphtxts[i], _ = base64.StdEncoding.DecodeString(clrtxts_b64[i])
		cipher.XORKeyStream(ciphtxts[i], ciphtxts[i])
		if len(ciphtxts[i]) > max_len {
			max_len = len(ciphtxts[i])
		}
	}

	key_stream := make([]byte, max_len)

	// We find sentences starting with the 3 same chars, and guess what can
	// they be
	/*
		for i, c1 := range ciphtxts {
			for j, c2 := range ciphtxts {
				if c1[0] == c2[0] && c1[1] == c2[1] && c1[2] == c2[2] && i != j {
					fmt.Println("match", i, j)
				}
			}
		}
	*/
	//n := 23
	set1.XorBytes(key_stream, ciphtxts[0], []byte("I "))
	set1.XorBytes(key_stream[2:], ciphtxts[2][2:], []byte("om"))
	set1.XorBytes(key_stream[4:], ciphtxts[1][4:], []byte("ng"))
	set1.XorBytes(key_stream[6:], ciphtxts[3][6:], []byte("en"))
	set1.XorBytes(key_stream[8:], ciphtxts[26][8:], []byte("ng"))
	set1.XorBytes(key_stream[10:], ciphtxts[28][10:], []byte("ve"))
	set1.XorBytes(key_stream[12:], ciphtxts[29][12:], []byte("d "))
	set1.XorBytes(key_stream[14:], ciphtxts[3][14:], []byte("tury"))
	set1.XorBytes(key_stream[18:], ciphtxts[30][18:], []byte("ad"))
	set1.XorBytes(key_stream[20:], ciphtxts[11][20:], []byte("he "))
	set1.XorBytes(key_stream[23:], ciphtxts[21][23:], []byte("l,"))
	set1.XorBytes(key_stream[25:], ciphtxts[5][25:], []byte("ds"))
	set1.XorBytes(key_stream[27:], ciphtxts[8][27:], []byte("ne"))
	set1.XorBytes(key_stream[29:], ciphtxts[0][29:], []byte("ay"))
	set1.XorBytes(key_stream[31:], ciphtxts[27][31:], []byte("nd,"))
	set1.XorBytes(key_stream[34:], ciphtxts[4][34:], []byte("ad"))
	set1.XorBytes(key_stream[36:], ciphtxts[37][36:], []byte("n,"))
	// We found that 23, 25 and 30 begin with The

	for i := 0; i < len(ciphtxts); i++ {
		set1.XorBytes(ciphtxts[i], ciphtxts[i], key_stream)
		fmt.Println(i, string(ciphtxts[i]))
	}

	/*
		col_len := len(ciphtxts)
		col := make([]byte, col_len)
		for i := 0; i < len(key_stream); i++ {
			ind := 0
			for j := 0; j < col_len; j++ {
				if len(ciphtxts[j]) > i {
					//fmt.Println("i", i, "j", j)
					//fmt.Println(ciphtxts[j])
					col[ind] = ciphtxts[j][i]
					ind++
				}
			}
			key_stream[i] = set1.BestByteXored(col[:ind])
		}
		fmt.Println("key stream:", key_stream)
		clrtxt := make([]byte, max_len)
		for i, c := range ciphtxts {
			fmt.Println("Cleartext", i, ":")
			set1.XorBytes(clrtxt, c, key_stream)
			fmt.Println(string(clrtxt[:len(c)]))
			fmt.Println()
		}
	*/
}
