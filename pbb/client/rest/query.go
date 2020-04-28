package rest

import (
	"fmt"
	"github.com/csmuller/uep-voting/pbb/internal/keeper"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/cosmos/cosmos-sdk/types/rest"
)

func ballotsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, keeper.QueryBallots)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func voterCredentialsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := fmt.Sprintf("custom/%s/%s", storeName, keeper.QueryVoterCredentials)
		res, _, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//func resolveNameHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r)
//		paramType := vars[restName]
//
//		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/ballots/%s", storeName,
//			paramType), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func whoIsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r)
//		paramType := vars[restName]
//
//		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/whois/%s", storeName, paramType), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
//
//func namesHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/names", storeName), nil)
//		if err != nil {
//			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
//			return
//		}
//		rest.PostProcessResponse(w, cliCtx, res)
//	}
//}
