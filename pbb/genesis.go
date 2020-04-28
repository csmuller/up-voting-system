package pbb

import (
	"github.com/csmuller/uep-voting/pbb/internal/keeper"
	"github.com/csmuller/uep-voting/pbb/internal/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func ValidateGenesis(genesisState types.GenesisState) error {
	// TODO: Validate parameters.
	//params := genesisState.Params
	return nil
}

func DefaultGenesisState() types.GenesisState {
	return types.GenesisState{
		Params: types.DefaultParams(),
	}
}

// InitGenesis - Init parameter store from genesis data
func InitGenesis(ctx sdk.Context, bk keeper.BulletinBoardKeeper, data types.GenesisState) {
	bk.SetParams(ctx, data.Params)
}

// ExportGenesis returns a GenesisState for a given context and bulletinBoardKeeper
func ExportGenesis(ctx sdk.Context, bk keeper.BulletinBoardKeeper) types.GenesisState {
	params := bk.GetParams(ctx)
	return types.NewGenesisState(params)
}
