package lcd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/abci/types"
	cryptoKeys "github.com/tendermint/go-crypto/keys"
	tmcfg "github.com/tendermint/tendermint/config"
	nm "github.com/tendermint/tendermint/node"
	p2p "github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/proxy"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmrpc "github.com/tendermint/tendermint/rpc/lib/server"
	tmtypes "github.com/tendermint/tendermint/types"
	pvm "github.com/tendermint/tendermint/types/priv_validator"
	"github.com/tendermint/tmlibs/cli"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"

	client "github.com/cosmos/cosmos-sdk/client"
	keys "github.com/cosmos/cosmos-sdk/client/keys"
	bapp "github.com/cosmos/cosmos-sdk/examples/basecoin/app"
	btypes "github.com/cosmos/cosmos-sdk/examples/basecoin/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	coinDenom  = "mycoin"
	coinAmount = int64(10000000)

	// XXX bad globals
	port     string // XXX: but it's the int ...
	name     string = "test"
	password string = "0123456789"
	seed     string
	sendAddr string
)

func TestKeys(t *testing.T) {

	// empty keys
	// XXX: the test comes with a key setup
	/*
		res, body := request(t, port, "GET", "/keys", nil)
		require.Equal(t, http.StatusOK, res.StatusCode, body)
		assert.Equal(t, "[]", body, "Expected an empty array")
	*/

	// get seed
	res, body := request(t, port, "GET", "/keys/seed", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	newSeed := body
	reg, err := regexp.Compile(`([a-z]+ ){12}`)
	require.Nil(t, err)
	match := reg.MatchString(seed)
	assert.True(t, match, "Returned seed has wrong foramt", seed)

	newName := "test_newname"
	newPassword := "0987654321"

	// add key
	var jsonStr = []byte(fmt.Sprintf(`{"name":"test_fail", "password":"%s"}`, password))
	res, body = request(t, port, "POST", "/keys", jsonStr)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode, "Account creation should require a seed")

	jsonStr = []byte(fmt.Sprintf(`{"name":"%s", "password":"%s", "seed": "%s"}`, newName, newPassword, newSeed))
	res, body = request(t, port, "POST", "/keys", jsonStr)

	require.Equal(t, http.StatusOK, res.StatusCode, body)
	addr := body
	assert.Len(t, addr, 40, "Returned address has wrong format", addr)

	// existing keys
	res, body = request(t, port, "GET", "/keys", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var m [2]keys.KeyOutput
	err = cdc.UnmarshalJSON([]byte(body), &m)
	require.Nil(t, err)

	assert.Equal(t, m[0].Name, name, "Did not serve keys name correctly")
	assert.Equal(t, m[0].Address, sendAddr, "Did not serve keys Address correctly")
	assert.Equal(t, m[1].Name, newName, "Did not serve keys name correctly")
	assert.Equal(t, m[1].Address, addr, "Did not serve keys Address correctly")

	// select key
	keyEndpoint := fmt.Sprintf("/keys/%s", newName)
	res, body = request(t, port, "GET", keyEndpoint, nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var m2 keys.KeyOutput
	err = cdc.UnmarshalJSON([]byte(body), &m2)
	require.Nil(t, err)

	assert.Equal(t, newName, m2.Name, "Did not serve keys name correctly")
	assert.Equal(t, addr, m2.Address, "Did not serve keys Address correctly")

	// update key
	jsonStr = []byte(fmt.Sprintf(`{"old_password":"%s", "new_password":"12345678901"}`, newPassword))
	res, body = request(t, port, "PUT", keyEndpoint, jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// here it should say unauthorized as we changed the password before
	res, body = request(t, port, "PUT", keyEndpoint, jsonStr)
	require.Equal(t, http.StatusUnauthorized, res.StatusCode, body)

	// delete key
	jsonStr = []byte(`{"password":"12345678901"}`)
	res, body = request(t, port, "DELETE", keyEndpoint, jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
}

func TestVersion(t *testing.T) {

	// node info
	res, body := request(t, port, "GET", "/version", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	reg, err := regexp.Compile(`\d+\.\d+\.\d+(-dev)?`)
	require.Nil(t, err)
	match := reg.MatchString(body)
	assert.True(t, match, body)
}

func TestNodeStatus(t *testing.T) {

	// node info
	res, body := request(t, port, "GET", "/node_info", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	var nodeInfo p2p.NodeInfo
	err := cdc.UnmarshalJSON([]byte(body), &nodeInfo)
	require.Nil(t, err, "Couldn't parse node info")

	assert.NotEqual(t, p2p.NodeInfo{}, nodeInfo, "res: %v", res)

	// syncing
	res, body = request(t, port, "GET", "/syncing", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// we expect that there is no other node running so the syncing state is "false"
	// we c
	assert.Equal(t, "false", body)
}

func TestBlock(t *testing.T) {

	waitForHeight(2)

	var resultBlock ctypes.ResultBlock

	res, body := request(t, port, "GET", "/blocks/latest", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err := cdc.UnmarshalJSON([]byte(body), &resultBlock)
	require.Nil(t, err, "Couldn't parse block")

	assert.NotEqual(t, ctypes.ResultBlock{}, resultBlock)

	// --

	res, body = request(t, port, "GET", "/blocks/1", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err = json.Unmarshal([]byte(body), &resultBlock)
	require.Nil(t, err, "Couldn't parse block")

	assert.NotEqual(t, ctypes.ResultBlock{}, resultBlock)

	// --

	res, body = request(t, port, "GET", "/blocks/1000000000", nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode, body)
}

func TestValidators(t *testing.T) {

	var resultVals ctypes.ResultValidators

	res, body := request(t, port, "GET", "/validatorsets/latest", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err := cdc.UnmarshalJSON([]byte(body), &resultVals)
	require.Nil(t, err, "Couldn't parse validatorset")

	assert.NotEqual(t, ctypes.ResultValidators{}, resultVals)

	// --

	res, body = request(t, port, "GET", "/validatorsets/1", nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &resultVals)
	require.Nil(t, err, "Couldn't parse validatorset")

	assert.NotEqual(t, ctypes.ResultValidators{}, resultVals)

	// --

	res, body = request(t, port, "GET", "/validatorsets/1000000000", nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

func TestCoinSend(t *testing.T) {

	// query empty
	res, body := request(t, port, "GET", "/accounts/8FA6AB57AD6870F6B5B2E57735F38F2F30E73CB6", nil)
	require.Equal(t, http.StatusNoContent, res.StatusCode, body)

	acc := getAccount(t, sendAddr)
	initialBalance := acc.GetCoins()

	// create TX
	receiveAddr, resultTx := doSend(t, port, seed)
	waitForHeight(resultTx.Height + 1)

	// check if tx was commited
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// query sender
	acc = getAccount(t, sendAddr)
	coins := acc.GetCoins()
	mycoins := coins[0]
	assert.Equal(t, coinDenom, mycoins.Denom)
	assert.Equal(t, initialBalance[0].Amount-1, mycoins.Amount)

	// query receiver
	acc = getAccount(t, receiveAddr)
	coins = acc.GetCoins()
	mycoins = coins[0]
	assert.Equal(t, coinDenom, mycoins.Denom)
	assert.Equal(t, int64(1), mycoins.Amount)
}

func TestIBCTransfer(t *testing.T) {

	acc := getAccount(t, sendAddr)
	initialBalance := acc.GetCoins()

	// create TX
	resultTx := doIBCTransfer(t, port, seed)

	waitForHeight(resultTx.Height + 1)

	// check if tx was commited
	assert.Equal(t, uint32(0), resultTx.CheckTx.Code)
	assert.Equal(t, uint32(0), resultTx.DeliverTx.Code)

	// query sender
	acc = getAccount(t, sendAddr)
	coins := acc.GetCoins()
	mycoins := coins[0]
	assert.Equal(t, coinDenom, mycoins.Denom)
	assert.Equal(t, initialBalance[0].Amount-1, mycoins.Amount)

	// TODO: query ibc egress packet state
}

func TestTxs(t *testing.T) {

	// TODO: re-enable once we can get txs by tag

	// query wrong
	// res, body := request(t, port, "GET", "/txs", nil)
	// require.Equal(t, http.StatusBadRequest, res.StatusCode, body)

	// query empty
	// res, body = request(t, port, "GET", fmt.Sprintf("/txs?tag=coin.sender='%s'", "8FA6AB57AD6870F6B5B2E57735F38F2F30E73CB6"), nil)
	// require.Equal(t, http.StatusOK, res.StatusCode, body)

	// assert.Equal(t, "[]", body)

	// create TX
	_, resultTx := doSend(t, port, seed)

	waitForHeight(resultTx.Height + 1)

	// check if tx is findable
	res, body := request(t, port, "GET", fmt.Sprintf("/txs/%s", resultTx.Hash), nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	// // query sender
	// res, body = request(t, port, "GET", fmt.Sprintf("/txs?tag=coin.sender='%s'", addr), nil)
	// require.Equal(t, http.StatusOK, res.StatusCode, body)

	// assert.NotEqual(t, "[]", body)

	// // query receiver
	// res, body = request(t, port, "GET", fmt.Sprintf("/txs?tag=coin.receiver='%s'", receiveAddr), nil)
	// require.Equal(t, http.StatusOK, res.StatusCode, body)

	// assert.NotEqual(t, "[]", body)
}

//__________________________________________________________
// helpers

// strt TM and the LCD in process, listening on their respective sockets
func startTMAndLCD() (*nm.Node, net.Listener, error) {

	dir, err := ioutil.TempDir("", "lcd_test")
	if err != nil {
		return nil, nil, err
	}
	viper.Set(cli.HomeFlag, dir)
	kb, err := keys.GetKeyBase() // dbm.NewMemDB()) // :(
	if err != nil {
		return nil, nil, err
	}
	var info cryptoKeys.Info
	info, seed, err = kb.Create(name, password, cryptoKeys.AlgoEd25519) // XXX global seed
	if err != nil {
		return nil, nil, err
	}

	pubKey := info.PubKey
	sendAddr = pubKey.Address().String() // XXX global

	config := GetConfig()
	config.Consensus.TimeoutCommit = 1000
	config.Consensus.SkipTimeoutCommit = false

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	// logger = log.NewFilter(logger, log.AllowError())
	privValidatorFile := config.PrivValidatorFile()
	privVal := pvm.LoadOrGenFilePV(privValidatorFile)
	db := dbm.NewMemDB()
	app := bapp.NewBasecoinApp(logger, db)
	cdc = bapp.MakeCodec() // XXX

	genesisFile := config.GenesisFile()
	genDoc, err := tmtypes.GenesisDocFromFile(genesisFile)
	if err != nil {
		return nil, nil, err
	}

	coins := sdk.Coins{{coinDenom, coinAmount}}
	appState := map[string]interface{}{
		"accounts": []*btypes.GenesisAccount{
			{
				Name:    "tester",
				Address: pubKey.Address(),
				Coins:   coins,
			},
		},
	}
	stateBytes, err := json.Marshal(appState)
	if err != nil {
		return nil, nil, err
	}
	genDoc.AppStateJSON = stateBytes

	// LCD listen address
	port = fmt.Sprintf("%d", 17377)                       // XXX
	listenAddr := fmt.Sprintf("tcp://localhost:%s", port) // XXX

	// XXX: need to set this so LCD knows the tendermint node address!
	viper.Set(client.FlagNode, config.RPC.ListenAddress)
	viper.Set(client.FlagChainID, genDoc.ChainID)

	node, err := startTM(config, logger, genDoc, privVal, app)
	if err != nil {
		return nil, nil, err
	}
	lcd, err := startLCD(logger, listenAddr)
	if err != nil {
		return nil, nil, err
	}

	waitForStart()

	return node, lcd, nil
}

// Create & start in-process tendermint node with memdb
// and in-process abci application.
// TODO: need to clean up the WAL dir or enable it to be not persistent
func startTM(cfg *tmcfg.Config, logger log.Logger, genDoc *tmtypes.GenesisDoc, privVal tmtypes.PrivValidator, app abci.Application) (*nm.Node, error) {
	genDocProvider := func() (*tmtypes.GenesisDoc, error) { return genDoc, nil }
	dbProvider := func(*nm.DBContext) (dbm.DB, error) { return dbm.NewMemDB(), nil }
	n, err := nm.NewNode(cfg,
		privVal,
		proxy.NewLocalClientCreator(app),
		genDocProvider,
		dbProvider,
		logger.With("module", "node"))
	if err != nil {
		return nil, err
	}

	err = n.Start()
	if err != nil {
		return nil, err
	}

	// wait for rpc
	waitForRPC()

	logger.Info("Tendermint running!")
	return n, err
}

// start the LCD. note this blocks!
func startLCD(logger log.Logger, listenAddr string) (net.Listener, error) {
	handler := createHandler(cdc)
	return tmrpc.StartHTTPServer(listenAddr, handler, logger)
}

func request(t *testing.T, port, method, path string, payload []byte) (*http.Response, string) {
	var res *http.Response
	var err error
	url := fmt.Sprintf("http://localhost:%v%v", port, path)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	require.Nil(t, err)
	res, err = http.DefaultClient.Do(req)
	//	res, err = http.Post(url, "application/json", bytes.NewBuffer(payload))
	require.Nil(t, err)

	output, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	require.Nil(t, err)

	return res, string(output)
}

func getAccount(t *testing.T, sendAddr string) sdk.Account {
	// get the account to get the sequence
	res, body := request(t, port, "GET", "/accounts/"+sendAddr, nil)
	require.Equal(t, http.StatusOK, res.StatusCode, body)
	var acc sdk.Account
	err := cdc.UnmarshalJSON([]byte(body), &acc)
	require.Nil(t, err)
	return acc
}

func doSend(t *testing.T, port, seed string) (receiveAddr string, resultTx ctypes.ResultBroadcastTxCommit) {

	// create receive address
	kb := client.MockKeyBase()
	receiveInfo, _, err := kb.Create("receive_address", "1234567890", cryptoKeys.CryptoAlgo("ed25519"))
	require.Nil(t, err)
	receiveAddr = receiveInfo.PubKey.Address().String()

	acc := getAccount(t, sendAddr)
	sequence := acc.GetSequence()

	// send
	jsonStr := []byte(fmt.Sprintf(`{ "name":"%s", "password":"%s", "sequence":%d, "amount":[{ "denom": "%s", "amount": 1 }] }`, name, password, sequence, coinDenom))
	res, body := request(t, port, "POST", "/accounts/"+receiveAddr+"/send", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return receiveAddr, resultTx
}

func doIBCTransfer(t *testing.T, port, seed string) (resultTx ctypes.ResultBroadcastTxCommit) {

	// create receive address
	kb := client.MockKeyBase()
	receiveInfo, _, err := kb.Create("receive_address", "1234567890", cryptoKeys.CryptoAlgo("ed25519"))
	require.Nil(t, err)
	receiveAddr := receiveInfo.PubKey.Address().String()

	// get the account to get the sequence
	acc := getAccount(t, sendAddr)
	sequence := acc.GetSequence()

	// send
	jsonStr := []byte(fmt.Sprintf(`{ "name":"%s", "password":"%s", "sequence":%d, "amount":[{ "denom": "%s", "amount": 1 }] }`, name, password, sequence, coinDenom))
	res, body := request(t, port, "POST", "/ibc/testchain/"+receiveAddr+"/send", jsonStr)
	require.Equal(t, http.StatusOK, res.StatusCode, body)

	err = cdc.UnmarshalJSON([]byte(body), &resultTx)
	require.Nil(t, err)

	return resultTx
}

func waitForHeight(height int64) {
	for {
		var resultBlock ctypes.ResultBlock

		url := fmt.Sprintf("http://localhost:%v%v", port, "/blocks/latest")
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		res.Body.Close()

		err = cdc.UnmarshalJSON([]byte(body), &resultBlock)
		if err != nil {
			fmt.Println("RES", res)
			fmt.Println("BODY", string(body))
			panic(err)
		}

		if resultBlock.Block.Height >= height {
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
}

// wait for 2 blocks
func waitForStart() {
	waitHeight := int64(2)
	for {
		time.Sleep(time.Second)

		url := fmt.Sprintf("http://localhost:%v%v", port, "/blocks/latest")
		res, err := http.Get(url)
		if err != nil {
			panic(err)
		}

		// waiting for server to start ...
		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			continue
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		res.Body.Close()

		resultBlock := new(ctypes.ResultBlock)
		err = cdc.UnmarshalJSON([]byte(body), &resultBlock)
		if err != nil {
			fmt.Println("RES", res)
			fmt.Println("BODY", string(body))
			panic(err)
		}

		if resultBlock.Block.Height >= waitHeight {
			return
		}
	}
}
