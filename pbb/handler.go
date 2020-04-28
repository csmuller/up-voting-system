package pbb

import "C"
import (
	"fmt"
	"github.com/csmuller/uep-voting/crypto"
	"github.com/csmuller/uep-voting/pbb/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	EventTypeBallot          = "ballot"
	EventTypeVoterCredential = "voterCredential"

	AttributeKeyElectionCredential = "electionCredential"
	AttributeKeyVote               = "vote"
	AttributeKeyVoterCredential    = "voterCredential"
)

// NewHandler returns a handler for bulletin board messages
func NewHandler(keeper BulletinBoardKeeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgPutBallot:
			return handleMsgPutBallot(ctx, keeper, msg)
		case MsgPutVoterCredential:
			return handleMsgPutVoterCredential(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized bulletin board message type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgPutBallot(ctx sdk.Context, keeper BulletinBoardKeeper, msg MsgPutBallot) sdk.Result {
	// TODO: Discard the ballot if the contained vote does not adhere to a specified format.
	if keeper.HasElectionCredential(ctx, msg.Ballot.UHat.BigInt()) {
		return types.ErrInvalidBallot("A ballot has already been stored for this election " +
			"credential.").Result()
	}
	params := keeper.GetParams(ctx)
	if !isMembershipProofValid(msg.Ballot, params.CommP, keeper.GetCredentialPolynomial(ctx)) {
		return types.ErrInvalidBallot("Invalid membership proof").Result()
	}
	if !isRepresentationProofValid(msg.Ballot, params) {
		return types.ErrInvalidBallot("Invalid proof of known representation").Result()
	}
	if !isPreimageProofValid(msg.Ballot, params) {
		return types.ErrInvalidBallot("Invalid pre-image proof").Result()
	}
	if err := keeper.StoreBallot(ctx, msg.Ballot); err != nil {
		return types.ErrInvalidBallot(err.Error()).Result()
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(EventTypeBallot,
		sdk.NewAttribute(AttributeKeyElectionCredential, msg.Ballot.UHat.String()),
		sdk.NewAttribute(AttributeKeyVote, msg.Ballot.V)))
	return sdk.Result{Code: sdk.CodeOK}
}

func isMembershipProofValid(b Ballot, commP crypto.PedersenCommitmentScheme,
	poly crypto.Polynomial) bool {
	ps := crypto.NewPolynomialEvaluationProofSystem(commP, poly)
	return ps.Verify(b.Proof1, b.C.BigInt(), b.V)
}

func isRepresentationProofValid(b Ballot, params Params) bool {
	ps := crypto.NewDoubleDiscreteLogProofSystem(params.CommP, params.CommQ, params.SecurityParam)
	return ps.Verify(b.Proof2, b.C.BigInt(), b.D.BigInt(), b.V)
}

func isPreimageProofValid(b Ballot, params Params) bool {
	ps := crypto.NewPreimageEqualityProofSystem(params.HHat.BigInt(), params.CommQ)
	return ps.Verify(b.Proof3, b.D.BigInt(), b.UHat.BigInt(), b.V)
}

func handleMsgPutVoterCredential(ctx sdk.Context, keeper BulletinBoardKeeper,
	msg types.MsgPutVoterCredential) sdk.Result {

	// TODO: When an identity management system is available check if the sender account belongs
	//  to an eligible voter and store the voters identity with the public credential.
	err := keeper.StoreVoterCredential(ctx, msg.Credential)
	if err != nil {
		return sdk.NewError(types.BulletinBoardCodespace, sdk.CodeInternal, "%v", err).Result()
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(EventTypeVoterCredential,
		sdk.NewAttribute(AttributeKeyVoterCredential, msg.Credential.String()),
		sdk.NewAttribute(sdk.AttributeKeySender, msg.Signer.String())))
	return sdk.Result{Code: sdk.CodeOK}
}
