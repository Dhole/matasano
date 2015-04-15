package main

import (
	"./hmac"
	crand "crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"time"
)

var passwords = []string{
	"123456",
	"password",
	"welcome",
	"ninja",
	"abc123",
	"123456789",
	"12345678",
	"sunshine",
	"princess",
	"qwerty",
}

var agree Agree

type Agree struct {
	N    big.Int
	g    big.Int
	k    big.Int
	I, P string
}

func initAgree() {
	_N, _ := hex.DecodeString("ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552bb9ed529077096966d670c354e4abc9804f1746c08ca237327ffffffffffffffff")
	_g, _ := hex.DecodeString("02")
	_k, _ := hex.DecodeString("03")

	agree.N.SetBytes(_N)
	agree.g.SetBytes(_g)
	agree.k.SetBytes(_k)

	agree.I = "john@mail.com"
	agree.P = passwords[rand.Intn(len(passwords))]
	fmt.Println("Using password:", agree.P)
}

func carol(con chan []byte, fin chan int, email, pass string) {

	n := int64(math.Ceil(float64(agree.N.BitLen()) / 8.0))

	// Carol computes her secret a la Diffie Hellman
	var a big.Int
	var A big.Int
	_a := make([]byte, n)
	_, err := crand.Read(_a)
	if err != nil {
		panic(err)
	}
	a.SetBytes(_a)
	// A=g**a % N
	A.Exp(&agree.g, &a, &agree.N)

	// Carol sends I and A to Steve
	con <- []byte(email)
	con <- A.Bytes()

	// Receive salt, B and u from Steve
	salt := <-con
	var B big.Int
	B.SetBytes(<-con)
	var u big.Int
	u.SetBytes(<-con)

	// Compute salted hash of password
	h := sha256.New()
	h.Write(salt)
	h.Write([]byte(pass))
	xH := h.Sum(nil)
	fmt.Println("Carol x:", xH)

	var x big.Int
	x.SetBytes(xH)

	// Carol generates S and K
	var S big.Int
	// S = B**(a + ux) % N
	S.Add(&a, big.NewInt(0).Mul(&u, &x))
	S.Mod(&S, &agree.N)
	S.Exp(&B, &S, &agree.N)
	fmt.Println("Carol S:", S.Bytes())

	h = sha256.New()
	h.Write(S.Bytes())
	K := h.Sum(nil)

	// Carol computes HMAC-SHA256(K, salt)
	h = sha256.New()
	sig := hmac.Calc(h, K, salt)

	// Carol sends hmac to Steve
	con <- sig

	// Carol receives response from Steve
	ack := string(<-con)

	if ack == "OK" {
		fmt.Println("Carol: My password was accepted by Steve")
	} else {
		fmt.Println("Carol: My password was REJECTED by Steve")
	}

	fin <- 0
}

func steve(con chan []byte, fin chan int) {
	// Steve:
	// Generate random salt
	salt := make([]byte, 32)
	_, err := crand.Read(salt)
	if err != nil {
		panic(err)
	}

	// We don't know the password yet
	/*
		h := sha256.New()
		h.Write(salt)
		h.Write([]byte(agree.P))
		xH := h.Sum(nil)
		fmt.Println("Steve x:", xH)

		var x big.Int
		x.SetBytes(xH)

		// Compute password verfier v
		var v big.Int
		v.Exp(&agree.g, &x, &agree.N)

		// Forget x and xH
		x.SetInt64(0)
		xH = nil
	*/

	n := int64(math.Ceil(float64(agree.N.BitLen()) / 8.0))

	// Receive I and A from Carol
	I := string(<-con)
	fmt.Println("Bad-Steve: Recived credentials:", I)
	var A big.Int
	A.SetBytes(<-con)

	/*
		if I != agree.I {
			panic("Wrong email")
		}
	*/

	// Steve computes his secret with the verifier
	var b big.Int
	var B big.Int
	_b := make([]byte, n)
	_, err = crand.Read(_b)
	if err != nil {
		panic(err)
	}
	b.SetBytes(_b)
	// B=g**b % N
	B.Exp(&agree.g, &b, &agree.N)

	// u = 128 bit random number
	var u big.Int
	_u := make([]byte, 128/8)
	_, err = crand.Read(_u)
	if err != nil {
		panic(err)
	}
	u.SetBytes(_u)

	// Steve sends salt, B and u to Carol
	con <- salt
	con <- B.Bytes()
	con <- u.Bytes()

	// We can't generate the HMAC because we don't know the password
	/*
		// Steve generates S and K
		var S big.Int
		// S = (A * v**u) ** b % N
		S.Mul(&A, big.NewInt(0).Exp(&v, &u, &agree.N))
		//S.Mod(&S, &agree.N)
		S.Exp(&S, &b, &agree.N)
		fmt.Println("Steve S:", S.Bytes())

		h = sha256.New()
		h.Write(S.Bytes())
		K := h.Sum(nil)

		// Steve computes HMAC-SHA256(K, salt)
		h = sha256.New()
		sig := hmac.Calc(h, K, salt)
	*/

	// Steve receives hmac from Carol
	sig2 := <-con

	// Send fake OK
	con <- []byte("OK")
	fmt.Println("Bad-Steve: Fake-accepting password from Carol")

	/*
		if subtle.ConstantTimeCompare(sig, sig2) == 1 {
			con <- []byte("OK")
			fmt.Println("Steve: Accepting password from Carol")
		} else {
			con <- []byte("KO")
			fmt.Println("Steve: REJECTING password from Carol")
		}
	*/

	fmt.Println("Bad-Steve: Cracking Carol password...")
	// We brute-force the hmac received from Carol against all the
	// passwords in our dictionary offline.
	for _, p := range passwords {
		h := sha256.New()
		h.Write(salt)
		h.Write([]byte(p))
		xH := h.Sum(nil)

		var x big.Int
		x.SetBytes(xH)

		// Compute password verfier v
		var v big.Int
		v.Exp(&agree.g, &x, &agree.N)

		// Steve generates S and K
		var S big.Int
		// S = (A * v**u) ** b % N
		S.Mul(&A, big.NewInt(0).Exp(&v, &u, &agree.N))
		//S.Mod(&S, &agree.N)
		S.Exp(&S, &b, &agree.N)

		h = sha256.New()
		h.Write(S.Bytes())
		K := h.Sum(nil)

		// Steve computes HMAC-SHA256(K, salt)
		h = sha256.New()
		sig := hmac.Calc(h, K, salt)

		if subtle.ConstantTimeCompare(sig, sig2) == 1 {
			fmt.Println("Bad-Steve: Found Carol's password:", p)
			fin <- 0
			return
		}
	}
	fmt.Println("Bad-Steve: Couldn't crack Carol's password :,(")

	fin <- 0
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initAgree()
	con_ab := make(chan []byte)
	fin := make(chan int)
	go carol(con_ab, fin, agree.I, agree.P)
	go steve(con_ab, fin)

	for i := 0; i < 2; i++ {
		<-fin
	}
}
