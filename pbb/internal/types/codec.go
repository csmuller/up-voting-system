package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/csmuller/up-voting-system/crypto"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Ballot{}, "pbb/Ballot", nil)
	cdc.RegisterConcrete(MsgPutBallot{}, "pbb/PutBallot", nil)
	cdc.RegisterConcrete(MsgPutVoterCredential{}, "pbb/PutVoterCredential", nil)
	cdc.RegisterConcrete(crypto.Polynomial{}, "pbb/Polynomial", nil)
	cdc.RegisterConcrete(crypto.GStarModPrime{}, "pbb/GStarModPrime", nil)
	cdc.RegisterConcrete(crypto.ZModPrime{}, "pbb/ZModPrime", nil)
	cdc.RegisterConcrete(crypto.ZStarModPrime{}, "pbb/ZStarModPrime", nil)
	cdc.RegisterConcrete(crypto.PedersenCommitmentScheme{}, "pbb/PedersenCommScheme", nil)
	cdc.RegisterConcrete(crypto.Int{}, "pbb/Int", nil)
	cdc.RegisterConcrete(crypto.DdLogProof{}, "pbb/DdLogProof", nil)
	cdc.RegisterConcrete(crypto.PolyEvalProof{}, "pbb/PolyEvalProof", nil)
	cdc.RegisterConcrete(crypto.PreimageEqualityProof{}, "pbb/PreimageEqualityProof", nil)
}
