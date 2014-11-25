package hamming

var oneBits = [...]int{0, 1, 1, 2, 1, 2, 2, 3, 1, 2, 2, 3, 2, 3, 3, 4}

func CountBitsByte(x byte) (c int) {
	return oneBits[x&0x0f] + oneBits[x>>4]
}

func DistByte(x, y byte) (c int) {
	return CountBitsByte(x ^ y)
}

func DistByteSlice(a, b []byte) (c int) {
	for i := 0; i < len(a); i++ {
		c += DistByte(a[i], b[i])
	}
	return c
}
