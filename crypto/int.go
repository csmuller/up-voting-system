package crypto

import (
	"encoding/json"
	"fmt"
	"math/big"
)

// Int wraps a Go big integer and is necessary for marshalling/unmarshalling big
// integers in Tendermint. We explicitly don't use github.com/cosmos/cosmos-sdk/types.Int because it
// is restricted to 256-bit integers. The UEP voting protocol requires bigger integers for
// sufficient security.
type Int struct {
	i *big.Int
}

// NewInt creates a new Int instance wrapping the given Go big integer.
func NewInt(i *big.Int) Int {
	return Int{i}
}

// BigInt returns the big integer wrapped in this Int.
func (i Int) BigInt() *big.Int {
	return i.i
}

// IntFromString tries to parses the given string as a big integer and returns a new instance of
// Int.
func IntFromString(s string) (Int, error) {
	i, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return Int{}, fmt.Errorf("failed parsing integer %s", s)
	}
	return Int{i}, nil
}

// MarshalAmino byte-serializes this big integer and returns the bytes as a string.
func (i Int) MarshalAmino() (string, error) {
	if i.i == nil {
		i.i = new(big.Int)
	}
	text, err := i.i.MarshalText()
	return string(text), err
}

// UnmarshalAmino deserializes the given string of bytes into this big integer.
func (i *Int) UnmarshalAmino(text string) error {
	if i.i == nil {
		i.i = new(big.Int)
	}
	if err := i.i.UnmarshalText([]byte(text)); err != nil {
		return err
	}
	return nil
}

// MarshalJSON serializes this big integer to JSON format
func (i Int) MarshalJSON() ([]byte, error) {
	if i.i == nil {
		i.i = new(big.Int)
	}
	text, err := i.i.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(text))
}

// UnmarshalJSON deserializes the given JSON bytes into this big integer.
func (i *Int) UnmarshalJSON(bytes []byte) error {
	if i.i == nil {
		i.i = new(big.Int)
	}
	var text string
	if err := json.Unmarshal(bytes, &text); err != nil {
		return err
	}

	if err := i.i.UnmarshalText([]byte(text)); err != nil {
		return err
	}
	return nil
}

// IsZero returns true if this integer is zero.
func (i Int) IsZero() bool {
	return i.i.Sign() == 0
}

// IsNegative returns true if this integer has a negative sign.
func (i Int) IsNegative() bool {
	return i.i.Sign() == -1
}

// String returns a string representation of this Int.
func (i Int) String() string {
	return i.i.String()
}
