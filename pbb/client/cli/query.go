package cli

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/csmuller/uep-voting/crypto"
	"github.com/csmuller/uep-voting/pbb/internal/keeper"
	"github.com/csmuller/uep-voting/pbb/internal/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"os"
	"path"
)

const (
	votesFileName             = "votes.txt"
	defaultParamsFileName     = "params.json"
	defaultPolynomialFileName = "poly.json"
)

func GetQueryCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	bulletinBoardQueryCmd := &cobra.Command{
		Use:                        types.BulletinBoardModuleName,
		Short:                      "Querying commands for the pbb module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	bulletinBoardQueryCmd.AddCommand(client.GetCommands(
		GetCmdVerifyBallots(storeKey, cdc),
		GetCmdVoterCredentials(storeKey, cdc),
		GetCmdParameters(storeKey, cdc),
		GetCmdCredentialPolynomial(storeKey, cdc),
	)...)
	return bulletinBoardQueryCmd
}

// GetCmdVerifyBallots retrieves the list of all ballots and verifies each of them.
func GetCmdVerifyBallots(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "verify",
		Short: "Retrieve all ballots stored on the bulletin board, verify them, " +
			"and store the valid votes in a file.",
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			ballots, err := QueryBallots(cliCtx, cdc)
			if err != nil {
				return err
			}
			params, err := QueryBulletinBoardParameters(cliCtx, cdc)
			if err != nil {
				return err
			}
			if params.HHat.BigInt() == nil {
				return errors.New("election generator has not been defined yet")
			}
			poly, err := QueryCredentialPolynomial(cliCtx, cdc)
			if err != nil {
				return err
			}
			commP := params.CommP
			commQ := params.CommQ
			// TODO: Make sure that the proof systems are immutable, i.e can be used multiple times
			//  for different proofs.
			ps1 := crypto.NewPolynomialEvaluationProofSystem(commP, poly)
			ps2 := crypto.NewDoubleDiscreteLogProofSystem(commP, commQ, params.SecurityParam)
			ps3 := crypto.NewPreimageEqualityProofSystem(params.HHat.BigInt(), commQ)

			filePath := path.Join(viper.GetString(cli.HomeFlag), votesFileName)
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("couldn't create or open file '%s'\n%v", filePath, err)
			}
			defer f.Close()
			for _, b := range ballots {
				v := true
				v = v && ps1.Verify(b.Proof1, b.C.BigInt(), b.V)
				v = v && ps2.Verify(b.Proof2, b.C.BigInt(), b.D.BigInt(), b.V)
				v = v && ps3.Verify(b.Proof3, b.D.BigInt(), b.UHat.BigInt(), b.V)
				if v {
					if _, err := f.WriteString(fmt.Sprintf("%s\n", b.V)); err != nil {
						return fmt.Errorf("couldn't write to file '%s'\n%v", filePath, err)
					}
				}
			}

			return nil
		},
	}
}

func QueryBallots(cliCtx context.CLIContext, cdc *codec.Codec) ([]types.Ballot, error) {
	route := fmt.Sprintf("custom/%s/%s", types.BulletinBoardModuleName, keeper.QueryBallots)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		msg := sdk.AppendMsgToErr("failed querying ballots", err.Error())
		return []types.Ballot{}, sdk.ErrInternal(msg)
	}
	var out []types.Ballot
	cdc.MustUnmarshalJSON(res, &out)
	return out, nil
}

func GetCmdVoterCredentials(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "credentials",
		Short: "Retrieve the set of voter credentials",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			route := fmt.Sprintf("custom/%s/%s", queryRoute, keeper.QueryVoterCredentials)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				msg := sdk.AppendMsgToErr("failed querying voter credentials", err.Error())
				return sdk.ErrInternal(msg)
			}
			var out types.QueryResVoterCredentials
			cdc.MustUnmarshalJSON(res, &out)
			return cliCtx.PrintOutput(out)
		},
	}
}

// GetCmdParameters fetches the bulletin board's set of parameters
func GetCmdParameters(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params [file]",
		Short: "Retrieve the set of bulletin board parameters and store them in a file.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params, err := QueryBulletinBoardParameters(cliCtx, cdc)
			if err != nil {
				return err
			}
			filePath := getFileName(args, 1, defaultParamsFileName)
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("couldn't create or open file '%s'\n%v", filePath, err)
			}
			defer f.Close()
			json, err := cdc.MarshalJSON(params)
			if err != nil {
				return fmt.Errorf("error marshalling params to json\n%v", err)
			}
			_, err = f.WriteString(string(json))
			if err != nil {
				return fmt.Errorf("error writing params to '%s'\n%v", filePath, err)
			}
			return cliCtx.PrintOutput(params)
		},
	}
}

func getFileName(args []string, argPosition int, defaultFileName string) string {
	if len(args) >= argPosition {
		return args[argPosition-1]
	} else {
		return path.Join(viper.GetString(cli.HomeFlag), defaultFileName)
	}
}

func QueryBulletinBoardParameters(cliCtx context.CLIContext, cdc *codec.Codec) (types.Params, error) {
	route := fmt.Sprintf("custom/%s/%s", types.BulletinBoardModuleName, keeper.QueryParameters)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		msg := sdk.AppendMsgToErr("failed querying parameters", err.Error())
		return types.Params{}, sdk.ErrInternal(msg)
	}
	var params types.Params
	cdc.MustUnmarshalJSON(res, &params)
	return params, nil
}

func GetCmdCredentialPolynomial(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "poly [file]",
		Short: "Retrieve credential polynomial and save it to file.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			poly, err := QueryCredentialPolynomial(cliCtx, cdc)
			if err != nil {
				return err
			}
			filePath := getFileName(args, 1, defaultPolynomialFileName)
			f, err := os.Create(filePath)
			if err != nil {
				return fmt.Errorf("couldn't create or open file '%s'\n%v", filePath, err)
			}
			defer f.Close()
			json, err := cdc.MarshalJSON(poly)
			if err != nil {
				return fmt.Errorf("error marshalling polynomial to json\n%v", err)
			}
			_, err = f.WriteString(string(json))
			if err != nil {
				return fmt.Errorf("error writing polynomial to '%s'\n%v", filePath, err)
			}
			return cliCtx.PrintOutput(&poly)
		},
	}
}

func QueryCredentialPolynomial(cliCtx context.CLIContext, cdc *codec.Codec) (crypto.Polynomial,
	error) {

	route := fmt.Sprintf("custom/%s/%s", types.BulletinBoardModuleName, keeper.QueryCredentialPolynomial)
	res, _, err := cliCtx.QueryWithData(route, nil)
	if err != nil {
		return crypto.Polynomial{}, fmt.Errorf("failed querying credential polynomial\n%v", err)
	}
	var poly crypto.Polynomial
	cdc.MustUnmarshalJSON(res, &poly)
	return poly, nil
}
