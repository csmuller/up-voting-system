package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/csmuller/up-voting-system/crypto"
	"math/big"
)

//--------------------------------------------------------------------------------------------------
// MspPutPublicCredential

var _ sdk.Msg = MsgPutVoterCredential{}

// MsgPutBallot defines the message for posting a ballot to the bulletin board.
type MsgPutVoterCredential struct {
	Credential crypto.Int     `json:"u"`
	Signer     sdk.AccAddress `json:"signer"`
}

// NewMsgPutBallot creates a new instance of the MsgPutBallot message.
func NewMsgPutVoterCredential(credential crypto.Int, signer sdk.AccAddress) MsgPutVoterCredential {
	return MsgPutVoterCredential{
		Credential: credential,
		Signer:     signer,
	}
}

// Route returns the name of the module.
func (msg MsgPutVoterCredential) Route() string {
	return BulletinBoardModuleName
}

// Type returns the action of the message.
func (msg MsgPutVoterCredential) Type() string {
	return "put_voter_credential"
}

// ValidateBasic runs stateless checks on the message
func (msg MsgPutVoterCredential) ValidateBasic() sdk.Error {
	if msg.Credential.IsZero() || msg.Credential.IsNegative() {
		return sdk.NewError(BulletinBoardCodespace, InvalidCredential,
			"voter credential value cannot be zero or negative")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgPutVoterCredential) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(msg)
}

// GetSigners defines whose signature is required
func (msg MsgPutVoterCredential) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

//--------------------------------------------------------------------------------------------------
// MspPutBallot

var _ sdk.Msg = MsgPutBallot{}

// MsgPutBallot defines the message for posting a ballot to the bulletin board.
type MsgPutBallot struct {
	Ballot Ballot
	Signer sdk.AccAddress `json:"signer"`
}

func NewMsgPutBallot(c *big.Int, d *big.Int, v string, uHat *big.Int, proof1 crypto.PolyEvalProof,
	proof2 crypto.DdLogProof, proof3 crypto.PreimageEqualityProof, signer sdk.AccAddress) MsgPutBallot {

	return MsgPutBallot{
		Ballot: NewBallot(c, d, v, uHat, proof1, proof2, proof3),
		Signer: signer,
	}
}

// Route returns the name of the module.
func (msg MsgPutBallot) Route() string {
	return BulletinBoardModuleName
}

// Type returns the action of the message.
func (msg MsgPutBallot) Type() string {
	return "put_ballot"
}

// ValidateBasic runs stateless checks on the message
func (msg MsgPutBallot) ValidateBasic() sdk.Error {
	if len(msg.Ballot.V) == 0 {
		return sdk.NewError(BulletinBoardCodespace, InvalidBallot, "vote cannot be empty")
	}
	if msg.Ballot.UHat.IsZero() || msg.Ballot.UHat.IsNegative() {
		return sdk.NewError(BulletinBoardCodespace, InvalidBallot,
			"voter's election credential cannot be zero or negative")
	}
	if msg.Ballot.C.IsZero() || msg.Ballot.C.IsNegative() {
		return sdk.NewError(BulletinBoardCodespace, InvalidBallot,
			"commitment to voter public credential cannot be zero or negative")
	}
	if msg.Ballot.D.IsZero() || msg.Ballot.D.IsNegative() {
		return sdk.NewError(BulletinBoardCodespace, InvalidBallot,
			"commitment to voter private credentials cannot be zero or negative")
	}
	// TODO: Do empty array checks on the proofs.
	return nil
}

// GetSignBytes encodes the message for signing
func (msg MsgPutBallot) GetSignBytes() []byte {
	return ModuleCdc.MustMarshalJSON(msg)
}

// GetSigners defines whose signature is required
func (msg MsgPutBallot) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
