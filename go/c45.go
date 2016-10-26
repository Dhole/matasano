package main

import (
	"./dsa"
	"./utils"
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/big"
)

func main() {

	p_str := "800000000000000089e1855218a0e7dac38136ffafa72eda7859f2171e25e65eac698c1702578b07dc2a1076da241c76c62d374d8389ea5aeffd3226a0530cc565f3bf6b50929139ebeac04f48c3c84afb796d61e5a4f9a8fda812ab59494232c7d2b4deb50aa18ee9e132bfa85ac4374d7f9091abc3d015efc871a584471bb1"
	q_str := "f4f47f05794b256174bba6e9b396a7707e563c5b"
	g_str := "5958c9d3898b224b12672c0b98e06c60df923cb8bc999d119458fef538b8fa4046c8db53039db620c094c9fa077ef389b5322a559946a71903f990f1f7e0e025e2d7f7cf494aff1a0470f5b64c36b625a097f1651fe775323556fe00b3608c887892878480e99041be601a62166ca6894bdd41a7054ec89f756ba9fc95302291"
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
	p := big.NewInt(0).SetBytes(_p)
	q := big.NewInt(0).SetBytes(_q)
	g := big.NewInt(0).SetBytes(_g)

	params := dsa.Params{sha1.New(), p, q, g}

	badParams1(&params)
	params.G = g
	badParams2(&params)
}

func BadSign(par *dsa.Params, keys *dsa.Keys, msg []byte) dsa.Signature {
	r := big.NewInt(0)
	s := big.NewInt(0)
	par.Hash.Reset()
	par.Hash.Write(msg)
	h_b := par.Hash.Sum(nil)
	h := big.NewInt(0).SetBytes(h_b)
	for {
		k, err := crand.Int(crand.Reader, par.Q)
		if err != nil {
			panic(err)
		}

		r.Exp(par.G, k, par.P)
		r.Mod(r, par.Q)
		//if r.Cmp(big.NewInt(0)) == 0 {
		//	continue
		//}
		//print_dec("r", r)

		k1 := big.NewInt(0).ModInverse(k, par.Q)
		s.Mul(keys.X, r)
		s.Add(h, s)
		s.Mul(k1, s)
		s.Mod(s, par.Q)
		//if s.Cmp(big.NewInt(0)) == 0 {
		//	continue
		//}

		break
	}
	return dsa.Signature{r, s}
}

func badParams1(params *dsa.Params) {

	keys := dsa.KeyGen(params)
	params.G = big.NewInt(0)

	msg_a := "I'll give you $10"
	msg_b := "I'll give you $9999"

	// With g = 0, r = 0, so the signature never finishes, because after
	// getting r = 0 it tries a new k again and again.  If the signature
	// algorithm proceeds we'd end up with a signature that always has r as
	// 0 and the verification will succeed for any message because since g =
	// 0, v = 0.
	//sig_a := dsa.Sign(params, &keys, []byte(msg_a))
	// We'll use a bad signature implementation that doesn't check if r is 0
	sig_a := BadSign(params, &keys, []byte(msg_a))

	utils.PrintHex("sig_a.r", sig_a.R)
	utils.PrintHex("sig_a.s", sig_a.S)

	if dsa.Verify(params, &sig_a, keys.Y, []byte(msg_a)) {
		fmt.Println("Valid signature for msg_a:", msg_a)
	}

	if dsa.Verify(params, &sig_a, keys.Y, []byte(msg_b)) {
		fmt.Println("Valid signature for msg_b:", msg_b)
	}

}

func badParams2(params *dsa.Params) {

	keys := dsa.KeyGen(params)
	utils.PrintHex("keys.x", keys.X)
	utils.PrintHex("keys.y", keys.Y)

	params.G.Add(params.P, big.NewInt(1))

	r := big.NewInt(0).Exp(keys.Y, big.NewInt(42), params.P)
	r.Mod(r, params.Q)
	s := big.NewInt(0).ModInverse(big.NewInt(42), params.Q)
	s.Mul(s, r)
	s.Mod(s, params.Q)
	master_sig := dsa.Signature{r, s}
	utils.PrintHex("master_sig.r", master_sig.R)
	utils.PrintHex("master_sig.s", master_sig.S)

	msg_a := "Hello, world"
	msg_b := "Goodbye, world"

	if dsa.Verify(params, &master_sig, keys.Y, []byte(msg_a)) {
		fmt.Println("Valid signature for msg_a:", msg_a)
	}

	if dsa.Verify(params, &master_sig, keys.Y, []byte(msg_b)) {
		fmt.Println("Valid signature for msg_b:", msg_b)
	}
}
