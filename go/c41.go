package main

import (
	"./rsa"
	"./sha1"
	crand "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"time"
)

var bitlen = 1024
var e, d, n *big.Int
var seen_cs [][]byte

func init() {
	e, d, n = rsa.GenKeyPair(bitlen)
}

func server_dec(c *big.Int) (p *big.Int, err error) {
	h := sha1.New()
	h.Write(c.Bytes())
	h_sum := h.Sum(nil)
	for _, seen_c := range seen_cs {
		if reflect.DeepEqual(seen_c, h_sum) {
			return p, errors.New("Ciphertext already seen")
		}
	}
	seen_cs = append(seen_cs, h_sum)
	p = rsa.Crypt(c, d, n)
	return p, nil
}

func user_do() (c *big.Int) {
	social := "555-55-5555"
	p_s := fmt.Sprintf("{time: %d, social: %s}", time.Now().Unix(), social)
	p := rsa.Encode(p_s)
	c = rsa.Crypt(p, e, n)
	p1, _ := server_dec(c)
	if p.Cmp(p1) != 0 {
		log.Fatal("Decrypted message doesn't match with plain text.")
	}
	return c
}

func main() {
	c := user_do()
	s, err := crand.Int(crand.Reader, n)
	if err != nil {
		log.Fatal(err)
	}
	c1 := rsa.Crypt(s, e, n)
	c1.Mul(c1, c)
	c1.Mod(c1, n)

	p1, err := server_dec(c1)
	if err != nil {
		log.Fatal(err)
	}

	p := new(big.Int).Mul(p1, rsa.InvModBig(s, n))
	p.Mod(p, n)
	fmt.Println(rsa.Decode(p))
}
