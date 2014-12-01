package main

import (
	"./set1"
	"./set2"
	"bytes"
	"fmt"
	"reflect"
)

type oracle func([]byte) []byte

func FindOracleBlockSize(ora oracle) (bs int) {
	bs = 0
	msg := bytes.Repeat([]byte{byte('A')}, 0)
	out := ora(msg)
	len_fst := len(out)
	for i := 1; i < 256; i++ {
		msg = bytes.Repeat([]byte{byte('A')}, i)
		out := ora(msg)
		if len(out) > len_fst {
			bs = len(out) - len_fst
			break
		}
	}
	return bs
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
	//fmt.Println(fill)
	for i := 1; i < bs+1; i++ {
		//fmt.Println(i, ":")
		win1 := fill[i-1 : bs-1+i]
		win2 := fill[i-1 : len_prev+bs-1+i]
		//fmt.Println(win)
		base := ora(win1[:bs-i])
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
	return blk
}

func FindOracleClear(ora oracle, bs int) (clear []byte) {
	len_clear := len(ora(make([]byte, 0)))
	clear = make([]byte, len_clear)
	n_blks := len(clear) / bs
	blk := make([]byte, bs)

	blk = FindOracleBlock(ora, make([]byte, 0), bs)
	copy(clear[:bs], blk)
	//n_blks = 2
	for i := 1; i < n_blks; i++ {
		blk = FindOracleBlock(ora, clear[:bs*i], bs)
		copy(clear[bs*i:bs*(i+1)], blk)
	}
	return clear
}

func main() {
	set2.SetOracle2Key()
	//input := make([]byte, 16*3)
	//output := set2.AESEncryptionOracle2(input)
	//fmt.Println(output)
	bs := FindOracleBlockSize(set2.AESEncryptionOracle2)
	if IsOracleECB(set2.AESEncryptionOracle2, bs) {
		fmt.Println("It's ECB")
	} else {
		fmt.Println("It's not ECB")
		return
	}
	blk := FindOracleClear(set2.AESEncryptionOracle2, bs)
	fmt.Println(string(blk))
	fmt.Println(blk)

}
