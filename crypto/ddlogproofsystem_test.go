package crypto

import (
	"math/big"
	"testing"
)

const (
	// Toy example for debugging.
	//q string = "11"
	//// p =  2 * q + 1
	//p string = "23"
	//// o = 4 * p + 1
	//o string = "47"

	// Q, 160 bit
	qTest string = "1081119563825030427708677600856959359670713108783"
	// P = a * Q + 1, 1024 bit
	pTest string = "132981118064499312972124229719551507064282251442693318094413647002876359530119444044769383265695686373097209253015503887096288112369989708235068428214124661556800389180762828009952422599372290980806417384771730325122099441368051976156139223257233269955912341167062173607119895128870594055324929155200165347329"
	// O = 981 * P + 1, 1034 bit
	oTest string = "130321495703209326712681745125160476922996606413839451732525374062818832339517055163873995600381772645635265067955193809354362350122589914070367059649842168325664381397147571449753374147384845161190289037076295718619657452540690936633016438792088604556794094343720930134977497226293182174218430572096162040382421"

	securityParam = 4
)

func TestDDLogProofSystem(t *testing.T) {
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

	// commitment c
	commToURand := commP.G.ZModOrder().RandomElement()
	commToU := commP.Commit(commToURand, voter1.U)

	// commitment d
	commToAandBRand := commQ.G.ZModOrder().RandomElement()
	commToAandB := commQ.Commit(commToAandBRand, voter1.A, voter1.B)

	ps := NewDoubleDiscreteLogProofSystem(commP, commQ, securityParam)
	proof := ps.Generate(voter1, commToU, commToURand, commToAandB,
		commToAandBRand, "yes")
	if !ps.Verify(proof, commToU, commToAandB, "yes") {
		t.Fail()
	}
}
