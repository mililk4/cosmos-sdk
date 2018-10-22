package keys

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/tendermint/tendermint/crypto"
	"net/http"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto/multisig"
	"github.com/tendermint/tendermint/libs/cli"
)

const (
	// FlagAddress is the flag for the user's address on the command line.
	FlagAddress = "address"
	// FlagPublicKey represents the user's public key on the command line.
	FlagPublicKey = "pubkey"
	// FlagBechPrefix defines a desired Bech32 prefix encoding for a key.
	FlagBechPrefix = "bech"

	flagMultiSigThreshold = "multisig-threshold"
)

var _ keys.Info = (keys.Info)(nil)

type multiSigKey struct {
	name string
	key  crypto.PubKey
}

func (m multiSigKey) GetName() string            { return m.name }
func (m multiSigKey) GetType() keys.KeyType      { return keys.TypeLocal }
func (m multiSigKey) GetPubKey() crypto.PubKey   { return m.key }
func (m multiSigKey) GetAddress() sdk.AccAddress { return sdk.AccAddress(m.key.Address()) }

func showKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show [name]",
		Short: "Show key info for the given name",
		Long:  `Return public details of one local key.`,
		Args:  cobra.MinimumNArgs(1),
		RunE:  runShowCmd,
	}

	cmd.Flags().String(FlagBechPrefix, "acc", "The Bech32 prefix encoding for a key (acc|val|cons)")
	cmd.Flags().Bool(FlagAddress, false, "output the address only (overrides --output)")
	cmd.Flags().Bool(FlagPublicKey, false, "output the public key only (overrides --output)")
	cmd.Flags().UintP(flagMultiSigThreshold, "m", 1, "K out of N required signatures")

	return cmd
}

func runShowCmd(cmd *cobra.Command, args []string) (err error) {
	var info keys.Info

	if len(args) == 1 {
		info, err = GetKeyInfo(args[0])
		if err != nil {
			return err
		}
	} else {
		pks := make([]crypto.PubKey, len(args))
		for i, keyName := range args {
			info, err := GetKeyInfo(keyName)
			if err != nil {
				return err
			}
			pks[i] = info.GetPubKey()
		}
		multikey := multisig.NewPubKeyMultisigThreshold(viper.GetInt(flagMultiSigThreshold), pks)
		info = multiSigKey{
			name: "multi",
			key:  multikey,
		}
	}

	isShowAddr := viper.GetBool(FlagAddress)
	isShowPubKey := viper.GetBool(FlagPublicKey)
	isOutputSet := cmd.Flag(cli.OutputFlag).Changed

	if isShowAddr && isShowPubKey {
		return errors.New("cannot use both --address and --pubkey at once")
	}

	if isOutputSet && (isShowAddr || isShowPubKey) {
		return errors.New("cannot use --output with --address or --pubkey")
	}

	bechKeyOut, err := getBechKeyOut(viper.GetString(FlagBechPrefix))
	if err != nil {
		return err
	}

	switch {
	case isShowAddr:
		printKeyAddress(info, bechKeyOut)
	case isShowPubKey:
		printPubKey(info, bechKeyOut)
	default:
		printKeyInfo(info, bechKeyOut)
	}

	return nil
}

func getBechKeyOut(bechPrefix string) (bechKeyOutFn, error) {
	switch bechPrefix {
	case "acc":
		return Bech32KeyOutput, nil
	case "val":
		return Bech32ValKeyOutput, nil
	case "cons":
		return Bech32ConsKeyOutput, nil
	}

	return nil, fmt.Errorf("invalid Bech32 prefix encoding provided: %s", bechPrefix)
}

///////////////////////////
// REST

// get key REST handler
func GetKeyRequestHandler(indent bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		bechPrefix := r.URL.Query().Get(FlagBechPrefix)

		if bechPrefix == "" {
			bechPrefix = "acc"
		}

		bechKeyOut, err := getBechKeyOut(bechPrefix)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		info, err := GetKeyInfo(name)
		// TODO: check for the error if key actually does not exist, instead of
		// assuming this as the reason
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(err.Error()))
			return
		}

		keyOutput, err := bechKeyOut(info)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		PostProcessResponse(w, cdc, keyOutput, indent)
	}
}
