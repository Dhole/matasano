package main

import (
	"./rsa"
	"encoding/base64"
	"fmt"
	"math/big"
)

var e, d, n *big.Int

func oracleIsOdd(c *big.Int) bool {
	m := rsa.Crypt(c, d, n)
	if m.Bit(0) == 1 {
		return true
	} else {
		return false
	}
}

func main() {
	m_b64 := "VGhhdCdzIHdoeSBJIGZvdW5kIHlvdSBkb24ndCBwbGF5IGFyb3VuZCB3aXRoIHRoZSBGdW5reSBDb2xkIE1lZGluYQ=="

	m_bytes, err := base64.StdEncoding.DecodeString(m_b64)
	if err != nil {
		panic(err)
	}

	m := big.NewInt(0)
	m.SetBytes(m_bytes)
	//m.SetBytes([]byte("HELLO WORLD!"))

	e, d, n = rsa.GenKeyPair(1024)
	c := rsa.Crypt(m, e, n)

	c1 := big.NewInt(0).Set(c)
	aux := rsa.Crypt(big.NewInt(2), e, n)
	aux.Mod(aux, n)

	top, bottom := big.NewInt(0), big.NewInt(0)
	top.Set(n)

	frac := big.NewInt(0)
	pow := big.NewInt(1)
	for i := 0; i < 1024; i++ {
		c1.Mul(c1, aux)
		c1.Mod(c1, n)
		pow.Mul(pow, big.NewInt(2))
		frac.Div(n, pow)
		if oracleIsOdd(c1) {
			bottom.Add(bottom, frac)
		} else {
			top.Sub(top, frac)
		}
		fmt.Println(rsa.Decode(top))
		//fmt.Println(top.Bytes())
	}
}
