package main

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
)

func pow_mod(a, b, n uint64) (res uint64) {
	if b == 0 {
		return 1
	}
	res = a
	for i := uint64(1); i < b; i++ {
		res = (res * b) % n
	}
	return res
}

func dh_small_nums() {
	p := uint64(37)
	g := uint64(5)

	// Alice
	a := uint64(rand.Intn(256)) % p
	A := pow_mod(g, a, p)

	// Bob
	b := uint64(rand.Intn(256)) % p
	B := pow_mod(g, b, p)

	// Alice and Bob share A and B

	// Alice
	s_a := pow_mod(B, a, p)

	// Bob
	s_b := pow_mod(A, b, p)

	fmt.Println("Alice - secret:", a, "public:", A)
	fmt.Println("Bob   - secret:", b, "public:", B)
	fmt.Println("Shared secret:", s_a, "=", s_b)
}

func dh_big_nums() {
	_p, _ := hex.DecodeString("ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552bb9ed529077096966d670c354e4abc9804f1746c08ca237327ffffffffffffffff")
	_g, _ := hex.DecodeString("02")

	var p big.Int
	var g big.Int
	p.SetBytes(_p)
	g.SetBytes(_g)
	if p.ProbablyPrime(10) {
		fmt.Println("p is probably prime")
	} else {
		fmt.Println("p is probably not prime!")
	}

	N := int64(math.Ceil(float64(p.BitLen()) / 8.0))

	// Alice
	var a big.Int
	var A big.Int
	_a := make([]byte, N)
	_, err := crand.Read(_a)
	if err != nil {
		panic(err)
	}
	a.SetBytes(_a)
	A.Exp(&g, &a, &p)

	// Bob
	var b big.Int
	var B big.Int
	_b := make([]byte, N)
	_, err = crand.Read(_b)
	if err != nil {
		panic(err)
	}
	b.SetBytes(_b)
	B.Exp(&g, &b, &p)

	// Alice and Bob share A and B

	// Alice
	var s_a big.Int
	s_a.Exp(&B, &a, &p)

	// Bob
	var s_b big.Int
	s_b.Exp(&A, &b, &p)

	fmt.Println("Shared secret:")
	fmt.Println(s_a.String())
	fmt.Println(s_b.String())
}

func main() {
	fmt.Println("\n\tSmall numbers\n")
	dh_small_nums()
	fmt.Println("\n\tBig numbers\n")
	dh_big_nums()
}
