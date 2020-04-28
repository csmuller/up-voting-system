package crypto

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/tendermint/go-amino"
	"math"
	"math/big"
	"time"
)

// PolynomialEvaluationProofSystem is the proof system used for set membership proofs.
// It is used to prove that a voter's public credential is in the set of eligible voters.
// It is based on the work "A practical system for globally revoking the unlinkable pseudonyms
// of unknown users" by Stefan Brands et al.
type PolynomialEvaluationProofSystem struct {
	CommScheme PedersenCommitmentScheme // The scheme used to commit to the public credential u.
	Polynomial Polynomial               // The credential polynomial containing all eligible voters.
	gStarModPr GStarModPrime
	zModPr     ZModPrime
	// d is calculated from the order D of the polynomial.
	// D = 2^(d+1) - 1 ==> d = ceil(log( D + 1)) - 1 = floor( log(D))
	d int
}

// PolyEvalProof represents a transcript of a polynomial evaluation proof.
type PolyEvalProof struct {
	CArr   []*big.Int // Length d. Doesn't include the commitment to credential u.
	CfArr  []*big.Int // Length d + 1
	CdArr  []*big.Int // Length d + 1
	CfuArr []*big.Int // Length d

	FBarArr  []*big.Int // Length d + 1
	RBarArr  []*big.Int // Length d + 1
	TBar     *big.Int
	XiBarArr []*big.Int // d
}

// NewPolynomialEvaluationProofSystem creates a new instance of the proof system.
// The given commitment scheme is the one used to commit to the voter's public credential u.
// The polynomial is the credential polynomial containing all eligible voter's public credentials.
func NewPolynomialEvaluationProofSystem(commScheme PedersenCommitmentScheme,
	poly Polynomial) PolynomialEvaluationProofSystem {

	return PolynomialEvaluationProofSystem{
		CommScheme: commScheme,
		Polynomial: poly,
		gStarModPr: commScheme.G,
		zModPr:     commScheme.G.ZModOrder(),
		d:          int(math.Floor(math.Log(float64(poly.Degree())) / math.Log(2))),
	}
}

var (
	// The credential polynomial evaluates Poly(u) = v = 0 for every public credential u that was
	// added to it. The polynomial evaluation proof system expects commitments to u and to v as
	// public input. Because v is always 0 the randomness (also input to the commitment) is chosen
	// to be 0 as well and the commitment thereby takes on the value 1.
	// These two values are required in the proof generation and verification.
	commToVRandomness = big.NewInt(0)
	commToV           = big.NewInt(1)
)

