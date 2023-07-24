package thorchain

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"

	// tcTypes "gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/thorchain/thorchain"
	tcTypes "gitlab.com/mayachain/mayanode/x/mayachain/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	ctypes "github.com/cosmos/cosmos-sdk/types"
	rpcclient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/rs/zerolog/log"
	"gitlab.com/mayachain/mayanode/bifrost/metrics"
	"gitlab.com/mayachain/mayanode/config"

	"gitlab.com/mayachain/mayanode/bifrost/mayaclient"
	"gitlab.com/mayachain/mayanode/cmd"
	"gitlab.com/mayachain/mayanode/common"
	. "gopkg.in/check.v1"
)

// -------------------------------------------------------------------------------------
// Mock FeeTx
// -------------------------------------------------------------------------------------

type MockFeeTx struct {
	fee ctypes.Coins
	gas uint64
}

func (m *MockFeeTx) GetMsgs() []ctypes.Msg {
	return nil
}

func (m *MockFeeTx) ValidateBasic() error {
	return nil
}

func (m *MockFeeTx) GetGas() uint64 {
	return m.gas
}

func (m *MockFeeTx) GetFee() ctypes.Coins {
	return m.fee
}

func (m *MockFeeTx) FeePayer() ctypes.AccAddress {
	return nil
}

func (m *MockFeeTx) FeeGranter() ctypes.AccAddress {
	return nil
}

// -------------------------------------------------------------------------------------
// Tests
// -------------------------------------------------------------------------------------

type BlockScannerTestSuite struct {
	m      *metrics.Metrics
	bridge mayaclient.MayachainBridge
	keys   *mayaclient.Keys
}

var _ = Suite(&BlockScannerTestSuite{})

func (s *BlockScannerTestSuite) SetUpSuite(c *C) {
	s.m = GetMetricForTest(c)
	c.Assert(s.m, NotNil)
	cfg := config.BifrostClientConfiguration{
		ChainID:         "thorchain",
		ChainHost:       "localhost",
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: "",
	}

	kb := cKeys.NewInMemory()
	_, _, err := kb.NewMnemonic(cfg.SignerName, cKeys.English, cmd.BASEChainHDPath, cfg.SignerPasswd, hd.Secp256k1)
	c.Assert(err, IsNil)
	thorKeys := mayaclient.NewKeysWithKeybase(kb, cfg.SignerName, cfg.SignerPasswd)
	c.Assert(err, IsNil)
	s.bridge, err = mayaclient.NewMayachainBridge(cfg, s.m, thorKeys)
	c.Assert(err, IsNil)
	s.keys = thorKeys
}

// TC has a fixed fee
// func (s *BlockScannerTestSuite) TestCalculateAverageGasFees(c *C) {
// 	cfg := config.BifrostBlockScannerConfiguration{ChainID: common.THORChain}
// 	blockScanner := CosmosBlockScanner{cfg: cfg}
//
// 	blockScanner.updateGasCache(&MockFeeTx{
// 		gas: GasLimit / 2,
// 		fee: ctypes.Coins{ctypes.NewCoin("rune", ctypes.NewInt(1000000))},
// 	})
// 	c.Check(len(blockScanner.feeCache), Equals, 1)
// 	c.Check(blockScanner.averageFee().Uint64(), Equals, uint64(2000000), Commentf("expected %s, got %d", blockScanner.averageFee().String(), 2000000))
//
// 	blockScanner.updateGasCache(&MockFeeTx{
// 		gas: GasLimit / 2,
// 		fee: ctypes.Coins{ctypes.NewCoin("rune", ctypes.NewInt(1000000))},
// 	})
// 	c.Check(len(blockScanner.feeCache), Equals, 2)
// 	c.Check(blockScanner.averageFee().Uint64(), Equals, uint64(2000000), Commentf("expected %s, got %d", blockScanner.averageFee().String(), 2000000))
//
// 	// two blocks at half fee should average to 75% of last
// 	blockScanner.updateGasCache(&MockFeeTx{
// 		gas: GasLimit,
// 		fee: ctypes.Coins{ctypes.NewCoin("rune", ctypes.NewInt(1000000))},
// 	})
// 	blockScanner.updateGasCache(&MockFeeTx{
// 		gas: GasLimit,
// 		fee: ctypes.Coins{ctypes.NewCoin("rune", ctypes.NewInt(1000000))},
// 	})
// 	c.Check(len(blockScanner.feeCache), Equals, 4)
// 	c.Check(blockScanner.averageFee().Uint64(), Equals, uint64(1500000), Commentf("expected %s, got %d", blockScanner.averageFee().String(), 1500000))
//
// 	// skip transactions with multiple coins
// 	blockScanner.updateGasCache(&MockFeeTx{
// 		gas: GasLimit,
// 		fee: ctypes.Coins{
// 			ctypes.NewCoin("rune", ctypes.NewInt(1000000)),
// 			ctypes.NewCoin("deletethis", ctypes.NewInt(1000000)),
// 		},
// 	})
// 	c.Check(len(blockScanner.feeCache), Equals, 4)
// 	c.Check(blockScanner.averageFee().Uint64(), Equals, uint64(1500000), Commentf("expected %s, got %d", blockScanner.averageFee().String(), 15000))
//
// 	// skip transactions with zero fee
// 	blockScanner.updateGasCache(&MockFeeTx{
// 		gas: GasLimit,
// 		fee: ctypes.Coins{
// 			ctypes.NewCoin("rune", ctypes.NewInt(0)),
// 		},
// 	})
// 	c.Check(len(blockScanner.feeCache), Equals, 4)
// 	c.Check(blockScanner.averageFee().Uint64(), Equals, uint64(1500000), Commentf("expected %s, got %d", blockScanner.averageFee().String(), 15000))
//
// 	// ensure we only cache the transaction limit number of blocks
// 	for i := 0; i < GasCacheTransactions; i++ {
// 		blockScanner.updateGasCache(&MockFeeTx{
// 			gas: GasLimit,
// 			fee: ctypes.Coins{
// 				ctypes.NewCoin("rune", ctypes.NewInt(1000000)),
// 			},
// 		})
// 	}
// 	c.Check(len(blockScanner.feeCache), Equals, GasCacheTransactions)
// 	c.Check(blockScanner.averageFee().Uint64(), Equals, uint64(1000000), Commentf("expected %s, got %d", blockScanner.averageFee().String(), 15000))
// }

