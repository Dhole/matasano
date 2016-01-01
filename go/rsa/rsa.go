package rsa

import (
	crand "crypto/rand"
	"math/big"
)

func InvModBig(a, n *big.Int) *big.Int {
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

func Egcd(a *big.Int, b *big.Int) (*big.Int, *big.Int) {
	if b.Cmp(big.NewInt(0)) == 0 {
		return big.NewInt(1), big.NewInt(0)
	}
	q := new(big.Int)
	r := new(big.Int)
	q.Quo(a, b)
	r.Mod(a, b)
	s, t := Egcd(b, r)
	return t, new(big.Int).Sub(s, new(big.Int).Mul(q, t))
}

func GenKeyPair(bit_len int) (e, d, n *big.Int) {
	e, d, n = new(big.Int), new(big.Int), new(big.Int)

	var et big.Int
	e.SetInt64(3)

	for {
		p, err := crand.Prime(crand.Reader, bit_len/2)
		if err != nil {
			panic(err)
		}

		q, err := crand.Prime(crand.Reader, bit_len/2)
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
			break
		}
	}

	d = InvModBig(e, &et)

	return e, d, n
}

func Crypt(m, e, n *big.Int) (c *big.Int) {
	c = new(big.Int)
	c.Exp(m, e, n)

	return c
}

func Encode(m string) (m0 *big.Int) {
	m0 = new(big.Int)
	m0.SetBytes([]byte(m))
	return m0
}

func Decode(m0 *big.Int) (m string) {
	m = string(m0.Bytes()[:])
	return m
}
