// from github.com/iden3/go-iden3-crypto/ff/poseidon

package poseidon2

import (
	"github.com/iden3/go-iden3-crypto/ff"
)

type Poseidon2 struct {
	RoundConstants [][]ff.Element
}

func NewPoseidon2(RoundConstants [][]ff.Element) *Poseidon2 {
	return &Poseidon2{RoundConstants: RoundConstants}
}

func (p *Poseidon2) Permutation(input []*ff.Element) []*ff.Element {

	currentState := make([]*ff.Element, len(input))
	copy(currentState, input)

	// Linear layer at beginning
	p.MatmulExternal(currentState)

	for r := 0; r < 4; r++ {
		currentState = p.AddRc(currentState, p.RoundConstants[r])
		currentState = p.Sbox(currentState)
		p.MatmulExternal(currentState)
	}

	pEnd := 4 + 56
	for r := 4; r < pEnd; r++ {
		currentState[0].Add(currentState[0], &p.RoundConstants[r][0])
		currentState[0] = p.SboxP(currentState[0])
		p.MatmulInternal(currentState)
	}

	for r := pEnd; r < 64; r++ {
		currentState = p.AddRc(currentState, p.RoundConstants[r])
		currentState = p.Sbox(currentState)
		p.MatmulExternal(currentState)
	}
	return currentState
}

func (p *Poseidon2) Sbox(input []*ff.Element) []*ff.Element {
	output := make([]*ff.Element, len(input))
	for i, el := range input {
		output[i] = new(ff.Element)
		output[i] = p.SboxP(el)
	}
	return output
}

func (p *Poseidon2) SboxP(input *ff.Element) *ff.Element {
	input2 := new(ff.Element)
	input2.Square(input)
	out := new(ff.Element)
	out.Square(input2)
	out.Mul(out, input)
	return out
}

// MatmulExternal is a matrix multiplication with the external matrix
// [2, 1, 1]
// [1, 2, 1]
// [1, 1, 2]
func (p *Poseidon2) MatmulExternal(input []*ff.Element) {
	sum := new(ff.Element)
	sum.Add(input[0], input[1])
	sum.Add(sum, input[2])
	input[0].Add(input[0], sum)
	input[1].Add(input[1], sum)
	input[2].Add(input[2], sum)
}

// MatmulInternal is a matrix multiplication with the internal matrix
// [2, 1, 1]
// [1, 2, 1]
// [1, 1, 3]
func (p *Poseidon2) MatmulInternal(input []*ff.Element) {
	sum := new(ff.Element)
	sum.Add(input[0], input[1])
	sum.Add(sum, input[2])
	input[0].Add(input[0], sum)
	input[1].Add(input[1], sum)
	input[2].Double(input[2])
	input[2].Add(input[2], sum)
}

func (p *Poseidon2) AddRc(input []*ff.Element, rc []ff.Element) []*ff.Element {
	output := make([]*ff.Element, len(input))
	for i := range input {
		output[i] = new(ff.Element)
		output[i].Add(input[i], &rc[i])
	}
	return output
}

// Implementing MerkleTreeHash interface
type MerkleTreeHash interface {
	Compress(input []*ff.Element) *ff.Element
}

func (p *Poseidon2) Compress(input []*ff.Element) *ff.Element {
	return p.Permutation([]*ff.Element{input[0], input[1], ff.NewElement()})[0]
}
