package random_test

import (
	"encoding/hex"
	"math/rand/v2"
	"testing"
	random "ws-load/pkg/random/pure"
)

func TestBytes_ShouldBePure(t *testing.T) {
	t.Parallel()

	genBytes := func() string {
		gen := rand.New(rand.NewPCG(2, 2)) //nolint:gosec

		return hex.EncodeToString(random.Bytes(gen, 16))
	}

	const n = 20

	base := genBytes()

	for range n {
		if base != genBytes() {
			t.Errorf("bytes should be pure, always return the same result")
		}
	}
}

func BenchmarkBytes(b *testing.B) {
	gen := rand.New(rand.NewPCG(8, 8)) //nolint:gosec

	for b.Loop() {
		random.Bytes(gen, 64)
	}
}
