package main

import (
	"./set2"
	"bytes"
	"fmt"
	"math"
	"reflect"
)

type oracle func([]byte) []byte

func FindRepBlock(input []byte, bs int) (i, j int) {
	n_blocks := len(input) / bs
	for i = 0; i < n_blocks-1; i++ {
		blk_a := input[bs*i : bs*(i+1)]
		for j = i + 1; j < n_blocks; j++ {
			blk_b := input[bs*j : bs*(j+1)]
			if bytes.Equal(blk_a, blk_b) {
				return i, j
			}
		}
	}
	return 0, 0
}

func FindOracleBlockSize(ora oracle, len_pref int) (bs, pad_size int) {
	bs = 0
	msg := make([]byte, len_pref+256)
	copy(msg[len_pref:], bytes.Repeat([]byte{byte('A')}, 256))
	out := ora(msg[:len_pref])
	len_fst := len(out)
	for i := 1; i < 256; i++ {
		out := ora(msg[:len_pref+i])
		if len(out) > len_fst {
			bs = len(out) - len_fst
			pad_size = i
			break
		}
	}
	return bs, pad_size
}

func FindOracleBlock(ora oracle, pref, prev []byte, len_skip, bs int) (blk []byte) {
	len_prev := len(prev)
	blk = make([]byte, bs)
	fill := bytes.Repeat([]byte{byte('A')}, bs*2+len_prev)
	copy(fill[bs-1:], prev)
	for i := 1; i < bs+1; i++ {
		win1 := fill[i-1 : bs-1]
		win2 := fill[i-1 : len_prev+bs-1+i]
		base := ora(append(pref, win1...))
		for j := 0; j < 256; j++ {
			win2[len_prev+bs-1] = byte(j)
			out := ora(append(pref, win2...))
			if reflect.DeepEqual(out[len_skip+len_prev:len_skip+len_prev+bs],
				base[len_skip+len_prev:len_skip+len_prev+bs]) {
				blk[i-1] = byte(j)
				break
			}
		}
	}
	//fmt.Println("block:\n", blk)
	return blk
}

func FindOraclePlain(ora oracle, bs int) (plain []byte) {
	len_pref := FindPrefSize(ora, bs)
	fmt.Println("> Prefix size:", len_pref)
	len_pad_pref := int(math.Ceil(float64(len_pref)/float64(bs)))*bs - len_pref
	pad_pref := bytes.Repeat([]byte{byte('B')}, len_pad_pref)
	fmt.Println("> Pad Prefix size:", len_pad_pref)
	len_skip := len_pref + len_pad_pref
	fmt.Println("> Skip size:", len_skip)

	_, pad_size := FindOracleBlockSize(ora, len_pad_pref)
	fmt.Println("> Pad size:", pad_size)

	len_plain := len(ora(pad_pref)) - len_pref
	plain = make([]byte, len_plain)
	n_blks := (len_plain - len_pad_pref) / bs
	blk := make([]byte, bs)

	copy(plain, pad_pref)
	blk = FindOracleBlock(ora, pad_pref, make([]byte, 0), len_skip, bs)
	fmt.Println("> 0th blk:", blk)
	//panic("stop")
	copy(plain[:bs], blk)
	//n_blks = 2
	for i := 1; i < n_blks; i++ {
		blk = FindOracleBlock(ora, pad_pref, plain[:bs*i], len_skip, bs)
		fmt.Println(">", i, "th blk:", blk)
		copy(plain[bs*i:], blk)
	}
	return plain[:len_plain-pad_size]
}

func FindPrefSize(ora oracle, bs int) (size int) {
	pad := bytes.Repeat([]byte{byte('A')}, bs*3)

	pad_size := 0
	ciphtxt := ora(pad[:bs*3])
	ifst, _ := FindRepBlock(ciphtxt, bs)
	for i := bs*3 - 1; i >= 0; i-- {
		ciphtxt := ora(pad[:i])
		ib, jb := FindRepBlock(ciphtxt, bs)
		if ib == 0 && jb == 0 {
			pad_size = i + 1 - bs*2
			break
		}
	}
	return ifst*bs - pad_size
}

func main() {
	bs := 16
	set2.SetRndSeed()
	set2.SetOracle2Key()
	set2.SetOracle3RndPrefix()
	//bs, pad_size := FindOracleBlockSize(set2.AESEncryptionOracle3, 10)
	//fmt.Println("Block size:", bs, ", Pad size:", pad_size)
	//fmt.Println("Pref size:", FindPrefSize(set2.AESEncryptionOracle3, bs))
	fmt.Print("Plaintext:\n", string(FindOraclePlain(set2.AESEncryptionOracle3, bs)))
	//fmt.Println("My Pad:", my_pad_size)
	//ciphtxt := set2.AESEncryptionOracle3([]byte(msg))
	//fmt.Println(ciphtxt)
}
