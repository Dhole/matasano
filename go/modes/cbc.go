package modes

import (
	"../set1"
	"crypto/cipher"
)

func dup(p []byte) []byte {
	q := make([]byte, len(p))
	copy(q, p)
	return q
}

// Heavily copied from src/pkg/crypto/cipher/cbc.go

type cbc struct {
	b         cipher.Block
	blockSize int
	iv        []byte
	tmp       []byte
}

func newCBC(b cipher.Block, iv []byte) *cbc {
	return &cbc{
		b:         b,
		blockSize: b.BlockSize(),
		iv:        dup(iv),
		tmp:       make([]byte, b.BlockSize()),
	}
}

type cbcEncrypter cbc

func NewCBCEncrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	if len(iv) != b.BlockSize() {
		panic("modes: IV lenght must equal block size")
	}
	return (*cbcEncrypter)(newCBC(b, iv))
}

func (x *cbcEncrypter) BlockSize() int { return x.blockSize }

func (x *cbcEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("modes: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("modes: output smaller than input")
	}
	iv := x.iv
	for len(src) > 0 {
		set1.XorBytes(dst[:x.blockSize], src[:x.blockSize], iv)
		x.b.Encrypt(dst[:x.blockSize], dst[:x.blockSize])
		iv = dst[:x.blockSize]
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
	copy(x.iv, iv)
}

type cbcDecrypter cbc

func NewCBCDecrypter(b cipher.Block, iv []byte) cipher.BlockMode {
	if len(iv) != b.BlockSize() {
		panic("modes: IV lenght must equal block size")
	}
	return (*cbcDecrypter)(newCBC(b, iv))
}

func (x *cbcDecrypter) BlockSize() int { return x.blockSize }

func (x *cbcDecrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("modes: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("modes: output smaller than input")
	}
	if len(src) == 0 {
		return
	}
	end := len(src)
	start := end - x.blockSize
	prev := start - x.blockSize

	copy(x.tmp, src[start:end])
	for start > 0 {
		x.b.Decrypt(dst[start:end], src[start:end])
		set1.XorBytes(dst[start:end], dst[start:end], src[prev:start])

		end = start
		start = prev
		prev -= x.blockSize
	}
	x.b.Decrypt(dst[start:end], src[start:end])
	set1.XorBytes(dst[start:end], dst[start:end], x.iv)
	x.iv, x.tmp = x.tmp, x.iv
}