// Generate generates a polynomial evaluation proof for the given voter public credential u.
// Provide a commitment commToU to the voter's public credential u with a random commitment element
// r.
func (ps *PolynomialEvaluationProofSystem) Generate(u, r, commToU *big.Int,
	vote string) PolyEvalProof {

	defer LogExecutionTime(time.Now(), "polynomial evaluation proof generation")

	rArr := make([]*big.Int, ps.d+1)
	fArr := make([]*big.Int, ps.d+1)
	sArr := make([]*big.Int, ps.d+1)
	tArr := make([]*big.Int, ps.d+1)
	xiArr := make([]*big.Int, ps.d)
	for i := 0; i < ps.d; i++ {
		xiArr[i] = ps.zModPr.RandomElement()
	}
	for i := 0; i < ps.d+1; i++ {
		rArr[i] = ps.zModPr.RandomElement()
		fArr[i] = ps.zModPr.RandomElement()
		sArr[i] = ps.zModPr.RandomElement()
		tArr[i] = ps.zModPr.RandomElement()
	}
	rArr[0] = r

	uArr := make([]*big.Int, ps.d+1)
	uArr[0] = u
	for i := 1; i < len(uArr); i++ {
		uArr[i] = ps.zModPr.Mul(uArr[i-1], uArr[i-1])
	}

	// 1. Create commitments
	// a) c_1 ... c_d
	cArr := make([]*big.Int, ps.d+1)
	for i := 1; i < len(cArr); i++ {
		cArr[i] = ps.CommScheme.Commit(rArr[i], uArr[i])
	}
	cArr = cArr[1:]

	// b) c_f_0 ... c_f_d
	cfArr := make([]*big.Int, ps.d+1)
	for i := 0; i < len(cfArr); i++ {
		cfArr[i] = ps.CommScheme.Commit(sArr[i], fArr[i])
	}

	// c) c_delta_0 ... c_delta_d
	dArr := ps.calcDeltas(uArr, fArr)
	cdArr := make([]*big.Int, ps.d+1)
	for i := 0; i < len(cdArr); i++ {
		cdArr[i] = ps.CommScheme.Commit(tArr[i], dArr[i])
	}

	// d) c_fu_0 ... c_fu_d-1
	cfuArr := make([]*big.Int, ps.d)
	for i := 0; i < len(cfuArr); i++ {
		cfuArr[i] = ps.CommScheme.Commit(xiArr[i], ps.zModPr.Mul(fArr[i], uArr[i]))
	}

	publicInput := []*big.Int{commToU, commToV}
	ch := ps.generateChallenge(vote, publicInput, cArr, cfArr, cdArr, cfuArr)

	// Response 1 & 2
	fBarArr := make([]*big.Int, ps.d+1)
	rBarArr := make([]*big.Int, ps.d+1)
	for i := 0; i < len(fBarArr); i++ {
		fBarArr[i] = ps.zModPr.Add(ps.zModPr.Mul(uArr[i], ch), fArr[i])
		rBarArr[i] = ps.zModPr.Add(ps.zModPr.Mul(rArr[i], ch), sArr[i])
	}

	// Response 3
	dPlusOne := big.NewInt(int64(ps.d + 1))
	tBar := ps.zModPr.Mul(ps.zModPr.Exp(ch, dPlusOne), commToVRandomness)
	xi := ps.zModPr.Exp(ch, big.NewInt(0))
	for i := 0; i <= ps.d; i++ {
		tBar = ps.zModPr.Add(tBar, ps.zModPr.Mul(tArr[i], xi))
		xi = ps.zModPr.Mul(xi, ch)
	}

	// Response 4
	xiBarArr := make([]*big.Int, ps.d)
	for i := 0; i < len(xiBarArr); i++ {
		term1 := ps.zModPr.Mul(rArr[i+1], ch)
		term2 := ps.zModPr.AdditiveInvert(ps.zModPr.Mul(fBarArr[i], rArr[i]))
		xiBarArr[i] = ps.zModPr.Add(ps.zModPr.Add(term1, term2), xiArr[i])
	}

	return PolyEvalProof{cArr, cfArr, cdArr, cfuArr, fBarArr, rBarArr, tBar, xiBarArr}
}

var xFactorPoly Polynomial

func (ps *PolynomialEvaluationProofSystem) calcDeltas(uArr, fArr []*big.Int) []*big.Int {
	ring := ps.Polynomial.ZModPr
	// This is X^i[j] of the product
	xFactorPoly = NewPolynomial([]*big.Int{big.NewInt(0), big.NewInt(1)}, ring)
	// This polynomial is used and extended in each step of the
	currPoly := NewPolynomial([]*big.Int{big.NewInt(1)}, ring)
	var result Polynomial
	result = ps.calcDeltaPolynomial(ps.d+1, 0, currPoly, uArr, fArr, result)
	return result.Coeffs
}

func (ps *PolynomialEvaluationProofSystem) calcDeltaPolynomial(lvl int, deg int,
	currPoly Polynomial, uArr []*big.Int, fArr []*big.Int, result Polynomial) Polynomial {

	if lvl == 0 {
		if deg == 0 {
			result = currPoly
		}
		return result.Add(currPoly.MulScalar(ps.Polynomial.Coeffs[deg]))
	} else {
		result = ps.calcDeltaPolynomial(lvl-1, deg, currPoly.Mul(xFactorPoly), uArr, fArr, result)

		nextDeg := deg + int(math.Pow(2, float64(lvl-1)))
		if nextDeg <= ps.Polynomial.Degree() {
			xuFactorPoly := NewPolynomial([]*big.Int{fArr[lvl-1], uArr[lvl-1]}, ps.Polynomial.ZModPr)
			result = ps.calcDeltaPolynomial(lvl-1, nextDeg, currPoly.Mul(xuFactorPoly), uArr, fArr,
				result)
		}
		return result
	}
}

