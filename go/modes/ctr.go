package modes

import (
	"../set1"
	"crypto/cipher"
	"encoding/binary"
)

type ctr struct {
	b            cipher.Block
	blockSize    int
	keyStreamBlk []byte
	uniCounter   []byte
	nonce        uint64
	counter      uint64
}

func NewCTR(b cipher.Block, nonce uint64) *ctr {
	return &ctr{
		b:            b,
		blockSize:    b.BlockSize(),
		keyStreamBlk: make([]byte, b.BlockSize()),
		uniCounter:   make([]byte, b.BlockSize()),
		nonce:        nonce,
		counter:      0,
	}
}

func NewCTR2(b cipher.Block, nonce uint64, counter uint64) *ctr {
	return &ctr{
		b:            b,
		blockSize:    b.BlockSize(),
		keyStreamBlk: make([]byte, b.BlockSize()),
		uniCounter:   make([]byte, b.BlockSize()),
		nonce:        nonce,
		counter:      counter,
	}
}

func (x *ctr) XORKeyStream(dst, src []byte) {
	if len(dst) < len(src) {
		panic("modes: output smaller than input")
	}
	binary.LittleEndian.PutUint64(x.uniCounter[:8], x.nonce)

	for {
		binary.LittleEndian.PutUint64(x.uniCounter[8:], x.counter)
		x.b.Encrypt(x.keyStreamBlk, x.uniCounter)
		set1.XorBytes(dst, src, x.keyStreamBlk)
		if len(src) < x.blockSize {
			break
		} else {
			src = src[x.blockSize:]
			dst = dst[x.blockSize:]
			x.counter++
		}
	}
	x.counter = 0
}
