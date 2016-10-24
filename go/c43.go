package main

import (
	"bytes"
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
)

type params struct {
	hash    hash.Hash
	p, q, g *big.Int
}

type key_pair struct {
	x, y *big.Int
}

type signature struct {
	r, s *big.Int
}

func parameter_gen(hash hash.Hash, L, N int) (*params, error) {
	if N > hash.Size()*8 {
		return nil, fmt.Errorf("N must be smaller or equal than " +
			"the hash size")
	}

	q, err := crand.Prime(crand.Reader, N)
	if err != nil {
		panic(err)
	}
	//print_dec("q", q)

	p := big.NewInt(0)
	rnd_buf := make([]byte, (L-N)/8)
	_, err = crand.Read(rnd_buf)
	if err != nil {
		panic(err)
	}
	r := big.NewInt(0).SetBytes(rnd_buf)
	//r := crand.Int(crand.Reader, L-N)
	// Make r even
	r.SetBit(r, 0, 0)
	//print_dec("r", r)
	// p := r*q + 1, where r is random
	p.Mul(r, q)
	p.Add(p, big.NewInt(1))
	//print_dec("r*q + 1", p)

	// We keep adding 2*q to p until it's a prime
	q2 := big.NewInt(0).Mul(big.NewInt(2), q)
	for {
		//print_dec("p1", p)
		if p.ProbablyPrime(80) {
			break
		}
		p.Add(p, q2)
	}
	//print_dec("p", p)

	// g = h^((p − 1)/q) mod p
	g := big.NewInt(0)
	p1 := big.NewInt(0)
	p1.Sub(p, big.NewInt(1))
	p1.Div(p1, q)
	// NOTE: The limit on h is only valid for N > 32
	for h := int64(2); h < 4294967296; h++ {
		if g.Exp(big.NewInt(h), p1, p).Cmp(big.NewInt(1)) != 0 {
			break
		}
	}
	//print_dec("g", g)

	return &params{hash, p, q, g}, nil
}

func key_gen(par *params) key_pair {
	x, err := crand.Int(crand.Reader, par.q)
	if err != nil {
		panic(err)
	}
	y := big.NewInt(0)
	y.Exp(par.g, x, par.p)

	return key_pair{x, y}
}

func sign(par *params, keys *key_pair, msg []byte) signature {
	r := big.NewInt(0)
	s := big.NewInt(0)
	par.hash.Reset()
	par.hash.Write(msg)
	h_b := par.hash.Sum(nil)
	h := big.NewInt(0).SetBytes(h_b)
	for {
		k, err := crand.Int(crand.Reader, par.q)
		if err != nil {
			panic(err)
		}

		r.Exp(par.g, k, par.p)
		r.Mod(r, par.q)
		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		//print_dec("r", r)

		k1 := big.NewInt(0).ModInverse(k, par.q)
		s.Mul(keys.x, r)
		s.Add(h, s)
		s.Mul(k1, s)
		s.Mod(s, par.q)
		if s.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		r1 := big.NewInt(0).ModInverse(r, par.q)
		x1 := x_from_k(r1, s, par.q, h, k)
		print_dec("x1", x1)

		break
	}
	return signature{r, s}
}

func verify(par *params, sig *signature, y *big.Int, msg []byte) bool {
	if sig.r.Cmp(big.NewInt(0)) != 1 && sig.r.Cmp(par.q) != -1 {
		return false
	}
	if sig.s.Cmp(big.NewInt(0)) != 1 && sig.s.Cmp(par.q) != -1 {
		return false
	}

	w := big.NewInt(0)
	u1 := big.NewInt(0)
	u2 := big.NewInt(0)
	v := big.NewInt(0)

	par.hash.Reset()
	par.hash.Write(msg)
	h_b := par.hash.Sum(nil)
	fmt.Printf("h = 0x%v\n", hex.EncodeToString(h_b))
	h := big.NewInt(0).SetBytes(h_b)

	w.ModInverse(sig.s, par.q)
	print_dec("w", w)

	u1.Mul(h, w)
	u1.Mod(u1, par.q)
	print_dec("u1", u1)

	u2.Mul(sig.r, w)
	u2.Mod(u2, par.q)
	print_dec("u2", u2)

	u1.Exp(par.g, u1, par.p)
	u2.Exp(y, u2, par.p)

	v.Mul(u1, u2)
	v.Mod(v, par.p)
	v.Mod(v, par.q)
	print_dec("v", v)
	print_dec("r", sig.r)
	print_hex("r", sig.r)

	if v.Cmp(sig.r) == 0 {
		return true
	}
	return false
}

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
	dsa_params, err := parameter_gen(sha1.New(), 64, 16)
	//dsa_params, err := parameter_gen(sha1.New(), 1024, 160)
	if err != nil {
		panic(err)
	}
	print_dec("p", dsa_params.p)
	print_dec("q", dsa_params.q)
	print_dec("g", dsa_params.g)
	dsa_keys := key_gen(dsa_params)
	print_dec("x", dsa_keys.x)
	print_dec("y", dsa_keys.y)

	msg := "OLA K ASE"
	msg_b := []byte(msg)
	sig := sign(dsa_params, &dsa_keys, msg_b)
	fmt.Printf("msg = \"%v\"\n", msg)
	print_dec("r", sig.r)
	print_dec("s", sig.s)
	ok := verify(dsa_params, &sig, dsa_keys.y, msg_b)
	fmt.Println("verify = ", ok)
	fmt.Println("End")

	test()
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

func test() {
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

	par := params{sha1.New(), p, q, g}
	sig := signature{r, s}
	check := verify(&par, &sig, y, []byte(msg))
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
