package rest

import (
	"io/ioutil"
	"net/http"

	cliclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/bank/client"

	"github.com/gorilla/mux"
)

type sendBody struct {
	// fees is not used currently
	// Fees             sdk.Coin  `json="fees"`
	Amount           sdk.Coins `json:"amount"`
	LocalAccountName string    `json:"name"`
	Password         string    `json:"password"`
	ChainID          string    `json:"chain_id"`
	AccountNumber    int64     `json:"account_number"`
	Sequence         int64     `json:"sequence"`
	Gas              string    `json:"gas"`
	GasAdjustment    string    `json:"gas_adjustment"`
}

var msgCdc = wire.NewCodec()

func init() {
	bank.RegisterWire(msgCdc)
}

// SendRequestHandlerFn - http request handler to send coins to a address
// nolint: gocyclo
func SendRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// collect data
		vars := mux.Vars(r)
		bech32addr := vars["address"]

		to, err := sdk.AccAddressFromBech32(bech32addr)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var m sendBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		err = msgCdc.UnmarshalJSON(body, &m)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		// build message
		msg := client.BuildMsg(sdk.AccAddress(info.GetPubKey().Address()), to, m.Amount)
		if err != nil { // XXX rechecking same error ?
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		simulateGas, gas, err := cliclient.ReadGasFlag(m.Gas)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		adjustment, ok := utils.ParseFloat64OrReturnBadRequest(w, m.GasAdjustment, cliclient.DefaultGasAdjustment)
		if !ok {
			return
		}
		txBldr := authtxb.TxBuilder{
			Codec:         cdc,
			Gas:           gas,
			GasAdjustment: adjustment,
			SimulateGas:   simulateGas,
			ChainID:       m.ChainID,
			AccountNumber: m.AccountNumber,
			Sequence:      m.Sequence,
		}

		if utils.HasDryRunArg(r) || txBldr.SimulateGas {
			newBldr, err := utils.EnrichCtxWithGas(txBldr, cliCtx, m.LocalAccountName, []sdk.Msg{msg})
			if err != nil {
				utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
			if utils.HasDryRunArg(r) {
				utils.WriteSimulationResponse(w, newBldr.Gas)
				return
			}
			txBldr = newBldr
		}

		if utils.HasGenerateOnlyArg(r) {
			utils.WriteGenerateStdTxResponse(w, txBldr, []sdk.Msg{msg})
			return
		}

		txBytes, err := txBldr.BuildAndSign(m.LocalAccountName, m.Password, []sdk.Msg{msg})
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusUnauthorized, err.Error())
			return
		}

		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		output, err := wire.MarshalJSONIndent(cdc, res)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		w.Write(output)
	}
}
