package crypto

import (
	"encoding/json"
	"fmt"
	"github.com/tendermint/go-amino"
	"math/big"
)

// Polynomial represents a polynomial in the polynomial ring Z_p[x], where Z_p is the ring of
// integers mod prime p. The coefficients are elements of Z_p and coefficient arithmetic is
// performed/modulo p.
//
// A polynomial holds a slice of big integer pointers (the coefficients) and an instance of
// ZModPrime instance, therefor copying a polynomial is not expensive.
//
// Adapted from https://github.com/mbottini/newton/blob/master/polynomial/polynomial.go.
type Polynomial struct {
	Coeffs []*big.Int // Coefficients ordered ascending by degree (i.e idx 0 is the constant term)
	ZModPr ZModPrime  // The polynomial's ring.
}

// NewPolynomial creates a new Polynomial instance with the given coefficients and the ring.
// Note that no deep copy of the coefficients is performed.
// If you don't want modificaitons of the poynomial to show in the original argument, make sure to
// create copies of the integers first.
func NewPolynomial(coeffs []*big.Int, zModPr ZModPrime) Polynomial {
	return Polynomial{
		Coeffs: coeffs,
		ZModPr: zModPr,
	}
}

// Trim removes trailing zero terms from the Polynomial and returns a new trimmed Polynomial.
// For example, 0x^3 + 0x^2 + 3x + 1 results in 3x + 1.
// The remaining coefficients are not copied, i.e. changes in the original polynomial will affect
// the trimmed polynomial and vice versa.
func (p Polynomial) Trim() Polynomial {
	return NewPolynomial(trim(p.Coeffs), p.ZModPr)
}

func trim(coeffs []*big.Int) []*big.Int {
	zero := big.NewInt(0)
	for i := len(coeffs) - 1; i >= 0; i-- {
		if coeffs[i].Cmp(zero) != 0 {
			return coeffs[0 : i+1]
		}
	}
	// We keep a 0 constant term if the polynomial really is 0.
	return coeffs[0:1]
}

// Degree gets the degree of this polynomial, i.e. the largest integer m for which coefficient
// a_m != 0
func (p Polynomial) Degree() int {
	return len(trim(p.Coeffs)) - 1
}

// String returns a string representation of this polynomial.
func (p Polynomial) String() string {
	polyMax := len(p.Coeffs) - 1
	resultStr := ""
	for i := range p.Coeffs {
		power := polyMax - i
		if power != 0 {
			resultStr += fmt.Sprintf("%dx^%d + ", p.Coeffs[power], power)
		} else {
			resultStr += fmt.Sprintf("%d", p.Coeffs[power])
		}
	}
	return resultStr
}

// Add adds the given polynomial to this one and returns the result. Both inputs are not modified
// and the new polynomial's coefficients are new Int instances, i.e. changes in the input
// polynomials do not affect the new polynomial and vice versa.
func (p Polynomial) Add(other Polynomial) Polynomial {
	if p.ZModPr.Modulus != other.ZModPr.Modulus {
		panic("Cannot add two polynomials in different rings.")
	}
	maxIndex := len(p.Coeffs)
	if maxIndex < len(other.Coeffs) {
		maxIndex = len(other.Coeffs)
	}
	newCoeffs := make([]*big.Int, maxIndex)
	var currCoeff *big.Int
	for i := 0; i < maxIndex; i++ {
		currCoeff = big.NewInt(0)
		if i < len(p.Coeffs) {
			currCoeff = p.ZModPr.Add(currCoeff, p.Coeffs[i])
		}
		if i < len(other.Coeffs) {
			currCoeff = p.ZModPr.Add(currCoeff, other.Coeffs[i])
		}
		newCoeffs[i] = currCoeff
	}
	return NewPolynomial(newCoeffs, p.ZModPr)
}

