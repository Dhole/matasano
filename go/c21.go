package main

import (
	"fmt"
)

type MT19937 struct {
	index int
	MT    []uint32
}

func (m *MT19937) Init(seed uint32) {
	m.MT = make([]uint32, 624)
	m.index = 0
	m.MT[0] = seed
	for i := uint32(1); i < 624; i++ {
		m.MT[i] = (0x6c078965*(m.MT[i-1]^(m.MT[i-1]>>30)) + i)
	}
}

func (m *MT19937) ExtractNumber() uint32 {
	if m.index == 0 {
		m.generateNumbers()
	}

	y := m.MT[m.index]
	y = y ^ (y >> 11)
	y = y ^ ((y << 7) & 0x9d2c5680)
	y = y ^ ((y << 15) & 0xefc60000)
	y = y ^ (y >> 18)

	m.index = (m.index + 1) % 264
	return y
}

func (m *MT19937) generateNumbers() {
	for i := 0; i < 624; i++ {
		y := (m.MT[i] & 0x80000000) + (m.MT[(i+1)%624] & 0x7fffffff)
		m.MT[i] = m.MT[(i+397)%624] ^ (y >> 1)
		if (y % 2) != 0 {
			m.MT[i] = m.MT[i] ^ 0x9908b0df
		}
	}
}

func main() {
	var m MT19937
	m.Init(0x1571)

	for i := 0; i < 10; i++ {
		fmt.Printf("0x%.8x\n", m.ExtractNumber())
	}
}
