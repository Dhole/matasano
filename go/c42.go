package main

import (
	"./rsa"
	"./sha1"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"
)

var ASN1_sha1 = []byte{0x30, 0x21, 0x30, 0x09, 0x06, 0x05, 0x2b, 0x0e, 0x03, 0x02, 0x1a, 0x05, 0x00, 0x04, 0x14}

var bitlen = 1024
var e, d, n *big.Int

func verify(m, sig []byte) (bool, error) {
	pkcs1 := rsa.Crypt(new(big.Int).SetBytes(sig), e, n).Bytes()
	fmt.Println("veryfing pkcs1 =", pkcs1)
	if len(pkcs1) < 5 {
		return false, errors.New("Signature too short")
	}
	if pkcs1[0] != 0x01 {
		return false, errors.New("Padding error")
	}
	pkcs1 = pkcs1[1:]
	for len(pkcs1) > 0 {
		if pkcs1[0] != 0xff {
			break
		}
		pkcs1 = pkcs1[1:]
	}
	if pkcs1[0] != 0x00 {
		return false, errors.New("Padding error")
	}
	pkcs1 = pkcs1[1:]
	i := 0
	for ; i < len(ASN1_sha1); i++ {
		if pkcs1[0] != ASN1_sha1[i] {
			return false, errors.New("Unrecognized hash function")
		}
		pkcs1 = pkcs1[1:]
	}
	h := sha1.New()
	h.Write(m)
	h_sum := h.Sum(nil)
	if len(pkcs1) < len(h_sum) {
		return false, errors.New("Hash too short")
	}
	for i = 0; i < len(h_sum); i++ {
		if pkcs1[i] != h_sum[i] {
			log.Printf("computed  hash =", h_sum)
			log.Printf("signature hash =", pkcs1)
			return false, errors.New("Hash doesn't match")
		}
	}
	return true, nil
}

func init() {
	e, d, n = rsa.GenKeyPair(bitlen)
}

func user_do() {
	m := []byte("I don't authorize the transfer 1234")
	h := sha1.New()
	h.Write(m)
	h_sum := h.Sum(nil)
	t := append(ASN1_sha1, h_sum...)
	fill := bitlen/8 - 2 - 1 - len(t)
	pkcs1 := append([]byte{0x00, 0x01}, bytes.Repeat([]byte{0xFF}, fill)...)
	pkcs1 = append(pkcs1, 0x00)
	pkcs1 = append(pkcs1, t...)
	fmt.Println("user pkcs1 =", pkcs1)
	sig := rsa.Crypt(new(big.Int).SetBytes(pkcs1), d, n)
	ok, err := verify(m, sig.Bytes())
	if !ok {
		log.Fatal(err)
	} else {
		log.Println("User got a verified signature")
	}
}

func forge_sig() {
	m := []byte("I authorize the transfer 1234")
	h := sha1.New()
	h.Write(m)
	h_sum := h.Sum(nil)
	t := append([]byte{0x00}, ASN1_sha1...)
	t = append(t, h_sum...)

	d_len := (1 + len(ASN1_sha1) + 20) * 8
	var disp int
	if bitlen/3 < d_len {
		disp = (3 * 8) / 8
	} else {
		disp = (bitlen/3 - d_len) / 8
	}

	pkcs1 := make([]byte, bitlen/8)
	pkcs1[0] = 0x00
	pkcs1[1] = 0x01
	for i := 2; i < disp; i++ {
		pkcs1[i] = 0xFF
	}
	for i := 0; i < len(t); i++ {
		pkcs1[i+disp] = t[i]
	}
	fmt.Println("forged pkcs1 =", pkcs1)
	sig := cubeRootBig(new(big.Int).SetBytes(pkcs1), bitlen)
	sig_3 := new(big.Int).Exp(sig, big.NewInt(3), n)
	fmt.Println("sig**3         =", sig_3.Bytes())
	if reflect.DeepEqual(sig_3.Bytes()[:disp+d_len/8-1], pkcs1[1:disp+d_len/8]) {
		fmt.Println("Forged signature looks good")
	} else {
		fmt.Println("Forged signature doesn't looks valid!")
	}

	ok, err := verify(m, sig.Bytes())
	if !ok {
		log.Fatal(err)
	} else {
		log.Println("Forger got a verified signature")
	}
}

func cubeRootBig(n *big.Int, n_bits int) (r *big.Int) {
	min := big.NewInt(0)
	max := new(big.Int).SetBytes(bytes.Repeat([]byte{0xFF}, n_bits/8))
	max.Add(max, big.NewInt(1))
	r = new(big.Int).Set(min)
	inc := new(big.Int).Sub(max, min)
	q := new(big.Int)
	zero := big.NewInt(0)
	two := big.NewInt(2)
	for inc.Cmp(zero) == 1 {
		q.Mul(r, r)
		q.Mul(q, r)
		if q.Cmp(n) == -1 {
			r.Add(r, inc)
		} else if q.Cmp(n) == 1 {
			r.Sub(r, inc)
		} else {
			break
		}
		inc.Div(inc, two)
	}
	// if r**3 < n**3 add 1 to r so that the displaced hash in the PKCS#1
	// is not modified
	if q.Cmp(n) == 1 {
		r.Add(r, big.NewInt(1))
	}
	return r
}

func main() {
	user_do()
	forge_sig()
	fmt.Println("FIN")
}
