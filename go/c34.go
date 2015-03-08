package main

import (
	"./modes"
	"./sha1"
	"bytes"
	"crypto/aes"
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
)

func alice(con chan []byte, fin chan int) {
	msg := []byte("Yellow Submarine")
	// Alice decides p and g
	_p, _ := hex.DecodeString("ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec6f44c42e9a637ed6b0bff5cb6f406b7edee386bfb5a899fa5ae9f24117c4b1fe649286651ece45b3dc2007cb8a163bf0598da48361c55d39a69163fa8fd24cf5f83655d23dca3ad961c62f356208552bb9ed529077096966d670c354e4abc9804f1746c08ca237327ffffffffffffffff")
	_g, _ := hex.DecodeString("02")

	var p big.Int
	var g big.Int
	p.SetBytes(_p)
	g.SetBytes(_g)

	N := int64(math.Ceil(float64(p.BitLen()) / 8.0))

	// Alice computes her secret
	var a big.Int
	var A big.Int
	_a := make([]byte, N)
	_, err := crand.Read(_a)
	if err != nil {
		panic(err)
	}
	a.SetBytes(_a)
	A.Exp(&g, &a, &p)

	// Alice sends the p, g and A to Bob
	con <- p.Bytes()
	con <- g.Bytes()
	con <- A.Bytes()

	// Alice receives B from Bob
	var B big.Int
	B.SetBytes(<-con)

	// Alice computes the shared secret
	var s big.Int
	s.Exp(&B, &a, &p)

	fmt.Println("Alice> shared secret:", s.String())

	// Key derivation from shared secret
	h := sha1.New()
	h.Write(s.Bytes())
	key := h.Sum(nil)[0:16]

	ciptxt := make([]byte, 16)
	iv := make([]byte, 16)
	_, err = crand.Read(iv)
	if err != nil {
		panic(err)
	}

	// Encrypt message with derivated key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cbc_enc := modes.NewCBCEncrypter(block, iv)
	cbc_enc.CryptBlocks(ciptxt, msg)

	// Send ciphertext to Bob
	fmt.Println("Alice> seding:", ciptxt, "->", string(msg))
	con <- iv
	con <- ciptxt

	// Receive reencripted message from Bob
	iv2 := <-con
	ciptxt2 := <-con

	// Decrypt message from Bob
	msg2 := make([]byte, len(ciptxt2))
	cbc_dec := modes.NewCBCDecrypter(block, iv2)
	cbc_dec.CryptBlocks(msg2, ciptxt2)
	fmt.Println("Alice> received:", ciptxt2, "->", string(msg2))

	if bytes.Compare(msg, msg2) == 0 {
		fmt.Println("Decripted sent message is the same as decripted received message")
	} else {
		fmt.Println("Reiceved decripted message is different")
	}

	fin <- 0
}

func bob(con chan []byte, fin chan int) {
	var p big.Int
	var g big.Int
	var A big.Int

	// Bob receives p, g and A from Alice
	p.SetBytes(<-con)
	g.SetBytes(<-con)
	A.SetBytes(<-con)

	N := int64(math.Ceil(float64(p.BitLen()) / 8.0))

	// Bob computes his secret
	var b big.Int
	var B big.Int
	_b := make([]byte, N)
	_, err := crand.Read(_b)
	if err != nil {
		panic(err)
	}
	b.SetBytes(_b)
	B.Exp(&g, &b, &p)

	// Bob sends B to Alice
	con <- B.Bytes()

	// Bob computes the shared secret
	var s big.Int
	s.Exp(&A, &b, &p)

	fmt.Println("Bob> shared secret:", s.String())

	// Key derivation from shared secret
	h := sha1.New()
	h.Write(s.Bytes())
	key := h.Sum(nil)[0:16]

	// Receive ciphertext from Alice
	iv := <-con
	ciptxt := <-con

	// Decrypt message from Alice
	msg := make([]byte, len(ciptxt))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	cbc_dec := modes.NewCBCDecrypter(block, iv)
	cbc_dec.CryptBlocks(msg, ciptxt)
	fmt.Println("Bob> received:", ciptxt, "->", string(msg))

	// Encrypt the message again
	ciptxt2 := make([]byte, 16)
	iv2 := make([]byte, 16)
	_, err = crand.Read(iv2)
	if err != nil {
		panic(err)
	}

	cbc_enc := modes.NewCBCEncrypter(block, iv2)
	cbc_enc.CryptBlocks(ciptxt2, msg)

	// Send the re encripted message back to Alice
	fmt.Println("Bob> sending:", ciptxt2, "->", string(msg))
	con <- iv2
	con <- ciptxt2

	fin <- 0
}

func mitm(con_a, con_b chan []byte) {
	// Forward messages
	/*for {
		select {
		case msg_a := <-con_a:
			con_b <- msg_a
		case msg_b := <-con_b:
			con_a <- msg_b
		}
	}*/

	// We will A and B for p. Since p**n mod p == 0 mod p, we know the
	// shared secret: It's always 0
	var s big.Int
	s.SetInt64(0)
	h := sha1.New()
	h.Write(s.Bytes())
	key := h.Sum(nil)[0:16]

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Receive p, g, A from Alice
	p := <-con_a
	g := <-con_a
	//A := <-con_a
	_ = <-con_a

	// Send p, g, p to Bob
	con_b <- p
	con_b <- g
	con_b <- p

	// Receive B from Bob
	//B := <-con_b
	_ = <-con_b

	// Send p to Alice
	con_a <- p

	// Relay the ciphered messages exchange
	iv_a := <-con_a
	ciptxt_a := <-con_a

	// Decipher Alice to Bob message:
	clrtxt_a := make([]byte, len(ciptxt_a))
	cbc_dec_a := modes.NewCBCDecrypter(block, iv_a)
	cbc_dec_a.CryptBlocks(clrtxt_a, ciptxt_a)
	fmt.Println("Mallory> message from Alice to Bob:", ciptxt_a, "->",
		string(clrtxt_a))

	con_b <- iv_a
	con_b <- ciptxt_a

	iv_b := <-con_b
	ciptxt_b := <-con_b

	// Decipher Bob to Alice message:
	clrtxt_b := make([]byte, len(ciptxt_b))
	cbc_dec_b := modes.NewCBCDecrypter(block, iv_b)
	cbc_dec_b.CryptBlocks(clrtxt_b, ciptxt_b)
	fmt.Println("Mallory> message from Bob to Alice:", ciptxt_b, "->",
		string(clrtxt_b))

	con_a <- iv_b
	con_a <- ciptxt_b
}

func main() {
	con_a := make(chan []byte)
	con_b := make(chan []byte)
	fin := make(chan int)
	go alice(con_a, fin)
	go bob(con_b, fin)
	go mitm(con_a, con_b)

	for i := 0; i < 2; i++ {
		<-fin
	}
}
