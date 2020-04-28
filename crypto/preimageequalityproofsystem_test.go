package crypto

import (
	"math/big"
	"testing"
)

func TestPreimgEqProofSystem(t *testing.T) {
	p, _ := new(big.Int).SetString(pTest, 10)
	q, _ := new(big.Int).SetString(qTest, 10)

	gQ := NewGStarModPrime(p, q)
	commQ := NewPedersenCommitmentScheme(gQ, gQ.RandomGenerator(),
		[]*big.Int{gQ.RandomGenerator(), gQ.RandomGenerator()}) // comm_q can take two messages

	hHat := gQ.RandomElement()

	voter := GenerateNewVoter(commQ)
	uHat := gQ.Exp(hHat, voter.B)

	// commitment d
	commToAandBRand := commQ.G.ZModOrder().RandomElement()
	commToAandB := commQ.Commit(commToAandBRand, voter.A, voter.B)

	ps := NewPreimageEqualityProofSystem(hHat, commQ)
	proof := ps.Generate(voter, commToAandB, commToAandBRand, uHat, "yes")
	if !ps.Verify(proof, commToAandB, uHat, "yes") {
		t.Fail()
	}
}
