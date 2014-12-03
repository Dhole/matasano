package main

import (
	"./set1"
	"./set2"
	"bytes"
	"fmt"
	"reflect"
)

type oracle func([]byte) []byte

func FindOracleBlockSize(ora oracle) (bs, pad_size int) {
	bs = 0
	msg := bytes.Repeat([]byte{byte('A')}, 0)
	out := ora(msg)
	len_fst := len(out)
	for i := 1; i < 256; i++ {
		msg = bytes.Repeat([]byte{byte('A')}, i)
		out := ora(msg)
		if len(out) > len_fst {
			bs = len(out) - len_fst
			pad_size = i
			break
		}
	}
	return bs, pad_size
}

func IsOracleECB(ora oracle, bs int) (is_ECB bool) {
	input := make([]byte, 16*3)
	output := ora(input)
	return set1.IsECB(output, bs)
}

func FindOracleBlock(ora oracle, prev []byte, bs int) (blk []byte) {
	len_prev := len(prev)
	blk = make([]byte, bs)
	fill := bytes.Repeat([]byte{byte('A')}, bs*2+len_prev)
	copy(fill[bs-1:], prev)
	for i := 1; i < bs+1; i++ {
		win1 := fill[i-1 : bs-1]
		win2 := fill[i-1 : len_prev+bs-1+i]
		base := ora(win1)
		for j := 0; j < 256; j++ {
			win2[len_prev+bs-1] = byte(j)
			out := ora(win2)
			if reflect.DeepEqual(out[len_prev:len_prev+bs],
				base[len_prev:len_prev+bs]) {
				blk[i-1] = byte(j)
				break
			}
		}
	}
	//fmt.Println("block:\n", blk)
	return blk
}

func FindOraclePlain(ora oracle, bs, pad_size int) (plain []byte) {
	len_plain := len(ora(make([]byte, 0)))
	plain = make([]byte, len_plain)
	n_blks := len(plain) / bs
	blk := make([]byte, bs)

	blk = FindOracleBlock(ora, make([]byte, 0), bs)
	copy(plain[:bs], blk)
	//n_blks = 2
	for i := 1; i < n_blks; i++ {
		blk = FindOracleBlock(ora, plain[:bs*i], bs)
		copy(plain[bs*i:], blk)
	}
	return plain[:len_plain-pad_size]
}

func main() {
	set2.SetOracle2Key()
	/*
		input := make([]byte, 0)
		output := set2.AESEncryptionOracle2(input)
		fmt.Println("output:\n", output)
	*/
	bs, pad_size := FindOracleBlockSize(set2.AESEncryptionOracle2)
	fmt.Println("bs:", bs, "pad size:", pad_size)
	if IsOracleECB(set2.AESEncryptionOracle2, bs) {
		fmt.Println("It's ECB")
	} else {
		fmt.Println("It's not ECB")
		return
	}
	blk := FindOraclePlain(set2.AESEncryptionOracle2, bs, pad_size)
	fmt.Println(string(blk))
	fmt.Println(blk)

}
