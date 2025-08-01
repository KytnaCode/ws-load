package random

import "math/rand/v2"

func Bytes(gen *rand.Rand, size int) []byte {
	b := make([]byte, 0, size)

	var v uint64

	var mask uint64 = 0b11111111

	for i := range size {
		if i%8 == 0 {
			v = gen.Uint64()
		}

		b = append(b, byte(v&mask))
		v >>= 8
	}

	return b
}
