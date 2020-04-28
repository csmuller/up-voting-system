package keeper

import (
	"encoding/binary"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/csmuller/uep-voting/pbb/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// Query endpoints supported by the pbb Querier
const (
	QueryBallots              = "ballots"
	QueryParameters           = "parameters"
	QueryVoterCredentials     = "voterCredentials"
	QueryCredentialPolynomial = "credentialPolynomial"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper BulletinBoardKeeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryBallots:
			return queryBallots(ctx, keeper)
		case QueryVoterCredentials:
			return queryVoterCredentials(ctx, keeper)
		case QueryParameters:
			return queryParameters(ctx, keeper)
		case QueryCredentialPolynomial:
			return queryCredentialPolynomial(ctx, keeper)
		default:
			return nil, sdk.ErrUnknownRequest(
				fmt.Sprintf("Unknown bulletin board query endpoint %s.", path[0]))
		}
	}
}

func queryBallots(ctx sdk.Context, keeper BulletinBoardKeeper) ([]byte, sdk.Error) {
	var ballots []types.Ballot

	it := keeper.GetBallotsIterator(ctx)

	for ; it.Valid(); it.Next() {
		var ballot types.Ballot
		keeper.cdc.MustUnmarshalBinaryBare(it.Value(), &ballot)
		ballots = append(ballots, ballot)
	}

	res, err := keeper.cdc.MarshalJSONIndent(ballots, "", "  ")
	if err != nil {
		panic("Could not marshal list of ballots to JSON.")
	}
	return res, nil
}

func queryVoterCredentials(ctx sdk.Context, keeper BulletinBoardKeeper) ([]byte, sdk.Error) {
	var results types.QueryResVoterCredentials

	it := keeper.GetVoterCredentialsIterator(ctx)

	for ; it.Valid(); it.Next() {
		var result types.QueryResVoterCredential
		keeper.cdc.MustUnmarshalBinaryBare(it.Key(), &result.Credential)
		result.BlockHeight, _ = binary.Varint(it.Value())
		results = append(results, result)
	}

	res, err := keeper.cdc.MarshalJSONIndent(results, "", "  ")
	if err != nil {
		panic("Could not marshal credentials and their respective publish timestamps to JSON.")
	}
	return res, nil
}

func queryParameters(ctx sdk.Context, keeper BulletinBoardKeeper) ([]byte, sdk.Error) {
	params := keeper.GetParams(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal parameters to JSON",
			err.Error()))
	}
	return res, nil
}

func queryCredentialPolynomial(ctx sdk.Context, keeper BulletinBoardKeeper) ([]byte, sdk.Error) {
	poly := keeper.GetCredentialPolynomial(ctx)
	res, err := codec.MarshalJSONIndent(types.ModuleCdc, poly)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal polynomial to JSON",
			err.Error()))
	}
	return res, nil
}
