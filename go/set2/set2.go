package set2

import (
	//"../set1"
	"crypto/aes"
	"crypto/cipher"
	crand "crypto/rand"
	"math/rand"
	//"time"
	"../modes"
	"../pkcs7"
)

func RandSlice(dst []byte) {
	for i := range dst {
		dst[i] = byte(rand.Intn(256))
	}
}

func AESEncryptionOracle(src []byte) (dst []byte) {
	key := make([]byte, 16)
	_, err := crand.Read(key)
	if err != nil {
		panic(err)
	}
	iv := make([]byte, 16)
	_, err = crand.Read(iv)
	if err != nil {
		panic(err)
	}

	//rand.Seed(time.Now().UnixNano())

	before := make([]byte, 5+rand.Intn(6))
	after := make([]byte, 5+rand.Intn(6))
	RandSlice(before)
	RandSlice(after)

	dst = make([]byte, len(before)+len(src)+len(after))
	copy(dst, before)
	copy(dst[len(before):], src)
	copy(dst[len(before)+len(src):], after)
	dst, err = pkcs7.Pad(dst, 16)
	if err != nil {
		panic(err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	encrypter := func() cipher.BlockMode {
		if rand.Intn(2) == 0 {
			return modes.NewECBEncrypter(block)
		} else {
			return modes.NewCBCEncrypter(block, iv)
		}
	}()
	/*
		if rand.Intn(2) == 0 {
			encrypter = modes.NewECBEncrypter(block)
		} else {
			encrypter = modes.NewCBCEncrypter(block, iv)
		}
	*/
	encrypter.CryptBlocks(dst, dst)
	return dst
}

/*
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
*/
/*
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
*/
