package types

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/csmuller/up-voting-system/crypto"
	"math/big"
	"strings"
)

const DefaultParamSpace = BulletinBoardModuleName

var (
	// Parameter keys
	CommPKey         = []byte("CommP")
	CommQKey         = []byte("CommQ")
	HKey             = []byte("HHat")
	SecurityParamKey = []byte("SecurityParam")
)

// Params implements the ParamSet interface
type Params struct {
	CommP         crypto.PedersenCommitmentScheme `json:"comm_p"`
	CommQ         crypto.PedersenCommitmentScheme `json:"comm_q"`
	HHat          crypto.Int                      `json:"h"` // election generator
	SecurityParam int                             `json:"k"`
}

// ParamSetPairs returns all the key/value pairs pairs of the bulletin board module's parameters.
func (p *Params) ParamSetPairs() subspace.ParamSetPairs {
	return subspace.ParamSetPairs{
		{CommPKey, &p.CommP},
		{CommQKey, &p.CommQ},
		{HKey, &p.HHat},
		{SecurityParamKey, &p.SecurityParam},
	}
}

func (p Params) String() string {
	var str strings.Builder
	str.WriteString("Parameters: {\n")
	str.WriteString(fmt.Sprintf("commP: %s,\n", p.CommP.String()))
	str.WriteString(fmt.Sprintf("commQ: %s,\n", p.CommQ.String()))
	str.WriteString(fmt.Sprintf("HHat: %s,\n", p.HHat.String()))
	str.WriteString("}")
	return str.String()
}

func ParamKeyTable() subspace.KeyTable {
	return subspace.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object which does not set the election generator by choice.
func NewParams(commP crypto.PedersenCommitmentScheme, commQ crypto.PedersenCommitmentScheme,
	h *big.Int, securityParam int) Params {

	return Params{
		CommP:         commP,
		CommQ:         commQ,
		HHat:          crypto.NewInt(h),
		SecurityParam: securityParam,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	o, _ := new(big.Int).SetString(crypto.O, 10)
	p, _ := new(big.Int).SetString(crypto.P, 10)
	q, _ := new(big.Int).SetString(crypto.Q, 10)

	gP := crypto.NewGStarModPrime(o, p)
	commP := crypto.NewPedersenCommitmentScheme(gP, gP.RandomGenerator(),
		[]*big.Int{gP.RandomGenerator()})

	gQ := crypto.NewGStarModPrime(p, q)
	commQ := crypto.NewPedersenCommitmentScheme(gQ, gQ.RandomGenerator(),
		[]*big.Int{gQ.RandomGenerator(), gQ.RandomGenerator()}) // comm_q can take two messages

	h := gQ.RandomElement()

	return NewParams(commP, commQ, h, crypto.SecurityParam)
}
