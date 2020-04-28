package pbb

import (
	"github.com/csmuller/up-voting-system/pbb/internal/keeper"
	"github.com/csmuller/up-voting-system/pbb/internal/types"
)

const (
	BulletinBoardModuleName = types.BulletinBoardModuleName
	BallotStoreKey          = types.BallotStoreKey
	VoterCredentialStoreKey = types.VoterCredentialStoreKey
	PolynomialStoreKey      = types.PolynomialStoreKey
	DefaultParamSpace       = types.DefaultParamSpace
)

var (
	NewBulletinBoardKeeper = keeper.NewBulletinBoardKeeper
	NewQuerier             = keeper.NewQuerier
	ModuleCdc              = types.ModuleCdc
	RegisterCodec          = types.RegisterCodec
)

type (
	BulletinBoardKeeper      = keeper.BulletinBoardKeeper
	Ballot                   = types.Ballot
	MsgPutBallot             = types.MsgPutBallot
	MsgPutVoterCredential    = types.MsgPutVoterCredential
	QueryResVoterCredentials = types.QueryResVoterCredentials
	Params                   = types.Params
)