// Verify verifies the given proof transcript.
func (ps *PolynomialEvaluationProofSystem) Verify(proof PolyEvalProof, commToU *big.Int,
	vote string) bool {

	defer LogExecutionTime(time.Now(), "polynomial evaluation proof verification")

	cArr := proof.CArr
	cfArr := proof.CfArr
	cdArr := proof.CdArr
	cfuArr := proof.CfuArr

	fBarArr := proof.FBarArr
	rBarArr := proof.RBarArr
	tBar := proof.TBar
	xiBarArr := proof.XiBarArr

	publicInput := []*big.Int{commToU, commToV}
	ch := ps.generateChallenge(vote, publicInput, cArr, cfArr, cdArr, cfuArr)

	cArr = append([]*big.Int{commToU}, cArr...)

	// Pre-compute c_j^x
	cxArr := make([]*big.Int, len(cArr))
	for i := 0; i < len(cArr); i++ {
		cxArr[i] = ps.gStarModPr.Exp(cArr[i], ch)
	}

	v := true
	for i := 0; i < ps.d+1 && v; i++ {
		comm := ps.CommScheme.Commit(rBarArr[i], fBarArr[i])
		v = v && (ps.gStarModPr.Mul(cxArr[i], cfArr[i]).Cmp(comm) == 0)
	}

	zero := big.NewInt(0)
	for i := 0; i < ps.d && v; i++ {
		comm := ps.CommScheme.Commit(xiBarArr[i], zero)
		cExpF := ps.gStarModPr.Exp(cArr[i], ps.zModPr.AdditiveInvert(fBarArr[i]))
		v = v && ps.gStarModPr.Mul(ps.gStarModPr.Mul(cxArr[i+1], cExpF), cfuArr[i]).Cmp(comm) == 0
	}

	left := ps.gStarModPr.Exp(commToV, ps.zModPr.Exp(ch, big.NewInt(int64(ps.d+1))))
	xi := ps.zModPr.Exp(ch, zero)
	for i := 0; i <= ps.d; i++ {
		left = ps.gStarModPr.Mul(left, ps.gStarModPr.Exp(cdArr[i], xi))
		xi = ps.zModPr.Mul(xi, ch)
	}

	dBar := ps.calcDeltaBar(fBarArr, ch)
	v = v && left.Cmp(ps.CommScheme.Commit(tBar, dBar)) == 0

	return v
}

func (ps *PolynomialEvaluationProofSystem) calcDeltaBar(fBarArr []*big.Int, ch *big.Int) *big.Int {
	result := big.NewInt(0)
	return ps.calcDeltaBarTreelike(ps.d+1, 0, big.NewInt(1), fBarArr, ch, result)
}

func (ps *PolynomialEvaluationProofSystem) calcDeltaBarTreelike(lvl int, deg int, value *big.Int,
	fBarArr []*big.Int, ch *big.Int, result *big.Int) *big.Int {
	if lvl == 0 {
		return ps.zModPr.Add(result, ps.zModPr.Mul(value, ps.Polynomial.Coeffs[deg]))
	} else {
		// right node
		result = ps.calcDeltaBarTreelike(lvl-1, deg, ps.zModPr.Mul(value, ch), fBarArr, ch,
			result)

		// left node
		nextDeg := deg + int(math.Pow(2, float64(lvl-1)))
		if nextDeg <= ps.Polynomial.Degree() {
			result = ps.calcDeltaBarTreelike(lvl-1, nextDeg, ps.zModPr.Mul(value, fBarArr[lvl-1]),
				fBarArr, ch, result)
		}
		return result
	}
}

func (ps *PolynomialEvaluationProofSystem) generateChallenge(vote string,
	inputs ...[]*big.Int) *big.Int {

	sha := sha256.New()
	for _, c := range inputs {
		for _, elem := range c {
			sha.Write(elem.Bytes())
		}
	}
	sha.Write([]byte(vote))
	hash := sha.Sum(nil)
	ch := new(big.Int).SetBytes(hash)
	return ch.Mod(ch, ps.zModPr.Modulus)
}

