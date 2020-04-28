package pbb

import (
	"encoding/json"
	"github.com/csmuller/up-voting-system/pbb/internal/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/csmuller/up-voting-system/pbb/client/cli"
	"github.com/csmuller/up-voting-system/pbb/client/rest"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	// Make sure the interfaces are implemented.
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string {
	return BulletinBoardModuleName
}

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis validates the Genesis block information given in the JSON message.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	err := ModuleCdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return err
	}
	return ValidateGenesis(data)
}

// Register rest routes
func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	rest.RegisterRoutes(ctx, rtr, BulletinBoardModuleName)
}

// Get the root query command of this module
func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetQueryCmd(BulletinBoardModuleName, cdc)
}

// Get the root tx command of this module
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(BulletinBoardModuleName, cdc)
}

type AppModule struct {
	AppModuleBasic
	bulletinBoardKeeper BulletinBoardKeeper
}

func NewAppModule(k BulletinBoardKeeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		bulletinBoardKeeper: k,
	}
}

func (AppModule) Name() string {
	return BulletinBoardModuleName
}

func (am AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {}

func (am AppModule) Route() string {
	return BulletinBoardModuleName
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.bulletinBoardKeeper)
}
func (am AppModule) QuerierRoute() string {
	return BulletinBoardModuleName
}

func (am AppModule) NewQuerierHandler() sdk.Querier {
	return NewQuerier(am.bulletinBoardKeeper)
}

func (am AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

func (am AppModule) EndBlock(sdk.Context, abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState types.GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.bulletinBoardKeeper, genesisState)
	return []abci.ValidatorUpdate{}
}

func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, am.bulletinBoardKeeper)
	return ModuleCdc.MustMarshalJSON(gs)
}
