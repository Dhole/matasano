package main

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
)

func invmod(a, n int) int {
	t0, r0, t1, r1 := 0, n, 1, a
	quotient := 0
	for r1 != 0 {
		quotient = r0 / r1
		t0, t1 = t1, t0-quotient*t1
		r0, r1 = r1, r0-quotient*r1
	}
	if r0 > 1 {
		return -1
	}
	if t0 < 0 {
		t0 = t0 + n
	}
	return t0
}

func invModBig(a, n *big.Int) *big.Int {
	t0 := big.NewInt(0)
	r0 := big.NewInt(0).Set(n)
	t1 := big.NewInt(1)
	r1 := big.NewInt(0).Set(a)
	quotient := big.NewInt(0)
	tmp := big.NewInt(0)
	zero := big.NewInt(0)
	for r1.Cmp(zero) != 0 {
		quotient.Set(tmp.Quo(r0, r1))
		tmp.Set(t1)
		t1.Set(t0.Sub(t0, big.NewInt(0).Mul(quotient, t1)))
		t0.Set(tmp)
		tmp.Set(r1)
		r1.Set(r0.Sub(r0, big.NewInt(0).Mul(quotient, r1)))
		r0.Set(tmp)
	}
	if r0.Cmp(big.NewInt(1)) == 1 {
		return big.NewInt(-1)
	}
	if t0.Cmp(zero) == -1 {
		t0 = t0.Add(t0, n)
	}
	return t0
}

func main() {
	bit_len := 8

	p, err := crand.Prime(crand.Reader, bit_len)
	if err != nil {
		panic(err)
	}

	q, err := crand.Prime(crand.Reader, bit_len)
	if err != nil {
		panic(err)
	}

	var n big.Int
	n.Mul(p, q)

	var et big.Int
	et.Mul(big.NewInt(0).Sub(p, big.NewInt(1)),
		big.NewInt(0).Sub(q, big.NewInt(1)))

	e := big.NewInt(3)

	fmt.Println("FIN", e)

	primes := []int64{3, 5, 7, 11, 13, 17, 23}
	for _, p := range primes {
		fmt.Println("3,", p, invModBig(big.NewInt(3), big.NewInt(p)))
		fmt.Println("3,", p, big.NewInt(0).ModInverse(big.NewInt(3), big.NewInt(p)))
		fmt.Println("")
	}

}
