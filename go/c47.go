package main

import (
	"./rsa"
	"./utils"
	crand "crypto/rand"
	"fmt"
	"math/big"
)

type Pair struct {
	a, b *big.Int
}

var e, d, n *big.Int

func oracleIsPKCS(c *big.Int) bool {
	m := rsa.Crypt(c, d, n)
	m_b := m.Bytes()
	if len(m_b) == n.BitLen()/8-1 && m_b[0] == 0x02 {
		return true
	} else {
		return false
	}
}

func PKCS1Pad(msg []byte, n int) ([]byte, error) {
	if len(msg)+11 > n {
		return nil, fmt.Errorf("Msg is too long")
	}
	p_msg := make([]byte, n)
	p_msg[0] = 0x00
	p_msg[1] = 0x02
	_, err := crand.Read(p_msg[2 : n-len(msg)-1])
	if err != nil {
		panic(err)
	}
	p_msg[n-len(msg)-1] = 0x00
	copy(p_msg[n-len(msg):], msg)
	return p_msg, nil
}

func main() {
	m_str := "kick it, CC"
	p_m_b, err := PKCS1Pad([]byte(m_str), 256/8)
	if err != nil {
		panic(err)
	}
	m := big.NewInt(0).SetBytes([]byte(m_str))
	p_m := big.NewInt(0).SetBytes(p_m_b)
	utils.PrintHex("m", m)
	utils.PrintHex("p_m", p_m)

	e, d, n = rsa.GenKeyPair(256)
	c := rsa.Crypt(p_m, e, n)
	fmt.Println(oracleIsPKCS(c))
	utils.PrintHex("c", c)

	m1 := CCA(c)
}

func CCA(c, e, n *big.Int) *big.Int {
	B_b := make([]byte, n.BitLen/8-1)
	B_b[0] = 0x01
	B := big.NewInt(0).SetBytes(B_b)
	B_2 := big.NewInt(2)
	B_2.Mul(B_2, B)
	B_3 := big.NewInt(3)
	B_3.Mul(B_3, B)

	M := make([][]Pair, 0)
	s := make([]*big.Int, 0)

	// Step 1
	// s0 <- 1
	s.append(big.NewInt(1))
	// c0 <- c(s0)**e mod n
	c0 := big.NewInt(0)
	c0.Exp(s[0], e, n)
	c0.Mul(c0, c)
	c0.Mod(c0, n)
	// M0 <- {[2B, 3B-1]}
	a := big.NewInt(0)
	a.Set(B_2)
	b := big.NewInt(0)
	b.Sub(B_3, big.NewInt(1))
	M.append(make([]Pair, 0))
	M[0].append(Pair{a, b})
	// i <- 1
	i := 1

	for {
		s.append(big.Int(0))
		// Step 2
		if i == 1 {
			// Step 2.a
			// s1 <- n/(3B)
			s[1].Div(n, B_3)
			c1 := findNextSPKCS(s[1], big.NewInt(0), c0, e, n)
		} else if i > 1 && len(M[i-1]) >= 2 {
			// Step 2.b
			// s_i <- s_{i-1} + 1
			s[i].Set(s[i-1])
			s[i].Add(s[i], big.NewInt(1))
			c1 := findNextSPKCS(s[i], big.NewInt(0), c0, e, n)
		} else if len(M[i-1]) == 1 {
			// Step 2.c
			a := M[i-1][0].a
			b := M[i-1][0].b
			// ri = 2*(b*s_{i-1} - 2B)/n
			ri := big.NewInt(0)
			ri.Mul(b, s[i-1])
			ri.Sub(ri, B_2)
			ri.Div(ri, n)
			ri.Mul(big.NewInt(2), ri)

			for {
				// si = (2B + ri*n)/b
				s[i].Mul(ri, n)
				s[i].Add(B_2, si)
				s[i].Div(s[i], b)

				// max_si = (3B + ri*n)/a
				max_si := big.NewInt(0)
				max_si.Mul(ri, n)
				max_si.Add(B_3, max_si)
				max_si.Div(max_si, a)
				c1, ok := findNextSPKCS(s[i], max_si, c, e, n)
				if ok {
					break
				}
				r1.Add(r1, big.NewInt(1))
			}

		} else {
			panic("unreachable")
		}

		// Step 3
		M.append(make([]Pair, 0))

		// Step 4
		if len(M[i]) == 1 && M[i][0].a.Cmp(M[i][0].b) == 0 {
			// m <- a(s_0)**{-1} mod n
			m := big.NewInt(0)
			m.ModInverse(s[0], n)
			m.Mul(M[i][0].a, m)
			m.Mod(m, n)
			return m
		} else {
			// i <- i + 1
			i += 1
		}
	}
}

func findNextSPKCS(s, max_s, c, e, n *big.Int) (*big.Int, bool) {
	c1 := big.NewInt(0)
	for {
		if max_s.BitLen() != 0 && s.Cmp(max_s) != -1 {
			return c1, false
		}
		// c1 <- c0(s1)**e mod n
		c1.Exp(s, e, n)
		c1.Mul(c1, c)
		c1.Mod(c1, n)
		if oracleIsPKCS(c1) {
			return c1, true
		}
		// s1 <- s1+1
		s.Add(s, big.NewInt(1))
	}
}