func (s *BlockScannerTestSuite) TestGetBlock(c *C) {
	cfg := config.BifrostBlockScannerConfiguration{ChainID: common.THORChain}
	mockRPC := NewMockTmServiceClient()

	blockScanner := CosmosBlockScanner{
		cfg:       cfg,
		tmService: mockRPC,
	}

	block, err := blockScanner.GetBlock(1)

	c.Assert(err, IsNil)
	c.Assert(len(block.Data.Txs), Equals, 1)
	c.Assert(block.Header.Height, Equals, int64(162))
}

func (s *BlockScannerTestSuite) TestProcessTxs(c *C) {
	cfg := config.BifrostBlockScannerConfiguration{ChainID: common.THORChain}
	mockTmServiceClient := NewMockTmServiceClient()
	c.Assert(os.Setenv("NET", "testnet"), IsNil)

	registry := s.bridge.GetContext().InterfaceRegistry
	registry.RegisterImplementations((*ctypes.Msg)(nil), &tcTypes.MsgSend{})

	cdc := codec.NewProtoCodec(registry)

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jsonFile, err := os.Open("./test-data/tx_results_by_height.json")
		if err != nil {
			c.Fatal("unable to load tx_results_by_height.json")
		}
		defer jsonFile.Close()

		res, _ := io.ReadAll(jsonFile)
		if _, err := w.Write(res); err != nil {
			c.Fatal("unable to write /block_result", err)
		}
	})
	server := httptest.NewServer(h)
	defer server.Close()

	rpcClient, err := rpcclient.NewWithClient(server.URL, "/websocket", server.Client())
	if err != nil {
		c.Fatal("fail to create tendermint rpcclient")
	}

	blockScanner := CosmosBlockScanner{
		cfg:       cfg,
		tmService: mockTmServiceClient,
		txService: rpcClient,
		cdc:       cdc,
		logger:    log.Logger.With().Str("module", "blockscanner").Str("chain", common.THORChain.String()).Logger(),
	}

	block, err := blockScanner.GetBlock(1)
	c.Assert(err, IsNil)

	txInItems, err := blockScanner.processTxs(1, block.Data.Txs)
	c.Assert(err, IsNil)

	// Make sure AccAddress to String conversion is working
	c.Assert(txInItems[0].Sender, Equals, "tthor18z343fsdlav47chtkyp0aawqt6sgxsh3gjtkrh")
	c.Assert(txInItems[0].To, Equals, "tthor1xdmg4gfz8yvermj5yudn987nvs7f6zafq72m0q")

	// proccessTxs should filter out everything besides the valid MsgSend
	c.Assert(len(txInItems), Equals, 1)
}
