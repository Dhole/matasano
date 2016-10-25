package main

import (
	"./dsa"
	"bytes"
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
	// FIPS 186-3 specifies L and N length pairs of
	// (1,024, 160), (2,048, 224), (2,048, 256), and (3,072, 256).
	// dsa_params, err := parameter_gen(sha1.New(), 1024, 160)
	dsa_params, err := dsa.ParameterGen(sha1.New(), 64, 16)
	//dsa_params, err := parameter_gen(sha1.New(), 1024, 160)
	if err != nil {
		panic(err)
	}
	print_dec("p", dsa_params.P)
	print_dec("q", dsa_params.Q)
	print_dec("g", dsa_params.G)
	dsa_keys := dsa.KeyGen(dsa_params)
	print_dec("x", dsa_keys.X)
	print_dec("y", dsa_keys.Y)

	msg := "OLA K ASE"
	msg_b := []byte(msg)
	sig := dsa.Sign(dsa_params, &dsa_keys, msg_b)
	fmt.Printf("msg = \"%v\"\n", msg)
	print_dec("r", sig.R)
	print_dec("s", sig.S)
	ok := dsa.Verify(dsa_params, &sig, dsa_keys.Y, msg_b)
	fmt.Println("verify = ", ok)
	fmt.Println("End")

	challenge()
}

func x_from_k(r1, s, q, h, k *big.Int) *big.Int {
	//r1 := big.NewInt(0).ModInverse(r, q)
	x1 := big.NewInt(0)
	x1.Mul(s, k)
	x1.Sub(x1, h)
	x1.Mul(x1, r1)
	x1.Mod(x1, q)

	return x1
}

func challenge() {
	p_str :=
		"800000000000000089e1855218a0e7dac38136ffafa72eda7859f2171e25e65eac698c1702578b07dc2a1076da241c76c62d374d8389ea5aeffd3226a0530cc565f3bf6b50929139ebeac04f48c3c84afb796d61e5a4f9a8fda812ab59494232c7d2b4deb50aa18ee9e132bfa85ac4374d7f9091abc3d015efc871a584471bb1"

	q_str := "f4f47f05794b256174bba6e9b396a7707e563c5b"

	g_str :=
		"5958c9d3898b224b12672c0b98e06c60df923cb8bc999d119458fef538b8fa4046c8db53039db620c094c9fa077ef389b5322a559946a71903f990f1f7e0e025e2d7f7cf494aff1a0470f5b64c36b625a097f1651fe775323556fe00b3608c887892878480e99041be601a62166ca6894bdd41a7054ec89f756ba9fc95302291"

	y_str :=
		"084ad4719d044495496a3201c8ff484feb45b962e7302e56a392aee4abab3e4bdebf2955b4736012f21a08084056b19bcd7fee56048e004e44984e2f411788efdc837a0d2e5abb7b555039fd243ac01f0fb2ed1dec568280ce678e931868d23eb095fde9d3779191b8c0299d6e07bbb283e6633451e535c45513b2d33c99ea17"

	r_str := "548099063082341131477253921760299949438196259240"
	s_str := "857042759984254168557880549501802188789837994940"

	msg := "For those that envy a MC it can be hazardous to your health\nSo be friendly, a matter of life and death, just like a etch-a-sketch\n"

	h_x_str := "0954edd5e0afe5542a4adf012611a91912a3ec16"

	_p, _ := hex.DecodeString(p_str)
	_q, _ := hex.DecodeString(q_str)
	_g, _ := hex.DecodeString(g_str)
	_y, _ := hex.DecodeString(y_str)
	h_x, _ := hex.DecodeString(h_x_str)
	fmt.Println("H(x) =", hex.EncodeToString(h_x))

	p := big.NewInt(0).SetBytes(_p)
	q := big.NewInt(0).SetBytes(_q)
	g := big.NewInt(0).SetBytes(_g)
	y := big.NewInt(0).SetBytes(_y)

	r, _ := big.NewInt(0).SetString(r_str, 10)
	s, _ := big.NewInt(0).SetString(s_str, 10)

	print_hex("p", p)
	print_hex("g", g)
	print_hex("q", q)
	print_hex("y", y)
	print_hex("s", s)
	print_hex("r", r)

	par := dsa.Params{sha1.New(), p, q, g}
	sig := dsa.Signature{r, s}
	check := dsa.Verify(&par, &sig, y, []byte(msg))
	fmt.Println("MC =", check)

	hash := sha1.New()
	hash.Reset()
	hash.Write([]byte(msg))
	h_b := hash.Sum(nil)
	fmt.Println("H(m) =", hex.EncodeToString(h_b))
	h := big.NewInt(0).SetBytes(h_b)

	r1 := big.NewInt(0)
	r1.ModInverse(r, q)
	for _k := int64(0); _k <= 65536; _k++ {
		k := big.NewInt(_k)

		x := x_from_k(r1, s, q, h, k)

		x_hex := hex.EncodeToString(x.Bytes())
		//fmt.Println(x_hex)
		//y := big.NewInt(0).Exp(g, x, p)
		hash.Reset()
		hash.Write([]byte(x_hex))
		h_x1 := hash.Sum(nil)
		if bytes.Compare(h_x, h_x1) == 0 {
			fmt.Println("Found the private key!")
			print_hex("x", x)
			break
		}
		if _k%int64(65536/100) == 0 {
			fmt.Println(int(float64(_k*100)/65536), "%")
		}
	}
}
