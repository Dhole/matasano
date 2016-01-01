package main

import (
	"./rsa"
	"bytes"
	"fmt"
	"math/big"
)

var bitlen = 512

func genSecrets() (c1, c2, c3, n1, n2, n3 *big.Int) {
	m := "Super secret message"
	m_b := rsa.Encode(m)
	e1, d1, n1 := rsa.GenKeyPair(bitlen)
	e2, d2, n2 := rsa.GenKeyPair(bitlen)
	e3, d3, n3 := rsa.GenKeyPair(bitlen)
	c1 = rsa.Crypt(m_b, e1, n1)
	c2 = rsa.Crypt(m_b, e2, n2)
	c3 = rsa.Crypt(m_b, e3, n3)
	fmt.Printf("m=%v\n", m_b)
	fmt.Printf("e1=%v\ne2=%v\ne3=%v\n", e1, e2, e3)
	fmt.Printf("d1=%v\nd2=%v\nd3=%v\n", d1, d2, d3)
	return c1, c2, c3, n1, n2, n3
}

func cubeRoot(n int) (r int) {
	min := 0
	max := 1024
	r = min
	inc := max - min
	for inc > 0 {
		q := r * r * r
		if q < n {
			r += inc
		} else if q > n {
			r -= inc
		} else {
			break
		}
		inc /= 2
	}
	return r
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
			fmt.Println("YAY!")
			break
		}
		inc.Div(inc, two)
	}
	return r
}

func main() {

	c1, c2, c3, n1, n2, n3 := genSecrets()
	/*
		c1.SetInt64(2)
		c2.SetInt64(3)
		c3.SetInt64(2)
		n1.SetInt64(3)
		n2.SetInt64(5)
		n3.SetInt64(7)
	*/

	fmt.Printf("c1=%v\nc2=%v\nc3=%v\n", c1, c2, c3)
	fmt.Printf("n1=%v\nn2=%v\nn3=%v\n", n1, n2, n3)

	ms1 := new(big.Int).Mul(n2, n3)
	ms2 := new(big.Int).Mul(n1, n3)
	ms3 := new(big.Int).Mul(n1, n2)
	fmt.Printf("ms1=%v\nms2=%v\nms3=%v\n", ms1, ms2, ms3)

	n123 := new(big.Int).Mul(n1, n2)
	n123.Mul(n123, n3)

	_, i1 := rsa.Egcd(n1, ms1)
	_, i2 := rsa.Egcd(n2, ms2)
	_, i3 := rsa.Egcd(n3, ms3)
	res := big.NewInt(0)
	acc := new(big.Int)

	acc.Mul(c1, ms1)
	acc.Mul(acc, i1)
	res.Add(res, acc)
	acc.Mul(c2, ms2)
	acc.Mul(acc, i2)
	res.Add(res, acc)
	acc.Mul(c3, ms3)
	acc.Mul(acc, i3)
	res.Add(res, acc)
	fmt.Printf("res=%v\n", res)

	res.Mod(res, n123)
	fmt.Printf("res mod n123=%v\n", res)

	m_b := cubeRootBig(res, bitlen)
	fmt.Printf("m_b=%v\n", m_b)
	m := rsa.Decode(m_b)
	fmt.Println(m)

	return
}
