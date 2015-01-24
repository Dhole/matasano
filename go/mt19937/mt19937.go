package mt19937

type MT19937 struct {
	Index int
	MT    []uint32
}

func New() *MT19937 {
	res := &MT19937{
		Index: 0,
		MT:    make([]uint32, 624),
	}
	return res
}

func (m *MT19937) Init(seed uint32) {
	m.Index = 0
	m.MT[0] = seed
	for i := uint32(1); i < 624; i++ {
		m.MT[i] = (0x6c078965*(m.MT[i-1]^(m.MT[i-1]>>30)) + i)
	}
}

func (m *MT19937) ExtractNumber() uint32 {
	if m.Index == 0 {
		m.generateNumbers()
	}

	y := m.MT[m.Index]
	y = y ^ (y >> 11)
	y = y ^ ((y << 7) & 0x9d2c5680)
	y = y ^ ((y << 15) & 0xefc60000)
	y = y ^ (y >> 18)

	m.Index = (m.Index + 1) % 624
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
