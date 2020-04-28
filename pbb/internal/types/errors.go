package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	BulletinBoardCodespace sdk.CodespaceType = BulletinBoardModuleName

	InvalidBallot     sdk.CodeType = 101
	InvalidCredential sdk.CodeType = 201
)

func ErrInvalidBallot(msg string) sdk.Error {
	return sdk.NewError(BulletinBoardCodespace, InvalidBallot, msg)
}
