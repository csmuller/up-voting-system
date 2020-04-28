package crypto

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/go-amino"
	"math/big"
	"strings"
)

// PedersenCommitmentScheme is used to commit to a message m given a random value r in the
// following way: Com(m, r) = h1^m * h2^r. Where h1 and h2 are generators of a cyclic group.
// The scheme can be used in a generalized manner, i.e. one can commit to multiple messages in one
// commitment. Thus, a instance of the scheme can have multiple message generators.
type PedersenCommitmentScheme struct {
	G  GStarModPrime // cyclic, multiplicative group cyclic group of the generators
	Hr *big.Int      // randomization generator, h_0
	Hm []*big.Int    // message generators, h_1 to h_n
}

// NewPedersenCommitmentScheme creates a new instance of the scheme. The instance is based on the
// given group and generators which must be generators of the group.
func NewPedersenCommitmentScheme(g GStarModPrime, hr *big.Int,
	hm []*big.Int) PedersenCommitmentScheme {

	return PedersenCommitmentScheme{
		G:  g,
		Hr: hr,
		Hm: hm,
	}
}

// Commit creates a commitment to the given messages msgs with the given randomness r.
func (s *PedersenCommitmentScheme) Commit(r *big.Int, msgs ...*big.Int) *big.Int {
	if len(msgs) != len(s.Hm) {
		panic("The number of messages is not equal to the number of message generators in this" +
			" commitment scheme.")
	}
	for _, msg := range msgs {
		if msg.Cmp(s.G.Order) == 1 || msg.Cmp(big.NewInt(0)) == -1 {
			panic(fmt.Sprintf("The message is not in Z_q, where q is %s", s.G.Order.String()))
		}
	}
	if r.Cmp(s.G.Order) == 1 || r.Cmp(big.NewInt(0)) == -1 {
		panic(fmt.Sprintf("The random value is not in Z_q, where q is %s", s.G.Order.String()))
	}

	product := s.G.Exp(s.Hr, r)
	for i, msg := range msgs {
		if msg.Cmp(s.G.Order) == 1 || msg.Cmp(big.NewInt(0)) == -1 {
			panic(fmt.Sprintf("The message is not in Z_q, where q is %s", s.G.Order.String()))
		}
		t := s.G.Exp(s.Hm[i], msg)
		product = s.G.Mul(product, t)
	}
	return product
}

// String returns a string representation of this scheme.
func (s *PedersenCommitmentScheme) String() string {
	var sb strings.Builder
	sb.WriteString("PedComSch: {\n")
	sb.WriteString(fmt.Sprintf("\th_r: %s\n", s.Hr))
	for i, b := range s.Hm {
		sb.WriteString(fmt.Sprintf("\th_%d: %s\n", i+1, b))
	}
	return sb.String()
}

// pedersenCommitmentSchemeDTO is used for Tendermint serialization and deserialization of the
// PedersenCommitmentScheme type.
type pedersenCommitmentSchemeDTO struct {
	GStarModPr GStarModPrime `json:"g"`
	Hr         Int           `json:"hr"`
	Hm         []Int         `json:"hm"`
}

func (s PedersenCommitmentScheme) MarshalAmino() (string, error) {
	hm := make([]Int, len(s.Hm))
	for i, hi := range s.Hm {
		hm[i] = NewInt(hi)
	}
	dto := pedersenCommitmentSchemeDTO{GStarModPr: s.G, Hr: NewInt(s.Hr), Hm: hm}
	return string(amino.MustMarshalBinaryBare(dto)), nil
}

func (s *PedersenCommitmentScheme) UnmarshalAmino(bytes []byte) error {
	var dto pedersenCommitmentSchemeDTO
	amino.MustUnmarshalBinaryBare(bytes, &dto)
	s.Hm = make([]*big.Int, len(dto.Hm))
	for i, hi := range dto.Hm {
		s.Hm[i] = hi.BigInt()
	}
	s.Hr = dto.Hr.BigInt()
	s.G = dto.GStarModPr
	return nil
}

// MarshalJSON serializes this big integer to JSON format
func (s PedersenCommitmentScheme) MarshalJSON() ([]byte, error) {
	hm := make([]Int, len(s.Hm))
	for i, hi := range s.Hm {
		hm[i] = NewInt(hi)
	}
	dto := pedersenCommitmentSchemeDTO{GStarModPr: s.G, Hr: NewInt(s.Hr), Hm: hm}
	return json.Marshal(dto)
}

func (s *PedersenCommitmentScheme) UnmarshalJSON(bytes []byte) error {
	var dto pedersenCommitmentSchemeDTO
	err := amino.UnmarshalJSON(bytes, &dto)
	if err != nil {
		return err
	}
	s.Hm = make([]*big.Int, len(dto.Hm))
	for i, hi := range dto.Hm {
		s.Hm[i] = hi.BigInt()
	}
	s.Hr = dto.Hr.BigInt()
	s.G = dto.GStarModPr
	return nil
}
