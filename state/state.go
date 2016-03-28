package state

import (
	"github.com/tendermint/basecoin/types"
	. "github.com/tendermint/go-common"
	"github.com/tendermint/go-wire"
	eyes "github.com/tendermint/merkleeyes/client"
)

type State struct {
	chainID    string
	eyesCli    *eyes.Client
	checkCache map[string]checkAccount

	LastBlockHeight uint64
	LastBlockHash   []byte
	GasLimit        int64
}

func NewState(eyesCli *eyes.Client) *State {
	s := &State{
		chainID:    "",
		eyesCli:    eyesCli,
		checkCache: make(map[string]checkAccount),
	}
	return s
}

func (s *State) SetChainID(chainID string) {
	s.chainID = chainID
}

func (s *State) GetChainID() string {
	if s.chainID == "" {
		PanicSanity("Expected to have set SetChainID")
	}
	return s.chainID
}

//----------------------------------------
// CheckTx state

type checkAccount struct {
	sequence int
	balance  int64
}

func (s *State) GetCheckAccount(addr []byte, defaultSequence int, defaultBalance int64) (sequence int, balance int64) {
	cAcc, ok := s.checkCache[string(addr)]
	if !ok {
		return defaultSequence, defaultBalance
	}
	return cAcc.sequence, cAcc.balance
}

func (s *State) SetCheckAccount(addr []byte, sequence int, balance int64) {
	s.checkCache[string(addr)] = checkAccount{sequence, balance}
}

func (s *State) ResetCacheState() {
	s.checkCache = make(map[string]checkAccount)
}

//----------------------------------------

func (s *State) GetAccount(addr []byte) *types.Account {
	res := s.eyesCli.GetSync(addr)
	if res.IsErr() {
		panic(Fmt("Error loading account addr %X error: %v", addr, res.Error()))
	}
	if len(res.Data) == 0 {
		return nil
	}
	var acc *types.Account
	err := wire.ReadBinaryBytes(res.Data, &acc)
	if err != nil {
		panic(Fmt("Error reading account %X error: %v", res.Data, err.Error()))
	}
	return acc
}

func (s *State) SetAccount(address []byte, acc *types.Account) {
	accBytes := wire.BinaryBytes(acc)
	res := s.eyesCli.SetSync(address, accBytes)
	if res.IsErr() {
		panic(Fmt("Error storing account addr %X error: %v", address, res.Error()))
	}
}
