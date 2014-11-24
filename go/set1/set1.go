package set1

func Xor(a, b []byte) (x []byte) {
	x = make([]byte, len(a))
	for i := range x {
		x[i] = a[i] ^ b[i]
	}
	return x
}

func Xor_key(a, b []byte) (x []byte) {
	x = make([]byte, len(a))
	len_b := len(b)
	for i := range x {
		x[i] = a[i] ^ b[i%len_b]
	}
	return x
}
