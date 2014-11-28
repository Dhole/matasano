package set2

import (
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
