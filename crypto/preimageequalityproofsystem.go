package crypto

import (
	"crypto/sha256"
	"encoding/json"
	"github.com/tendermint/go-amino"
	"math/big"
	"time"
)

// PreimageEqualityProofSystem is used to proof equality of preimages. In the case of the UEP voting
// protocol it is used to proof that a voter's private credential beta was used in generating
//commitment d and election credential uHat.
type PreimageEqualityProofSystem struct {
	HHat *big.Int // Election generator used to generate the voter's election credential
	// Commitment scheme used to commit to the voter's private credentials alpha and beta.
	CommScheme PedersenCommitmentScheme
	gStarModPr GStarModPrime
	zModPr     ZModPrime
}

// NewPreimageEqualityProofSystem sets up a new instance of the proof system. Parameter hHat is the
// election generator used to generate the voter's election credential uHat. The parameter
// commScheme is the commitment scheme used to commit to the voter's private credentials alpha and
// beta.
func NewPreimageEqualityProofSystem(hHat *big.Int,
	commScheme PedersenCommitmentScheme) PreimageEqualityProofSystem {

	return PreimageEqualityProofSystem{
		HHat:       hHat,
		CommScheme: commScheme,
		gStarModPr: commScheme.G,
		zModPr:     commScheme.G.ZModOrder(),
	}
}

// PreimageEqualityProof represents a transcript of a preimage equality proof.
type PreimageEqualityProof struct {
	Comm     *big.Int
	CommHHat *big.Int
	RespA    *big.Int
	RespB    *big.Int
	RespS    *big.Int
}

// Generate generates a preimage equality proof for the given commitment to alpha and beta
// (commToAandB) and the election credential uHat.
func (ps *PreimageEqualityProofSystem) Generate(voter Voter, commToAandB *big.Int,
	commToAandBRand *big.Int, uHat *big.Int, vote string) PreimageEqualityProof {

	defer LogExecutionTime(time.Now(), "equality preimage proof generation")

	ra := ps.zModPr.RandomElement()
	rb := ps.zModPr.RandomElement()
	rs := ps.zModPr.RandomElement()

	comm1 := ps.CommScheme.Commit(rs, ra, rb)
	comm2 := ps.gStarModPr.Exp(ps.HHat, rb)

	ch := ps.generateChallenge(commToAandB, uHat, []*big.Int{comm1, comm2}, vote)

	a := voter.A
	b := voter.B
	s := commToAandBRand

	return PreimageEqualityProof{
		Comm:     comm1,
		CommHHat: comm2,
		RespA:    ps.zModPr.Add(ra, ps.zModPr.Mul(a, ch)),
		RespB:    ps.zModPr.Add(rb, ps.zModPr.Mul(b, ch)),
		RespS:    ps.zModPr.Add(rs, ps.zModPr.Mul(s, ch)),
	}
}

// Verify verifies the given preimage equality proof transcript.
func (ps *PreimageEqualityProofSystem) Verify(proof PreimageEqualityProof, commToAandB *big.Int,
	uHat *big.Int, vote string) bool {

	defer LogExecutionTime(time.Now(), "preimage equality proof verification")

	ch := ps.generateChallenge(commToAandB, uHat, []*big.Int{proof.Comm, proof.CommHHat}, vote)

	res1 := make([]*big.Int, 2)
	res1[0] = ps.CommScheme.Commit(proof.RespS, proof.RespA, proof.RespB)
	res1[1] = ps.gStarModPr.Exp(ps.HHat, proof.RespB)

	res2 := make([]*big.Int, 2)
	res2[0] = ps.gStarModPr.Mul(proof.Comm, ps.gStarModPr.Exp(commToAandB, ch))
	res2[1] = ps.gStarModPr.Mul(proof.CommHHat, ps.gStarModPr.Exp(uHat, ch))

	return res1[0].Cmp(res2[0]) == 0 && res1[1].Cmp(res2[1]) == 0
}

func (ps *PreimageEqualityProofSystem) generateChallenge(commToAandB *big.Int, uHat *big.Int,
	commitments []*big.Int, vote string) *big.Int {
	inputs := [][]*big.Int{{commToAandB, uHat}, commitments}
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

// preimageEqualityProofDTO is required for Tendermint serialization and deserialization.
type preimageEqualityProofDTO struct {
	Comm     Int `json:"comm"`
	CommHHat Int `json:"comm_h_hat"`
	RespA    Int `json:"resp_a"`
	RespB    Int `json:"resp_b"`
	RespS    Int `json:"resp_s"`
}

func (p PreimageEqualityProof) MarshalAmino() (string, error) {
	dto := preimageEqualityProofDTO{}
	p.wrapInDTO(&dto)
	return string(amino.MustMarshalBinaryBare(dto)), nil
}

func (p *PreimageEqualityProof) UnmarshalAmino(bytes []byte) error {
	var dto preimageEqualityProofDTO
	amino.MustUnmarshalBinaryBare(bytes, &dto)
	p.unwrapDTO(dto)
	return nil
}

func (p PreimageEqualityProof) MarshalJSON() ([]byte, error) {
	dto := preimageEqualityProofDTO{}
	p.wrapInDTO(&dto)
	return json.Marshal(dto)
}

func (p *PreimageEqualityProof) UnmarshalJSON(bytes []byte) error {
	var dto preimageEqualityProofDTO
	err := amino.UnmarshalJSON(bytes, &dto)
	if err != nil {
		return err
	}
	p.unwrapDTO(dto)
	return nil
}

func (p PreimageEqualityProof) wrapInDTO(dto *preimageEqualityProofDTO) {
	dto.Comm = NewInt(p.Comm)
	dto.CommHHat = NewInt(p.CommHHat)
	dto.RespA = NewInt(p.RespA)
	dto.RespB = NewInt(p.RespB)
	dto.RespS = NewInt(p.RespS)
}

func (p *PreimageEqualityProof) unwrapDTO(dto preimageEqualityProofDTO) {
	p.Comm = dto.Comm.BigInt()
	p.CommHHat = dto.CommHHat.BigInt()
	p.RespA = dto.RespA.BigInt()
	p.RespB = dto.RespB.BigInt()
	p.RespS = dto.RespS.BigInt()
}

func (p PreimageEqualityProof) String() string {
	return ""
}
