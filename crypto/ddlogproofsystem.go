package crypto

import "C"
import (
	"crypto/sha256"
	"encoding/json"
	"github.com/tendermint/go-amino"
	"math/big"
	"time"
)

// DoubleDiscreteLogProofSystem is used to prove knowledge of a representation of a committed value
// without revealing the representation nor the committed value. The proof is made non-interactive.
// It is based on the paper "Proof-of-Knowledge of Representation of Committed Value and Its
// Applications" by Man Ho Au et al.
// The class name is due to it being a generalization of double discrete log proofs.
type DoubleDiscreteLogProofSystem struct {
	CommSchemeInGp PedersenCommitmentScheme // The commitment scheme used for the committed value
	CommSchemeInGq PedersenCommitmentScheme // The commitment scheme used for the representation.
	SecurityParam  int                      // Security parameter determining the security level of the proof system.
	zp             ZModPrime
	zq             ZModPrime
	gp             GStarModPrime
	gq             GStarModPrime
}

// NewDoubleDiscreteLogProofSystem creates a new instance of the proof system with the given
// commitment schemes and security parameter.
func NewDoubleDiscreteLogProofSystem(commSchemeInGp PedersenCommitmentScheme,
	commSchemeInGq PedersenCommitmentScheme, securityParam int) DoubleDiscreteLogProofSystem {

	p := commSchemeInGp.G.Order
	q := commSchemeInGq.G.Order
	// Check if p = rq + 1
	if new(big.Int).GCD(nil, nil, new(big.Int).Sub(p, big.NewInt(1)), q).Cmp(q) != 0 {
		panic("Order p of cyclic group G_p must satisfy p = bq + 1 for some r and order q of" +
			" cyclic group G_q.")
	}
	// Check if 2^k < p, where k is the security parameter.
	if p.Cmp(new(big.Int).Exp(big.NewInt(2), big.NewInt(int64(securityParam)), nil)) < 1 {
		panic("Order p of cyclic group G_p must be bigger than 2^k where k is the security" +
			" parameter.")
	}

	return DoubleDiscreteLogProofSystem{
		CommSchemeInGp: commSchemeInGp,
		CommSchemeInGq: commSchemeInGq,
		SecurityParam:  securityParam,
		zp:             commSchemeInGp.G.ZModOrder(),
		zq:             commSchemeInGq.G.ZModOrder(),
		gp:             commSchemeInGp.G,
		gq:             commSchemeInGq.G,
	}
}

// DdLogProof represents a proof transcript for a proof of known representation of a
// committed value.
type DdLogProof struct {
	T     *big.Int
	T1Arr []*big.Int // Length equal to the security parameter.
	T2Arr []*big.Int // Length equal to the security parameter.

	ZX    *big.Int
	ZR    *big.Int
	ZMArr [][]*big.Int // A 2D array with len(zMArr) == securityParam and len(Zm[i]) == 2.
	ZSArr []*big.Int   // Length equal to the security parameter.
	ZRArr []*big.Int   // Length equal to the security parameter.
}

