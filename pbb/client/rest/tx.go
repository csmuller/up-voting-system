package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/csmuller/uep-voting/crypto"
	"github.com/csmuller/uep-voting/pbb/internal/types"
	"net/http"
)

type putVoterCredentialsReq struct {
	BaseReq    rest.BaseReq `json:"base_req"`
	Name       string       `json:"tx_name"` // name of the tx
	Credential crypto.Int   `json:"credential"`
}

func putVoterCredentialsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req putVoterCredentialsReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}
		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}
		msg := types.NewMsgPutVoterCredential(req.Credential, cliCtx.GetFromAddress())
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		utils.WriteGenerateStdTxResponse(w, cliCtx, baseReq, []sdk.Msg{msg})
	}
}
