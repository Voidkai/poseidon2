package poseidon2

import (
	"math/big"
	"testing"

	"github.com/iden3/go-iden3-crypto/ff"
	"github.com/iden3/go-iden3-crypto/utils"
	"github.com/stretchr/testify/assert"
)

func TestPoseidonHashFixed(t *testing.T) {
	b := make([]*ff.Element, 16)
	b[0] = from_hex("0x00")
	b[1] = from_hex("0x01")
	b[2] = from_hex("0x02")

	poseidon2 := NewPoseidon2(RoundConstants)

	assert.Equal(t, from_hex("0x00"), b[0])
	assert.Equal(t, from_hex("0x01"), b[1])
	assert.Equal(t, from_hex("0x02"), b[2])

	h := poseidon2.Permutation([]*ff.Element{b[0], b[1], b[2]})
	assert.Equal(t,
		from_hex("0x0bb61d24daca55eebcb1929a82650f328134334da98ea4f847f760054f4a3033"),
		h[0])
	assert.Equal(t,
		from_hex("0x303b6f7c86d043bfcbcc80214f26a30277a15d3f74ca654992defe7ff8d03570"),
		h[1])
	assert.Equal(t,
		from_hex("0x1ed25194542b12eef8617361c3ba7c52e660b145994427cc86296242cf766ec8"),
		h[2])
}

func BenchmarkPoseidon2Hash(b *testing.B) {
	b0 := ff.NewElement().SetBigInt(big.NewInt(0))
	b1 := ff.NewElement().SetBigInt(utils.NewIntFromString("12242166908188651009877250812424843524687801523336557272219921456462821518061")) //nolint:lll
	b2 := ff.NewElement().SetBigInt(utils.NewIntFromString("12242166908188651009877250812424843524687801523336557272219921456462821518061")) //nolint:lll

	poseidon2 := NewPoseidon2(RoundConstants)

	bigArray4 := []*ff.Element{b1, b2, b0}

	for i := 0; i < b.N; i++ {
		poseidon2.Permutation(bigArray4) //nolint:errcheck,gosec
	}
}
