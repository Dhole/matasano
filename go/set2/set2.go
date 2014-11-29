package set2

import (
	"../set1"
	"crypto/aes"
	"fmt"
)

func PKCS7Padding(blk []byte, blk_size int) (blk_pad []byte, err error) {
	if len(blk) > blk_size {
		return nil, fmt.Errorf("Block is larger than block size")
	}
	if blk_size > 256 {
		return nil, fmt.Errorf("Block size is larger than 256")
	}
	pad := byte(blk_size - len(blk))
	blk_pad = make([]byte, blk_size)
	copy(blk_pad, blk)
	for i := len(blk); i < blk_size; i++ {
		blk_pad[i] = pad
	}
	return blk_pad, nil
}

func AESCBCDecrypt(ciphertext, iv, key []byte) (plaintext []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		panic("ciphertext is not multiple of the block size")
	}

	if len(ciphertext) < bs {
		panic("ciphertext too short")
	}

	plaintext = make([]byte, len(ciphertext))
	ptxt := plaintext
	ctxt := ciphertext
	prev_cblk := iv
	for len(ctxt) > 0 {
		block.Decrypt(ptxt[:bs], ctxt[:bs])
		copy(ptxt[:bs], set1.XorKey(ptxt[:bs], prev_cblk))
		prev_cblk = ctxt[:bs]
		ptxt = ptxt[bs:]
		ctxt = ctxt[bs:]
	}
	return plaintext
}

func AESCBCEncrypt(plaintext, iv, key []byte) (ciphertext []byte) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	bs := block.BlockSize()
	if len(plaintext)%bs != 0 {
		panic("plaintext is not multiple of the block size")
	}

	if len(plaintext) < bs {
		panic("plaintext too short")
	}

	ciphertext = make([]byte, len(plaintext))
	ptxt := plaintext
	ctxt := ciphertext
	last_cblk := iv
	for len(ptxt) > 0 {
		copy(ptxt[:bs], set1.XorKey(ptxt[:bs], last_cblk))
		block.Encrypt(ctxt[:bs], ptxt[:bs])
		last_cblk = ctxt[:bs]
		ptxt = ptxt[bs:]
		ctxt = ctxt[bs:]
	}
	return plaintext
}
