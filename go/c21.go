package main

import (
	"./mt19937"
	"fmt"
)

func main() {
	m := mt19937.New()
	m.Init(0x1571)

	for i := 0; i < 10; i++ {
		fmt.Printf("0x%.8x\n", m.ExtractNumber())
	}
}
