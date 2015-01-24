package main

import (
	"./mt19937"
	"fmt"
	"time"
)

func untemper(v uint32) (res uint32) {
	res = v
	res = (res & 0xFFFFC000) | (((res >> 18) ^ v) & 0xFFFFFFFF)

	v = res
	res = (res & 0x00007FFF) | ((((res << 15) & 0xefc60000) ^ v) & 0x7FFFFFFF)
	res = (res & 0x00007FFF) | ((((res << 15) & 0xefc60000) ^ v) & 0xFFFFFFFF)

	v = res
	res = (res & 0x0000007F) | ((((res << 7) & 0x9d2c5680) ^ v) & 0x00003FFF)
	res = (res & 0x0000007F) | ((((res << 7) & 0x9d2c5680) ^ v) & 0x001FFFFF)
	res = (res & 0x0000007F) | ((((res << 7) & 0x9d2c5680) ^ v) & 0x0FFFFFFF)
	res = (res & 0x0000007F) | ((((res << 7) & 0x9d2c5680) ^ v) & 0xFFFFFFFF)

	v = res
	res = (res & 0xFFE00000) | (((res >> 11) ^ v) & 0xFFFFFC00)
	res = (res & 0xFFE00000) | (((res >> 11) ^ v) & 0xFFFFFFFE)
	res = (res & 0xFFE00000) | (((res >> 11) ^ v) & 0xFFFFFFFF)

	return res
}

func main() {

	num := 624

	m := mt19937.New()
	m.Init(uint32(time.Now().Unix()))
	//m.Init(0)

	rnd := make([]uint32, 624)
	my_MT := make([]uint32, 624)

	fmt.Println("Storing 624 random numbers...")
	for i := 0; i < num; i++ {
		rnd[i] = m.ExtractNumber()
	}
	fmt.Println("Calculating state...")
	for i := 0; i < num; i++ {
		my_MT[i] = untemper(rnd[i])
	}

	my_m := mt19937.New()
	my_m.MT = my_MT
	my_m.Index = 0

	fmt.Println("Generating new values from original and replica...")
	for i := 0; i < 10; i++ {
		fmt.Println("Original:", m.ExtractNumber())
		fmt.Println("Replica: ", my_m.ExtractNumber())
	}
}