// Generate generates a proof of known representation of committed value.
// The parameters include the voters public and private credentials, commitments to the these
// credentials, the randomness values used in the commitments and the voter's vote.
func (ps *DoubleDiscreteLogProofSystem) Generate(voter Voter, commToU *big.Int,
	commToURand *big.Int, commToAandB *big.Int, commToAandBRand *big.Int,
	vote string) DdLogProof {

	defer LogExecutionTime(time.Now(), "double discrete log proof generation")

	// 1. Create commitment
	rhoX := ps.zp.RandomElement()
	rhoR := ps.zp.RandomElement()
	t := ps.CommSchemeInGp.Commit(rhoR, rhoX)

	rhoMArr := make([][]*big.Int, ps.SecurityParam)
	for i := range rhoMArr {
		rhoMArr[i] = make([]*big.Int, len(ps.CommSchemeInGq.Hm))
	}
	rhoSArr := make([]*big.Int, ps.SecurityParam)
	rhoRArr := make([]*big.Int, ps.SecurityParam)
	t1Arr := make([]*big.Int, ps.SecurityParam)
	t2Arr := make([]*big.Int, ps.SecurityParam)

	for i := 0; i < ps.SecurityParam; i++ {
		rhoSArr[i] = ps.zq.RandomElement()
		rhoRArr[i] = ps.zp.RandomElement()
		hProduct := big.NewInt(1)
		for j := 0; j < len(ps.CommSchemeInGq.Hm); j++ {
			rhoMArr[i][j] = ps.zq.RandomElement()
			hExp := ps.gq.Exp(ps.CommSchemeInGq.Hm[j], rhoMArr[i][j])
			hProduct = ps.gq.Mul(hProduct, hExp)
		}
		t1Arr[i] = ps.CommSchemeInGp.Commit(rhoRArr[i], hProduct)
		t2Arr[i] = ps.CommSchemeInGq.Commit(rhoSArr[i], rhoMArr[i]...)
	}

	// 2. Create challenge
	// The challenge is a hash of all inputs modulo the p of Z_p. Later in the proof generation the
	// challenge's k first bits are used individually, where k is the security parameter. Therefore,
	// it is important that p is at least 2^k. It doesn't matter if the challenge has more than k
	// bits. In both, the proof generation and verification we only use the k first bits.
	ch := ps.generateChallenge(commToU, commToAandB, t, t1Arr, t2Arr, vote)

	// 3. Create response
	x := voter.U                      // voter credential u
	r := commToURand                  // randomness in commitment to u
	s := commToAandBRand              // randomness in commitment to a and b
	m := []*big.Int{voter.A, voter.B} // voter private credentials a and b

	zX := ps.zp.Add(rhoX, ps.zp.AdditiveInvert(ps.zp.Mul(x, ch)))
	zR := ps.zp.Add(rhoR, ps.zp.AdditiveInvert(ps.zp.Mul(r, ch)))

	zMArr := make([][]*big.Int, ps.SecurityParam)
	zSArr := make([]*big.Int, ps.SecurityParam)
	zRArr := make([]*big.Int, ps.SecurityParam)

	// len(m) and len(ps.CommSchemeInGq.Hm) must be the same since len(m) is the number of messages
	// in commitment scheme with group G_q and len(ps.CommSchemeInGq. Hm) is the number of message
	// generators in that same scheme.
	for i := 0; i < ps.SecurityParam; i++ {
		zMiArr := make([]*big.Int, len(m))
		bit := big.NewInt(int64(ch.Bit(i)))
		for j := 0; j < len(m); j++ {
			zMiArr[j] = ps.zq.Add(rhoMArr[i][j], ps.zq.AdditiveInvert(ps.zq.Mul(m[j], bit)))
		}
		zMArr[i] = zMiArr
		zSArr[i] = ps.zq.Add(rhoSArr[i], ps.zq.AdditiveInvert(ps.zq.Mul(s, bit)))
		hProduct := big.NewInt(1)
		for j := 0; j < len(ps.CommSchemeInGq.Hm); j++ {
			hExp := ps.gq.Exp(ps.CommSchemeInGq.Hm[j], zMArr[i][j])
			hProduct = ps.gq.Mul(hProduct, hExp)
		}
		zRArr[i] = ps.zp.Add(rhoRArr[i], ps.zp.AdditiveInvert(ps.zp.Mul(ps.zp.Mul(bit, hProduct), r)))
	}

	return DdLogProof{t, t1Arr, t2Arr, zX, zR, zMArr, zSArr, zRArr}
}

// Verify verifies a proof of known representation of committed values. Next to the proof
// transcript, the committed value and the commitment to the representation are the input
// parameters.
func (ps *DoubleDiscreteLogProofSystem) Verify(proof DdLogProof, commToU *big.Int,
	commToAandB *big.Int, vote string) bool {

	defer LogExecutionTime(time.Now(), "double discrete log proof verification")

	t := proof.T
	t1 := proof.T1Arr
	t2 := proof.T2Arr

	zR := proof.ZR
	zX := proof.ZX
	zMArr := proof.ZMArr
	zSArr := proof.ZSArr
	zRArr := proof.ZRArr

	// 2. Create challenge
	ch := ps.generateChallenge(commToU, commToAandB, proof.T, proof.T1Arr, proof.T2Arr, vote)

	comm := ps.CommSchemeInGp.Commit(zR, zX)
	v := t.Cmp(ps.gp.Mul(ps.gp.Exp(commToU, ch), comm)) == 0

	for i := 0; i < ps.SecurityParam; i++ {
		bit := big.NewInt(int64(ch.Bit(i)))
		// T2
		comm := ps.CommSchemeInGq.Commit(zSArr[i], zMArr[i]...)
		v = v && t2[i].Cmp(ps.gq.Mul(ps.gq.Exp(commToAandB, bit), comm)) == 0
		// T1
		hProduct := big.NewInt(1)
		for j := 0; j < len(ps.CommSchemeInGq.Hm); j++ {
			hExp := ps.gq.Exp(ps.CommSchemeInGq.Hm[j], zMArr[i][j])
			hProduct = ps.gq.Mul(hProduct, hExp)
		}
		if bit.Cmp(big.NewInt(0)) == 0 {
			v = v && t1[i].Cmp(ps.CommSchemeInGp.Commit(zRArr[i], hProduct)) == 0
		} else {
			g := ps.CommSchemeInGp.Hr
			temp := ps.gp.Mul(ps.gp.Exp(commToU, hProduct), ps.gp.Exp(g, zRArr[i]))
			v = v && t1[i].Cmp(temp) == 0
		}
	}

	return v
}

