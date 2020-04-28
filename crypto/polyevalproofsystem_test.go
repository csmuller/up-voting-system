package crypto

import (
	"math/big"
	"testing"
)

func TestPolyEvalProofSystem(t *testing.T) {
	o, _ := new(big.Int).SetString(oTest, 10)
	p, _ := new(big.Int).SetString(pTest, 10)
	q, _ := new(big.Int).SetString(qTest, 10)

	gP := NewGStarModPrime(o, p)
	commP := NewPedersenCommitmentScheme(gP, gP.RandomGenerator(),
		[]*big.Int{gP.RandomGenerator()})

	gQ := NewGStarModPrime(p, q)
	commQ := NewPedersenCommitmentScheme(gQ, gQ.RandomGenerator(),
		[]*big.Int{gQ.RandomGenerator(), gQ.RandomGenerator()}) // comm_q can take two messages

	voter1 := GenerateNewVoter(commQ)
	voter2 := GenerateNewVoter(commQ)
	voter3 := GenerateNewVoter(commQ)

	coeffs := []*big.Int{big.NewInt(1)}
	poly := NewPolynomial(coeffs, gP.ZModOrder())
	poly = poly.IncludeCredential(voter1.U)
	poly = poly.IncludeCredential(voter2.U)
	poly = poly.IncludeCredential(voter3.U)

	// commitment c
	commToPubRand := commP.G.ZModOrder().RandomElement()
	commToPub := commP.Commit(commToPubRand, voter1.U)

	// pi_1
	ps := NewPolynomialEvaluationProofSystem(commP, poly)
	proof := ps.Generate(voter1.U, commToPubRand, commToPub, "yes")
	if !ps.Verify(proof, commToPub, "yes") {
		t.Fail()
	}
}
