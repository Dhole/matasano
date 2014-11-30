package pkcs7

import (
	"errors"
	"fmt"
)

func Pad(buf []byte, blockSize int) (buf_padded []byte, err error) {
	if blockSize > 256 {
		return nil, fmt.Errorf("Block size is larger than 256")
	}
	n := blockSize - (len(buf) % blockSize)
	if n == 0 {
		n = 0xff
	}
	padding := make([]byte, n)
	for i := 0; i < n; i++ {
		padding[i] = byte(n)
	}
	buf_padded = append(buf, padding...)
	return buf_padded, nil
}

func Unpad(buf []byte) (buf_unpadded []byte, err error) {
	err = errors.New("Input is not a PKCS#7 padded block")
	if len(buf) == 0 {
		return nil, err
	}
	pad_byte := buf[len(buf)-1]
	n := int(pad_byte)
	if n > len(buf) {
		return nil, err
	}
	for i := len(buf) - n; i < len(buf); i++ {
		if buf[i] != pad_byte {
			return nil, err
		}
	}
	buf_unpadded = buf[:len(buf)-n]
	return buf_unpadded, nil
}
