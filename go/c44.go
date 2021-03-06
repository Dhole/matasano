package main

import (
	"./dsa"
	"bytes"
	"io/ioutil"
	"strings"
	//crand "crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	//"hash"
	"math/big"
)

func print_hex(pre string, v *big.Int) {
	fmt.Printf("%v = 0x%v\n", pre, hex.EncodeToString(v.Bytes()))
}

func print_dec(pre string, v *big.Int) {
	fmt.Printf("%v = %v\n", pre, v.String())
}

func main() {
	challenge()
}

func challenge() {
	p_str :=
		"800000000000000089e1855218a0e7dac38136ffafa72eda7859f2171e25e65eac698c1702578b07dc2a1076da241c76c62d374d8389ea5aeffd3226a0530cc565f3bf6b50929139ebeac04f48c3c84afb796d61e5a4f9a8fda812ab59494232c7d2b4deb50aa18ee9e132bfa85ac4374d7f9091abc3d015efc871a584471bb1"
	q_str := "f4f47f05794b256174bba6e9b396a7707e563c5b"
	g_str :=
		"5958c9d3898b224b12672c0b98e06c60df923cb8bc999d119458fef538b8fa4046c8db53039db620c094c9fa077ef389b5322a559946a71903f990f1f7e0e025e2d7f7cf494aff1a0470f5b64c36b625a097f1651fe775323556fe00b3608c887892878480e99041be601a62166ca6894bdd41a7054ec89f756ba9fc95302291"
	y_str := "2d026f4bf30195ede3a088da85e398ef869611d0f68f0713d51c9c1a3a26c95105d915e2d8cdf26d056b86b8a7b85519b1c23cc3ecdc6062650462e3063bd179c2a6581519f674a61f1d89a1fff27171ebc1b93d4dc57bceb7ae2430f98a6a4d83d8279ee65d71c1203d2c96d65ebbf7cce9d32971c3de5084cce04a2e147821"

	h_x_str := "ca8f6f7c66fa362d40760d135b763eb8527d3d52"

	_p, err := hex.DecodeString(p_str)
	if err != nil {
		panic(err)
	}
	_q, err := hex.DecodeString(q_str)
	if err != nil {
		panic(err)
	}
	_g, err := hex.DecodeString(g_str)
	if err != nil {
		panic(err)
	}
	_y, err := hex.DecodeString(y_str)
	if err != nil {
		panic(err)
	}
	h_x, _ := hex.DecodeString(h_x_str)
	if err != nil {
		panic(err)
	}
	p := big.NewInt(0).SetBytes(_p)
	q := big.NewInt(0).SetBytes(_q)
	g := big.NewInt(0).SetBytes(_g)
	y := big.NewInt(0).SetBytes(_y)
	print_hex("p", p)
	print_hex("g", g)
	print_hex("q", q)
	print_hex("y", y)
	fmt.Println("H(x) =", hex.EncodeToString(h_x))

	content, err := ioutil.ReadFile("44.txt")
	if err != nil {
		panic(err)
	}
	lines := strings.Split(string(content), "\n")
	n := len(lines) / 4
	msgs := make([]string, n)
	ss := make([]*big.Int, n)
	rs := make([]*big.Int, n)
	ms := make([]*big.Int, n)

	for i := 0; i < n; i++ {
		msgs[i] = lines[i*4+0][5:]
		r, ok := big.NewInt(0).SetString(lines[i*4+2][3:], 10)
		if ok != true {
			panic("")
		}
		s, ok := big.NewInt(0).SetString(lines[i*4+1][3:], 10)
		if ok != true {
			fmt.Println(lines[i*4+1][3:])
			panic("")
		}
		_m, err := hex.DecodeString(lines[i*4+3][3:])
		if err != nil {
			fmt.Println(lines[i*4+3][3:])
			panic(err)
		}
		m := big.NewInt(0).SetBytes(_m)
		ss[i] = s
		rs[i] = r
		ms[i] = m
	}

	// Find repeated r (signature generated by the same k)
	a, b := -1, -1
DoubleLoop:
	for i := 0; i < n-1; i++ {
		//print_hex("i", rs[i])
		for j := i + 1; j < n; j++ {
			//print_hex("  j", rs[j])
			if bytes.Compare(rs[i].Bytes(), rs[j].Bytes()) == 0 {
				a, b = i, j
				fmt.Println("Signature", a, "and", b,
					"used the same k")
				break DoubleLoop
			}
		}
	}

	// Get k from the two signatures that used the same k:
	// k = (h1-h2)/(s1-s2) mod q

	k := big.NewInt(0).Sub(ms[a], ms[b])
	div := big.NewInt(0).Sub(ss[a], ss[b])
	div.ModInverse(div, q)
	k.Mul(k, div)
	k.Mod(k, q)

	// Now we recover x from k
	r1 := big.NewInt(0).ModInverse(rs[a], q)
	x := dsa.XFromK(r1, ss[a], q, ms[a], k)

	x_hex := hex.EncodeToString(x.Bytes())
	hash := sha1.New()
	hash.Reset()
	hash.Write([]byte(x_hex))
	h_x1 := hash.Sum(nil)
	if bytes.Compare(h_x, h_x1) == 0 {
		fmt.Println("Found the private key!")
		print_hex("x", x)
	} else {
		fmt.Println("Private key not found :(")
	}
}
