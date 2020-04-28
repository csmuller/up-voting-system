package main

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/x/genaccounts"
	genaccscli "github.com/cosmos/cosmos-sdk/x/genaccounts/client/cli"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/csmuller/up-voting-system/crypto"
	"github.com/csmuller/up-voting-system/pbb"
	tlog "github.com/tendermint/tendermint/libs/log"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	sdk "github.com/cosmos/cosmos-sdk/types"
	genutilcli "github.com/cosmos/cosmos-sdk/x/genutil/client/cli"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	// NodeHomeDirectory is the default directory where the PBB data and configuration will be
	// stored.
	NodeHomeDirectory = os.ExpandEnv("$HOME/.pbbd")
	// DefaultClientHomeDirectory is the default directory for pbbd commands that require access
	// to client data, e.g. accounts needed for signing genesis transactions.
	DefaultClientHomeDirectory = os.ExpandEnv("$HOME/.acli")
)

func main() {
	crypto.SetupLogger("pbb-log.txt")
	cobra.EnableCommandSorting = false
	cdc := pbb.MakeCodec()

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(sdk.Bech32PrefixAccAddr, sdk.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(sdk.Bech32PrefixValAddr, sdk.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(sdk.Bech32PrefixConsAddr, sdk.Bech32PrefixConsPub)
	config.Seal()

	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "pbbd",
		Short:             "PBB App Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}
	// CLI commands to initialize the chain
	rootCmd.AddCommand(
		genutilcli.InitCmd(ctx, cdc, pbb.ModuleManager, NodeHomeDirectory),
		genutilcli.CollectGenTxsCmd(ctx, cdc, genaccounts.AppModuleBasic{}, NodeHomeDirectory),
		genutilcli.GenTxCmd(
			ctx, cdc, pbb.ModuleManager, staking.AppModuleBasic{},
			genaccounts.AppModuleBasic{}, NodeHomeDirectory, DefaultClientHomeDirectory,
		),
		genutilcli.ValidateGenesisCmd(ctx, cdc, pbb.ModuleManager),
		// AddGenesisAccountCmd allows users to add accounts to the genesis file
		genaccscli.AddGenesisAccountCmd(ctx, cdc, NodeHomeDirectory, DefaultClientHomeDirectory),
	)

	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "BulletinBoard", NodeHomeDirectory)
	err := executor.Execute()
	if err != nil {
		panic(err)
	}
}

func newApp(logger tlog.Logger, db dbm.DB, traceStore io.Writer) abci.Application {
	return pbb.NewBulletinBoardApp(logger, db) //baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))),
	//baseapp.SetMinGasPrices(viper.GetString(server.FlagMinGasPrices)),
	//baseapp.SetHaltHeight(viper.GetUint64(server.FlagHaltHeight)),
	//baseapp.SetHaltTime(viper.GetUint64(server.FlagHaltTime)),

}

func exportAppStateAndTMValidators(
	logger tlog.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool,
	jailWhiteList []string,
) (json.RawMessage, []tmtypes.GenesisValidator, error) {

	if height != -1 {
		nsApp := pbb.NewBulletinBoardApp(logger, db)
		err := nsApp.LoadHeight(height)
		if err != nil {
			return nil, nil, err
		}
		return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
	}

	nsApp := pbb.NewBulletinBoardApp(logger, db)

	return nsApp.ExportAppStateAndValidators(forZeroHeight, jailWhiteList)
}
