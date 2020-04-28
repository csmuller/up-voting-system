package cli

import (
	"errors"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/csmuller/uep-voting/crypto"
	"github.com/csmuller/uep-voting/pbb/internal/types"
	"github.com/spf13/cobra"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultPubCredFileName  = "cred.pub"
	defaultPrivCredFileName = "cred.priv"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	bulletinBoardTxCmd := &cobra.Command{
		Use:                        types.BulletinBoardModuleName,
		Short:                      "Bulletin Board transaction commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	bulletinBoardTxCmd.AddCommand(client.PostCommands(
		GetCmdGenerateAndPutVoterCredential(cdc),
		GetCmdGenerateAndPutBallot(cdc),
	)...)

	return bulletinBoardTxCmd
}

func GetCmdGenerateAndPutVoterCredential(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "new-voter [pub key file] [priv key file] [params file]",
		Short: "Generate new voter credentials and post the public credential to the bulletin board.",
		Args:  cobra.MaximumNArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params, err := readParameters(getFileName(args, 3, defaultParamsFileName), cdc)
			if err != nil {
				return err
			}
			voter := crypto.GenerateNewVoter(params.CommQ)
			// Create and send public credential transaction.
			msg := types.NewMsgPutVoterCredential(crypto.NewInt(voter.U), cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			txBuilder := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			if err := utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg}); err != nil {
				return err
			}
			// Write credentials to file.
			pubCredFileName := getFileName(args, 1, defaultPubCredFileName)
			privCredFileName := getFileName(args, 2, defaultPrivCredFileName)
			if err := writeCredentials(pubCredFileName, voter.U); err != nil {
				return err
			}
			if err := writeCredentials(privCredFileName, voter.A, voter.B); err != nil {
				return err
			}
			return nil
		},
	}
}

func readParameters(paramsFileName string, cdc *codec.Codec) (types.Params, error) {
	paramBytes, err := readFile(paramsFileName)
	if err != nil {
		return types.Params{}, err
	}
	var params types.Params
	if err := cdc.UnmarshalJSON(paramBytes, &params); err != nil {
		return types.Params{}, fmt.Errorf("failed unmarschalling parameters.\n%v", err)
	}
	return params, nil
}

func readPolynomial(polynomialFileName string, cdc *codec.Codec) (crypto.Polynomial, error) {
	polyBytes, err := readFile(polynomialFileName)
	if err != nil {
		return crypto.Polynomial{}, err
	}
	var poly crypto.Polynomial
	if err := cdc.UnmarshalJSON(polyBytes, &poly); err != nil {
		return crypto.Polynomial{}, fmt.Errorf("failed unmarschalling polynomial.\n%v", err)
	}
	return poly, nil
}

func writeCredentials(filePath string, credentials ...*big.Int) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("couldn't create or open file '%s'\n%v", filePath, err)
	}
	defer f.Close()
	out := new(strings.Builder)
	for i, cred := range credentials {
		out.WriteString(fmt.Sprintf("%s", cred.String()))
		if i != len(credentials)-1 {
			out.WriteString("\n")
		}
	}
	_, err = f.WriteString(out.String())
	if err != nil {
		return fmt.Errorf("error writing the private credentials to '%s'\n%v", filePath, err)
	}
	return nil
}

func GetCmdGenerateAndPutBallot(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "vote [vote] [pub key file] [priv key file] [params file] [poly file]",
		Short: "Generate a ballot with the given vote and post it to the bulletin board.",
		Args:  cobra.RangeArgs(1, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params, err := readParameters(getFileName(args, 4, defaultParamsFileName), cdc)
			if err != nil {
				return err
			}
			if params.HHat.BigInt() == nil {
				return errors.New("election generator has not been defined yet; " +
					"you might need to query the parameters again")
			}
			poly, err := readPolynomial(getFileName(args, 5, defaultPolynomialFileName), cdc)
			if err != nil {
				return err
			}
			pubKeyFile := getFileName(args, 2, defaultPubCredFileName)
			privKeyFile := getFileName(args, 3, defaultPrivCredFileName)
			voter, err := getVoterFromFiles(pubKeyFile, privKeyFile)
			if err != nil {
				return fmt.Errorf("failed fetching voter's credentials\n%v", err)
			}

			commP := params.CommP
			commQ := params.CommQ

			// election credential u_hat
			uHat := commQ.G.Exp(params.HHat.BigInt(), voter.B)

			// commitment c
			commToURand := commP.G.ZModOrder().RandomElement()
			commToU := commP.Commit(commToURand, voter.U)

			// commitment d
			commToAandBRand := commQ.G.ZModOrder().RandomElement()
			commToAandB := commQ.Commit(commToAandBRand, voter.A, voter.B)

			vote := args[0]

			// 1. proof
			ps1 := crypto.NewPolynomialEvaluationProofSystem(commP, poly)
			proof1 := ps1.Generate(voter.U, commToURand, commToU, vote)

			// 2. proof
			ps2 := crypto.NewDoubleDiscreteLogProofSystem(commP, commQ, params.SecurityParam)
			proof2 := ps2.Generate(voter, commToU, commToURand, commToAandB, commToAandBRand, vote)

			// 3. proof
			ps := crypto.NewPreimageEqualityProofSystem(params.HHat.BigInt(), commQ)
			proof3 := ps.Generate(voter, commToAandB, commToAandBRand, uHat, vote)

			// Create and send public credential transaction.
			msg := types.NewMsgPutBallot(commToU, commToAandB, vote, uHat, proof1, proof2,
				proof3, cliCtx.GetFromAddress())
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			txBuilder := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBuilder, []sdk.Msg{msg})
		},
	}
}

func getVoterFromFiles(pubCredFileName, privCredsFileName string) (crypto.Voter, error) {
	pubCredString, err := readFile(pubCredFileName)
	if err != nil {
		return crypto.Voter{}, err
	}
	pubCred, b := new(big.Int).SetString(string(pubCredString), 10)
	if !b {
		return crypto.Voter{}, fmt.Errorf("failed parsing public credential.\n%v", err)
	}

	privCredsString, err := readFile(privCredsFileName)
	if err != nil {
		return crypto.Voter{}, err
	}
	privCreds := strings.Split(string(privCredsString), "\n")
	alpha, _ := new(big.Int).SetString(privCreds[0], 10)
	beta, _ := new(big.Int).SetString(privCreds[1], 10)
	return crypto.NewVoter(alpha, beta, pubCred), nil
}

func readFile(arg string) ([]byte, error) {
	absPath, err := filepath.Abs(arg)
	if err != nil {
		return nil, fmt.Errorf("couldn't get absolute path for file %s\n%v", arg, err)
	}
	file, err := os.Open(absPath)
	if err != nil {
		return nil, fmt.Errorf("couldn't open file %s\n%v", arg, err)
	}
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("couldn't read all from file %s\n%v", arg, err)
	}
	return b, nil
}