// polyEvalProofDTO is required for Tendermint serialization and deserialization.
type polyEvalProofDTO struct {
	CArr   []Int `json:"c"`   // Length d. Doesn't include the commitment to credential u.
	CfArr  []Int `json:"cf"`  // Length d + 1
	CdArr  []Int `json:"cd"`  // Length d + 1
	CfuArr []Int `json:"cfu"` // Length d

	FBarArr  []Int `json:"fBar"` // Length d + 1
	RBarArr  []Int `json:"rBar"` // Length d + 1
	TBar     Int   `json:"tBar"`
	XiBarArr []Int `json:"xiBar"` // d
}

func (p PolyEvalProof) MarshalAmino() (string, error) {
	dto := polyEvalProofDTO{}
	p.wrapInDTO(&dto)
	return string(amino.MustMarshalBinaryBare(dto)), nil
}

func (p *PolyEvalProof) UnmarshalAmino(bytes []byte) error {
	var dto polyEvalProofDTO
	amino.MustUnmarshalBinaryBare(bytes, &dto)
	p.unwrapDTO(dto)
	return nil
}

func (p PolyEvalProof) MarshalJSON() ([]byte, error) {
	dto := polyEvalProofDTO{}
	p.wrapInDTO(&dto)
	return json.Marshal(dto)
}

func (p *PolyEvalProof) UnmarshalJSON(bytes []byte) error {
	var dto polyEvalProofDTO
	err := amino.UnmarshalJSON(bytes, &dto)
	if err != nil {
		return err
	}
	p.unwrapDTO(dto)
	return nil
}

func (p PolyEvalProof) wrapInDTO(dto *polyEvalProofDTO) {
	dto.CArr = make([]Int, len(p.CArr))
	for i, v := range p.CArr {
		dto.CArr[i] = NewInt(v)
	}
	dto.CfArr = make([]Int, len(p.CfArr))
	for i, v := range p.CfArr {
		dto.CfArr[i] = NewInt(v)
	}
	dto.CdArr = make([]Int, len(p.CdArr))
	for i, v := range p.CdArr {
		dto.CdArr[i] = NewInt(v)
	}
	dto.CfuArr = make([]Int, len(p.CfuArr))
	for i, v := range p.CfuArr {
		dto.CfuArr[i] = NewInt(v)
	}
	dto.FBarArr = make([]Int, len(p.FBarArr))
	for i, v := range p.FBarArr {
		dto.FBarArr[i] = NewInt(v)
	}
	dto.RBarArr = make([]Int, len(p.RBarArr))
	for i, v := range p.RBarArr {
		dto.RBarArr[i] = NewInt(v)
	}
	dto.TBar = NewInt(p.TBar)
	dto.XiBarArr = make([]Int, len(p.RBarArr))
	for i, v := range p.XiBarArr {
		dto.XiBarArr[i] = NewInt(v)
	}
}

func (p *PolyEvalProof) unwrapDTO(dto polyEvalProofDTO) {
	p.CArr = make([]*big.Int, len(dto.CArr))
	for i, v := range dto.CArr {
		p.CArr[i] = v.BigInt()
	}
	p.CfArr = make([]*big.Int, len(dto.CfArr))
	for i, v := range dto.CfArr {
		p.CfArr[i] = v.BigInt()
	}
	p.CdArr = make([]*big.Int, len(dto.CdArr))
	for i, v := range dto.CdArr {
		p.CdArr[i] = v.BigInt()
	}
	p.CfuArr = make([]*big.Int, len(dto.CfuArr))
	for i, v := range dto.CfuArr {
		p.CfuArr[i] = v.BigInt()
	}
	p.FBarArr = make([]*big.Int, len(dto.FBarArr))
	for i, v := range dto.FBarArr {
		p.FBarArr[i] = v.BigInt()
	}
	p.RBarArr = make([]*big.Int, len(dto.RBarArr))
	for i, v := range dto.RBarArr {
		p.RBarArr[i] = v.BigInt()
	}
	p.TBar = dto.TBar.BigInt()
	p.XiBarArr = make([]*big.Int, len(dto.RBarArr))
	for i, v := range dto.XiBarArr {
		p.XiBarArr[i] = v.BigInt()
	}
}

func (p PolyEvalProof) String() string {
	return ""
}
