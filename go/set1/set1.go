package set1

func Xor(a, b []byte) (x []byte) {
	x = make([]byte, len(a))
	for i := range x {
		x[i] = a[i] ^ b[i]
	}
	return x
}
