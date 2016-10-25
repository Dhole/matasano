package dsa

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"hash"
	"math/big"
)

func print_hex(pre string, v *big.Int) {
	fmt.Printf("%v = 0x%v\n", pre, hex.EncodeToString(v.Bytes()))
}

func print_dec(pre string, v *big.Int) {
	fmt.Printf("%v = %v\n", pre, v.String())
}

type Params struct {
	Hash    hash.Hash
	P, Q, G *big.Int
}

type Keys struct {
	X, Y *big.Int
}

type Signature struct {
	R, S *big.Int
}

func main() {
	fmt.Println("vim-go")
}

func ParameterGen(hash hash.Hash, L, N int) (*Params, error) {
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

	return &Params{hash, p, q, g}, nil
}

func KeyGen(par *Params) Keys {
	x, err := crand.Int(crand.Reader, par.Q)
	if err != nil {
		panic(err)
	}
	y := big.NewInt(0)
	y.Exp(par.G, x, par.P)

	return Keys{x, y}
}

func Sign(par *Params, keys *Keys, msg []byte) Signature {
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
		if r.Cmp(big.NewInt(0)) == 0 {
			continue
		}
		//print_dec("r", r)

		k1 := big.NewInt(0).ModInverse(k, par.Q)
		s.Mul(keys.X, r)
		s.Add(h, s)
		s.Mul(k1, s)
		s.Mod(s, par.Q)
		if s.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		//r1 := big.NewInt(0).ModInverse(r, par.Q)
		//x1 := x_from_k(r1, s, par.Q, h, k)
		//print_dec("x1", x1)

		break
	}
	return Signature{r, s}
}

func Verify(par *Params, sig *Signature, y *big.Int, msg []byte) bool {
	if sig.R.Cmp(big.NewInt(0)) != 1 && sig.R.Cmp(par.Q) != -1 {
		return false
	}
	if sig.S.Cmp(big.NewInt(0)) != 1 && sig.S.Cmp(par.Q) != -1 {
		return false
	}

	w := big.NewInt(0)
	u1 := big.NewInt(0)
	u2 := big.NewInt(0)
	v := big.NewInt(0)

	par.Hash.Reset()
	par.Hash.Write(msg)
	h_b := par.Hash.Sum(nil)
	fmt.Printf("h = 0x%v\n", hex.EncodeToString(h_b))
	h := big.NewInt(0).SetBytes(h_b)

	w.ModInverse(sig.S, par.Q)
	print_dec("w", w)

	u1.Mul(h, w)
	u1.Mod(u1, par.Q)
	print_dec("u1", u1)

	u2.Mul(sig.R, w)
	u2.Mod(u2, par.Q)
	print_dec("u2", u2)

	u1.Exp(par.G, u1, par.P)
	u2.Exp(y, u2, par.P)

	v.Mul(u1, u2)
	v.Mod(v, par.P)
	v.Mod(v, par.Q)
	print_dec("v", v)
	print_dec("r", sig.R)

	if v.Cmp(sig.R) == 0 {
		return true
	}
	return false
}

func XFromK(r1, s, q, h, k *big.Int) *big.Int {
	//r1 := big.NewInt(0).ModInverse(r, q)
	x1 := big.NewInt(0)
	x1.Mul(s, k)
	x1.Sub(x1, h)
	x1.Mul(x1, r1)
	x1.Mod(x1, q)

	return x1
}