// MulScalar multiplies each coefficient in this polynomial by scalar and returns the resulting
// polynomial. The input polynomial is not modified and the new polynomial's coefficients are new
// Int instances, i.e. changes in the input polynomial do not affect the new polynomial and vice
// versa.
func (p Polynomial) MulScalar(scalar *big.Int) Polynomial {
	newCoeffs := make([]*big.Int, len(p.Coeffs))
	for i, coefficient := range p.Coeffs {
		newCoeffs[i] = p.ZModPr.Mul(coefficient, scalar)
	}
	return NewPolynomial(trim(newCoeffs), p.ZModPr)
}

// Mul multiplies two Polynomials together and returns the resulting Polynomial.
// The input polynomials are not modified and the new polynomial's coefficients are new
// Int instances, i.e. changes in the input polynomials do not affect the new polynomial and vice
// versa.
func (p Polynomial) Mul(other Polynomial) Polynomial {
	if p.ZModPr.Modulus != other.ZModPr.Modulus {
		panic("Cannot multiply two polynomials in different rings.")
	}
	newDegree := p.Degree() + other.Degree()
	newCoeffs := make([]*big.Int, newDegree+1) // +1 for the constant term
	for i, thisCoeff := range p.Coeffs {
		for j, otherCoeff := range other.Coeffs {
			tmp := p.ZModPr.Mul(thisCoeff, otherCoeff)
			if newCoeffs[i+j] == nil {
				newCoeffs[i+j] = big.NewInt(0)
			}
			newCoeffs[i+j] = p.ZModPr.Add(newCoeffs[i+j], tmp)
		}
	}
	return NewPolynomial(newCoeffs, p.ZModPr)
}

// IncludeCredential includes the given credential u into this polynomial and returns the resulting
// polynomial. The credential is included by multiplying the polynomial with (X - u). The input
// polynomial is not modified and the new polynomial's coefficients are new Int instances, i.e.
// changes in the input polynomial do not affect the new polynomial and vice versa.
func (p Polynomial) IncludeCredential(u *big.Int) Polynomial {
	coeffs := []*big.Int{p.ZModPr.AdditiveInvert(u), big.NewInt(1)}
	p2 := NewPolynomial(coeffs, p.ZModPr)
	return p.Mul(p2)
}

// PolynomialDTO wraps the coefficients of a Polynomial into the amino serializable Int class.
type polynomialDTO struct {
	Coeffs    []Int     `json:"coeffs"`
	ZModPrDTO ZModPrime `json:"zmod"`
}

func (p Polynomial) MarshalAmino() (string, error) {
	dto := polynomialDTO{}
	p.wrapInDTO(&dto)
	return string(amino.MustMarshalBinaryBare(dto)), nil
}

func (p *Polynomial) UnmarshalAmino(bytes []byte) error {
	var dto polynomialDTO
	amino.MustUnmarshalBinaryBare(bytes, &dto)
	p.unwrapDTO(dto)
	return nil
}

func (p Polynomial) MarshalJSON() ([]byte, error) {
	dto := polynomialDTO{}
	p.wrapInDTO(&dto)
	return json.Marshal(dto)
}

func (p *Polynomial) UnmarshalJSON(bytes []byte) error {
	var dto polynomialDTO
	err := amino.UnmarshalJSON(bytes, &dto)
	if err != nil {
		return err
	}
	p.unwrapDTO(dto)
	return nil
}

func (p Polynomial) wrapInDTO(dto *polynomialDTO) {
	dto.Coeffs = make([]Int, len(p.Coeffs))
	for i, coeff := range p.Coeffs {
		dto.Coeffs[i] = NewInt(coeff)
	}
	dto.ZModPrDTO = p.ZModPr
}

func (p *Polynomial) unwrapDTO(dto polynomialDTO) {
	p.Coeffs = make([]*big.Int, len(dto.Coeffs))
	for i, coeff := range dto.Coeffs {
		p.Coeffs[i] = coeff.BigInt()
	}
	p.ZModPr = dto.ZModPrDTO
}
