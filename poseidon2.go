// from github.com/iden3/go-iden3-crypto/ff/poseidon

package poseidon2

import "github.com/iden3/go-iden3-crypto/ff"

// const NROUNDSF = 8 //nolint:golint

// var NROUNDSP = 56    //nolint:golint
// const MAX_WIDTH = 18 // len(NROUNDSP)+2

// func zero() *ff.Element {
// 	return ff.NewElement()
// }

// // exp5 performs x^5 mod p
// // https://eprint.iacr.org/2019/458.pdf page 8
// func exp5(a *ff.Element) {
// 	a.Exp(*a, big.NewInt(5)) //nolint:gomnd
// }

// // exp5state perform exp5 for whole state
// func exp5state(state []ff.Element, t int) {
// 	for i := 0; i < t; i++ {
// 		exp5(&state[i])
// 	}
// }

// // ark computes Add-Round Key, from the paper https://eprint.iacr.org/2019/458.pdf
// func ark(state []ff.Element, c []*ff.Element, it int, t int) {
// 	for i := 0; i < t; i++ {
// 		state[i].Add(&state[i], c[it+i])
// 	}
// }

// // mix returns [[matrix]] * [vector]
// func mix(state []ff.Element, t int, m [][]*ff.Element) []ff.Element {
// 	mul := zero()
// 	newState := make([]ff.Element, MAX_WIDTH)
// 	for i := 0; i < t; i++ {
// 		newState[i].SetUint64(0)
// 		for j := 0; j < t; j++ {
// 			mul.Mul(m[j][i], &state[j])
// 			newState[i].Add(&newState[i], mul)
// 		}
// 	}

// 	for i := 0; i < t; i++ {
// 		state[i].Set(&newState[i])
// 	}
// 	return state
// }

// func permute(state []ff.Element, t int) []ff.Element {

// 	nRoundsF := NROUNDSF
// 	nRoundsP := NROUNDSP[t-2]
// 	C := c.c[t-2]
// 	S := c.s[t-2]
// 	M := c.m[t-2]
// 	P := c.p[t-2]

// 	ark(state, C, 0, t)

// 	for i := 0; i < nRoundsF/2-1; i++ {
// 		exp5state(state, t)
// 		ark(state, C, (i+1)*t, t)
// 		state = mix(state, t, M)
// 	}
// 	exp5state(state, t)
// 	ark(state, C, (nRoundsF/2)*t, t)
// 	state = mix(state, t, P)

// 	for i := 0; i < nRoundsP; i++ {
// 		exp5(&state[0])
// 		state[0].Add(&state[0], C[(nRoundsF/2+1)*t+i])

// 		mul := zero()
// 		newState0 := zero()
// 		for j := 0; j < t; j++ {
// 			mul.Mul(S[(t*2-1)*i+j], &state[j])
// 			newState0.Add(newState0, mul)
// 		}

// 		for k := 1; k < t; k++ {
// 			state[k].Add(&state[k], mul.Mul(&state[0], S[(t*2-1)*i+t+k-1]))
// 		}
// 		state[0].Set(newState0)
// 	}

// 	for i := 0; i < nRoundsF/2-1; i++ {
// 		exp5state(state, t)
// 		ark(state, C, (nRoundsF/2+1)*t+nRoundsP+i*t, t)
// 		state = mix(state, t, M)
// 	}
// 	exp5state(state, t)
// 	return mix(state, t, M)
// }

// // for short, use size of inpBI as cap
// func Hash(inpBI []*big.Int, width int) (*big.Int, error) {
// 	return HashWithCap(inpBI, width, int64(len(inpBI)))
// }

// // Hash using possible sponge specs specified by width (rate from 1 to 15), the size of input is applied as capacity
// // (notice we do not include width in the capacity )
// func HashWithCap(inpBI []*big.Int, width int, nBytes int64) (*big.Int, error) {
// 	if width < 2 {
// 		return nil, fmt.Errorf("width must be ranged from 2 to 16")
// 	}
// 	if width > MAX_WIDTH {
// 		return nil, fmt.Errorf("invalid inputs width %d, max %d", width, MAX_WIDTH) //nolint:gomnd,lll
// 	}

// 	// capflag = nBytes * 2^64
// 	pow64 := big.NewInt(1)
// 	pow64.Lsh(pow64, 64)
// 	capflag := ff.NewElement().SetBigInt(big.NewInt(nBytes))
// 	capflag.Mul(capflag, ff.NewElement().SetBigInt(pow64))

// 	// initialize the state
// 	state := make([]ff.Element, MAX_WIDTH)
// 	state[0] = *capflag

// 	rate := width - 1
// 	i := 0
// 	// always perform one round of permutation even when input is empty
// 	for {
// 		// each round absorb at most `rate` elements from `inpBI`
// 		for j := 0; j < rate && i < len(inpBI); i, j = i+1, j+1 {
// 			state[j+1].Add(&state[j+1], ff.NewElement().SetBigInt(inpBI[i]))
// 		}
// 		state = permute(state, width)
// 		if i == len(inpBI) {
// 			break
// 		}
// 	}

// 	// squeeze
// 	rE := state[0]
// 	r := big.NewInt(0)
// 	rE.ToBigIntRegular(r)
// 	return r, nil

// }

// // Hash computes the Poseidon hash for the given fixed-size inputs, with specified domain field
// func HashFixedWithDomain(inpBI []*big.Int, domain *big.Int) (*big.Int, error) {
// 	t := len(inpBI) + 1
// 	if len(inpBI) == 0 || len(inpBI) > len(NROUNDSP) {
// 		return nil, fmt.Errorf("invalid inputs length %d, max %d", len(inpBI), len(NROUNDSP)) //nolint:gomnd,lll
// 	}
// 	if !utils.CheckBigIntArrayInField(inpBI[:]) {
// 		return nil, errors.New("inputs values not inside Finite Field")
// 	}
// 	inp := make([]ff.Element, MAX_WIDTH)
// 	for idx, bi := range inpBI {
// 		inp[idx].SetBigInt(bi)
// 	}

// 	state := make([]ff.Element, MAX_WIDTH)
// 	state[0] = *ff.NewElement().SetBigInt(domain)
// 	copy(state[1:], inp[:])

// 	state = permute(state, t)

// 	rE := state[0]
// 	r := big.NewInt(0)
// 	rE.ToBigIntRegular(r)
// 	return r, nil
// }

// // Deprecated HashFixed entry, with domain field is 0
// func HashFixed(inpBI []*big.Int) (*big.Int, error) {
// 	log.Warn("called a deprecated method for poseidon fixed hash", "inputs", inpBI)
// 	return HashFixedWithDomain(inpBI, big.NewInt(0))
// }

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
