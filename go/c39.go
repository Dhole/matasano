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

func genKeyPair(bit_len int) (e, d, n *big.Int) {
	e, d, n = new(big.Int), new(big.Int), new(big.Int)

	var et big.Int
	e.SetInt64(3)

	for {
		p, err := crand.Prime(crand.Reader, bit_len)
		if err != nil {
			panic(err)
		}

		q, err := crand.Prime(crand.Reader, bit_len)
		if err != nil {
			panic(err)
		}

		if p.Cmp(q) == 0 {
			continue
		}

		n.Mul(p, q)

		et.Mul(new(big.Int).Sub(p, big.NewInt(1)),
			new(big.Int).Sub(q, big.NewInt(1)))

		if new(big.Int).Rem(&et, e).Cmp(big.NewInt(0)) != 0 {
			fmt.Printf("p = %d, q = %d, et = %d\n", p, q, &et)
			break
		}
	}

	d = invModBig(e, &et)

	return e, d, n
}

func crypt(m, e, n *big.Int) (c *big.Int) {
	c = new(big.Int)
	c.Exp(m, e, n)

	return c
}

func encode(m string) (m0 *big.Int) {
	m0 = new(big.Int)
	m0.SetBytes([]byte(m))
	return m0
}

func decode(m0 *big.Int) (m string) {
	m = string(m0.Bytes()[:])
	return m
}

func testInvModBig() {
	e := big.NewInt(3)
	primes := []int64{3, 5, 7, 11, 13, 17, 23}
	for _, p := range primes {
		fmt.Println("3,", p, invmod(3, int(p)))
		fmt.Println("3,", p, invModBig(e, big.NewInt(p)))
		fmt.Println("3,", p, new(big.Int).ModInverse(e, big.NewInt(p)))
		fmt.Println("")
	}
	fmt.Println("17, 3120,", invmod(17, 3120))
	fmt.Println("17, 3120,", invModBig(big.NewInt(17), big.NewInt(3120)))
	fmt.Println("17, 3120,", new(big.Int).ModInverse(big.NewInt(17), big.NewInt(3120)))
	fmt.Println("")
}

func main() {
	e, d, n := genKeyPair(512)
	fmt.Printf("e = %d, d = %d, n = %d\n", e, d, n)

	m := "Welcome home"
	m0 := encode(m)
	c := crypt(m0, e, n)
	m1 := decode(crypt(c, d, n))
	fmt.Printf("m0 = %s, c = %d, m1 = %s\n", m, c, m1)

}
