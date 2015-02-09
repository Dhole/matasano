package hmac

import (
	"../set1"
	"bytes"
	//"fmt"
	"hash"
)

func Calc(h hash.Hash, msg []byte, key []byte) (hmac []byte) {
	bs := h.BlockSize()
	if len(key) > bs {
		h.Reset()
		h.Write(key)
		key = h.Sum(nil)
	} else if len(key) < bs {
		key = append(key, bytes.Repeat([]byte{byte(0x00)}, bs-len(key))...)
	}
	//fmt.Println("Key len:", len(key))
	//fmt.Println("key:", key)
	o_key_pad := make([]byte, bs)
	i_key_pad := make([]byte, bs)
	set1.XorBytes(o_key_pad, bytes.Repeat([]byte{byte(0x5C)}, bs), key)
	set1.XorBytes(i_key_pad, bytes.Repeat([]byte{byte(0x36)}, bs), key)

	h.Reset()
	h.Write(i_key_pad)
	h.Write(msg)
	hmac = h.Sum(nil)

	h.Reset()
	h.Write(o_key_pad)
	h.Write(hmac)
	hmac = h.Sum(nil)

	return hmac
}
