// from github.com/iden3/go-iden3-crypto/ff/poseidon

package poseidon2

import "github.com/iden3/go-iden3-crypto/ff"

type Poseidon2 struct {
	params *Poseidon2Params
}

func NewPoseidon2(param *Poseidon2Params) *Poseidon2 {
	return &Poseidon2{params: param}
}

func (p *Poseidon2) GetT() int {
	return p.params.t
}

func (p *Poseidon2) Permutation(input []*ff.Element) []*ff.Element {
	t := p.GetT()
	if len(input) != t {
		panic("input length must equal t")
	}

	currentState := make([]*ff.Element, len(input))
	copy(currentState, input)

	// Linear layer at beginning
	p.MatmulExternal(currentState)

	for r := 0; r < p.params.roundsFBeginning; r++ {
		currentState = p.AddRc(currentState, p.params.RoundConstants[r])
		currentState = p.Sbox(currentState)
		p.MatmulExternal(currentState)
	}

	pEnd := p.params.roundsFBeginning + p.params.roundsP
	for r := p.params.roundsFBeginning; r < pEnd; r++ {
		currentState[0].Add(currentState[0], &p.params.RoundConstants[r][0])
		currentState[0] = p.SboxP(currentState[0])
		p.MatmulInternal(currentState, p.params.MatInternalDiagM1)
	}

	for r := pEnd; r < p.params.rounds; r++ {
		currentState = p.AddRc(currentState, p.params.RoundConstants[r])
		currentState = p.Sbox(currentState)
		p.MatmulExternal(currentState)
	}
	return currentState
}

func (p *Poseidon2) Sbox(input []*ff.Element) []*ff.Element {
	output := make([]*ff.Element, len(input))
	for i, el := range input {
		output[i] = p.SboxP(el)
	}
	return output
}

func (p *Poseidon2) SboxP(input *ff.Element) *ff.Element {
	input2 := input
	input2.Square(input2)

	switch p.params.d {
	case 3:
		out := input2
		out.Mul(out, input)
		return out
	case 5:
		out := input2
		out.Square(out)
		out.Mul(out, input)
		return out
	case 7:
		out := input2
		out.Square(out)
		out.Mul(out, input2)
		out.Mul(out, input)
		return out
	default:
		panic("invalid sbox degree")
	}
}

func (p *Poseidon2) MatmulExternal(input []*ff.Element) {
	t := p.GetT()

	switch t {
	case 2:
		sum := input[0]
		sum.Add(sum, input[1])
		input[0].Add(input[0], sum)
		input[1].Add(input[1], sum)
	case 3:
		sum := input[0]
		sum.Add(sum, input[1])
		sum.Add(sum, input[2])
		input[0].Add(input[0], sum)
		input[1].Add(input[1], sum)
		input[2].Add(input[2], sum)
	case 4:
		p.MatmulM4(input)
	case 8, 12, 16, 20, 24:
		p.MatmulM4(input)
		t4 := t / 4
		stored := make([]*ff.Element, 4)
		for l := 0; l < 4; l++ {
			stored[l] = input[l]
			for j := 1; j < t4; j++ {
				stored[l].Add(stored[l], input[4*j+l])
			}
		}
		for i := 0; i < len(input); i++ {
			input[i].Add(input[i], stored[i%4])
		}
	default:
		panic("invalid t value")
	}
}

func (p *Poseidon2) MatmulM4(input []*ff.Element) {
	t := p.GetT()
	t4 := t / 4
	for i := 0; i < t4; i++ {
		startIndex := i * 4
		t0 := input[startIndex]
		t0.Add(t0, input[startIndex+1])
		t1 := input[startIndex+2]
		t1.Add(t1, input[startIndex+3])
		t2 := input[startIndex+1]
		t2.Double(t2)
		t2.Add(t2, t1)
		t3 := input[startIndex+3]
		t3.Double(t3)
		t3.Add(t3, t0)
		t4 := t1
		t4.Double(t4)
		t4.Double(t4)
		t4.Add(t4, t3)
		t5 := t0
		t5.Double(t5)
		t5.Double(t5)
		t5.Add(t5, t2)
		t6 := t3
		t6.Add(t6, t5)
		t7 := t2
		t7.Add(t7, t4)
		input[startIndex] = t6
		input[startIndex+1] = t5
		input[startIndex+2] = t7
		input[startIndex+3] = t4
	}
}

func (p *Poseidon2) MatmulInternal(input []*ff.Element, matInternalDiagM1 []ff.Element) {
	t := p.GetT()

	switch t {
	case 2:
		sum := input[0]
		sum.Add(sum, input[1])
		input[0].Add(input[0], sum)
		input[1].Double(input[1])
		input[1].Add(input[1], sum)
	case 3:
		sum := input[0]
		sum.Add(sum, input[1])
		sum.Add(sum, input[2])
		input[0].Add(input[0], sum)
		input[1].Add(input[1], sum)
		input[2].Double(input[2])
		input[2].Add(input[2], sum)
	case 4, 8, 12, 16, 20, 24:
		sum := input[0]
		for _, el := range input[1:] {
			sum.Add(sum, el)
		}
		for i := 0; i < len(input); i++ {
			input[i].Mul(input[i], &matInternalDiagM1[i])
			input[i].Add(input[i], sum)
		}
	default:
		panic("invalid t value")
	}
}

func (p *Poseidon2) AddRc(input []*ff.Element, rc []ff.Element) []*ff.Element {
	output := make([]*ff.Element, len(input))
	for i := range input {
		output[i] = input[i]
		output[i].Add(output[i], &rc[i])
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
