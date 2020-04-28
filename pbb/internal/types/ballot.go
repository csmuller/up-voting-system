package types

import (
	"fmt"
	"github.com/csmuller/uep-voting/crypto"
	"math/big"
	"strings"
)

type Ballot struct {
	C      crypto.Int                   `json:"c"`     // commitment to u, value in Gp
	D      crypto.Int                   `json:"d"`     // commitment to a and b, value in Gq
	V      string                       `json:"v"`     // vote
	UHat   crypto.Int                   `json:"u_hat"` // election credential
	Proof1 crypto.PolyEvalProof         `json:"p1"`
	Proof2 crypto.DdLogProof            `json:"p2"`
	Proof3 crypto.PreimageEqualityProof `json:"p3"`
}

func NewBallot(c *big.Int, d *big.Int, v string, uHat *big.Int,
	proof1 crypto.PolyEvalProof, proof2 crypto.DdLogProof, proof3 crypto.PreimageEqualityProof) Ballot {

	return Ballot{
		C:      crypto.NewInt(c),
		D:      crypto.NewInt(d),
		V:      v,
		UHat:   crypto.NewInt(uHat),
		Proof1: proof1,
		Proof2: proof2,
		Proof3: proof3,
	}
}

func (b Ballot) String() string {
	var str strings.Builder
	str.WriteString("Ballot: {\n")
	str.WriteString(fmt.Sprintf("\tc: %s\n", b.C.String()))
	str.WriteString(fmt.Sprintf("\td: %s\n", b.D.String()))
	str.WriteString(fmt.Sprintf("\tv: %s\n", b.V))
	str.WriteString(fmt.Sprintf("\tu_hat: %s\n", b.UHat))
	str.WriteString(fmt.Sprintf("\tproof1: %s\n", b.Proof1.String()))
	str.WriteString(fmt.Sprintf("\tproof2: %s\n", b.Proof2.String()))
	str.WriteString(fmt.Sprintf("\tproof3: %s\n", b.Proof3.String()))
	str.WriteString("}")
	return str.String()
}