func (ps *DoubleDiscreteLogProofSystem) generateChallenge(commToU, commToAandB, t *big.Int, t1Arr,
	t2Arr []*big.Int, vote string) *big.Int {
	inputs := [][]*big.Int{{commToU, commToAandB, t}, t1Arr, t2Arr}
	sha := sha256.New()
	for _, c := range inputs {
		for _, elem := range c {
			sha.Write(elem.Bytes())
		}
	}
	sha.Write([]byte(vote))
	hash := sha.Sum(nil)
	ch := new(big.Int).SetBytes(hash)
	return ch.Mod(ch, ps.zp.Modulus)
}

// ddLogProofDTO is needed for Tendermint serialization and deserialization.
type ddLogProofDTO struct {
	T     Int   `json:"t"`
	T1Arr []Int `json:"t1_ar" ` // Length equal to the security parameter.
	T2Arr []Int `json:"t2_arr"` // Length equal to the security parameter.

	ZX    Int     `json:"zx"`
	ZR    Int     `json:"zr"`
	ZMArr [][]Int `json:"zm_arr"` // A 2D array with length (securityParam, 2).
	ZSArr []Int   `json:"zs_arr"` // Length equal to the security parameter.
	ZRArr []Int   `json:"zr_arr"` // Length equal to the security parameter.
}

func (p DdLogProof) MarshalAmino() (string, error) {
	dto := ddLogProofDTO{}
	p.wrapInDTO(&dto)
	return string(amino.MustMarshalBinaryBare(dto)), nil
}

func (p *DdLogProof) UnmarshalAmino(bytes []byte) error {
	var dto ddLogProofDTO
	amino.MustUnmarshalBinaryBare(bytes, &dto)
	p.unwrapDTO(dto)
	return nil
}

func (p DdLogProof) MarshalJSON() ([]byte, error) {
	dto := ddLogProofDTO{}
	p.wrapInDTO(&dto)
	return json.Marshal(dto)
}

func (p *DdLogProof) UnmarshalJSON(bytes []byte) error {
	var dto ddLogProofDTO
	err := amino.UnmarshalJSON(bytes, &dto)
	if err != nil {
		return err
	}
	p.unwrapDTO(dto)
	return nil
}

func (p DdLogProof) wrapInDTO(dto *ddLogProofDTO) {
	dto.T = NewInt(p.T)
	dto.T1Arr = make([]Int, len(p.T1Arr))
	for i, v := range p.T1Arr {
		dto.T1Arr[i] = NewInt(v)
	}
	dto.T2Arr = make([]Int, len(p.T2Arr))
	for i, v := range p.T2Arr {
		dto.T2Arr[i] = NewInt(v)
	}
	dto.ZX = NewInt(p.ZX)
	dto.ZR = NewInt(p.ZR)
	dto.ZMArr = make([][]Int, len(p.ZMArr))
	for i, arr := range p.ZMArr {
		dto.ZMArr[i] = make([]Int, len(p.ZMArr[i]))
		for j, v := range arr {
			dto.ZMArr[i][j] = NewInt(v)
		}
	}
	dto.ZSArr = make([]Int, len(p.ZSArr))
	for i, v := range p.ZSArr {
		dto.ZSArr[i] = NewInt(v)
	}
	dto.ZRArr = make([]Int, len(p.ZRArr))
	for i, v := range p.ZRArr {
		dto.ZRArr[i] = NewInt(v)
	}
}

func (p *DdLogProof) unwrapDTO(dto ddLogProofDTO) {
	p.T = dto.T.BigInt()
	p.T1Arr = make([]*big.Int, len(dto.T1Arr))
	for i, v := range dto.T1Arr {
		p.T1Arr[i] = v.BigInt()
	}
	p.T2Arr = make([]*big.Int, len(dto.T2Arr))
	for i, v := range dto.T2Arr {
		p.T2Arr[i] = v.BigInt()
	}
	p.ZX = dto.ZX.BigInt()
	p.ZR = dto.ZR.BigInt()
	p.ZMArr = make([][]*big.Int, len(dto.ZMArr))
	for i, arr := range dto.ZMArr {
		p.ZMArr[i] = make([]*big.Int, len(dto.ZMArr[i]))
		for j, v := range arr {
			p.ZMArr[i][j] = v.BigInt()
		}
	}
	p.ZSArr = make([]*big.Int, len(dto.ZSArr))
	for i, v := range dto.ZSArr {
		p.ZSArr[i] = v.BigInt()
	}
	p.ZRArr = make([]*big.Int, len(dto.ZRArr))
	for i, v := range dto.ZRArr {
		p.ZRArr[i] = v.BigInt()
	}
}

func (p DdLogProof) String() string {
	return ""
}
