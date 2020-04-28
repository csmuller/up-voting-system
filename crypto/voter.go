package crypto

import "math/big"

// Voter encapsulates the credentials of a voter. I.e. the private credentials alpha and beta and
// the public credential u. The election credential u_hat is excluded.
type Voter struct {
	A *big.Int // private credential alpha, random element in Z_q
	B *big.Int // private credential beta, random element in Z_q
	U *big.Int // public credential u = h1^a h2^b mod q, h1 and h2 generators of G_q
}

// GenerateNewVoter generates new private and public credentials and returns a Voter instance
// with those credentials.
func GenerateNewVoter(commQ PedersenCommitmentScheme) Voter {
	a := commQ.G.ZModOrder().RandomElement()
	b := commQ.G.ZModOrder().RandomElement()
	h1 := commQ.G.Exp(commQ.Hm[0], a)
	h2 := commQ.G.Exp(commQ.Hm[1], b)
	u := commQ.G.Mul(h1, h2)
	return NewVoter(a, b, u)
}

// NewVoter instantiates a new voter from the given credentials
func NewVoter(a, b, u *big.Int) Voter {
	return Voter{
		A: a,
		B: b,
		U: u,
	}
}
