package dash

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	ctypes "gitlab.com/mayachain/binance-sdk/common/types"
	"gitlab.com/mayachain/dashd-go/btcjson"
	"gitlab.com/mayachain/mayanode/bifrost/mayaclient"
	"gitlab.com/mayachain/mayanode/bifrost/mayaclient/types"
	"gitlab.com/mayachain/mayanode/bifrost/metrics"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/shared/utxo"
	"gitlab.com/mayachain/mayanode/cmd"
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/config"
	ttypes "gitlab.com/mayachain/mayanode/x/mayachain/types"
	. "gopkg.in/check.v1"
)

const (
	bob      = "bob"
	password = "password"
)

func TestPackage(t *testing.T) { TestingT(t) }

type DashSuite struct {
	client *Client
	server *httptest.Server
	bridge mayaclient.MayachainBridge
	cfg    config.BifrostChainConfiguration
	m      *metrics.Metrics
}

var _ = Suite(
	&DashSuite{},
)

var m *metrics.Metrics

func GetMetricForTest(c *C) *metrics.Metrics {
	if m == nil {
		var err error
		m, err = metrics.NewMetrics(config.BifrostMetricsConfiguration{
			Enabled:      false,
			ListenPort:   9000,
			ReadTimeout:  time.Second / 10,
			WriteTimeout: time.Second / 10,
			Chains:       common.Chains{common.DASHChain},
		})
		c.Assert(m, NotNil)
		c.Assert(err, IsNil)
	}
	return m
}

func (s *DashSuite) SetUpTest(c *C) {
	s.m = GetMetricForTest(c)
	s.cfg = config.BifrostChainConfiguration{
		ChainID:     "DASH",
		UserName:    "bob",
		Password:    "password",
		DisableTLS:  true,
		HTTPostMode: true,
		BlockScanner: config.BifrostBlockScannerConfiguration{
			StartBlockHeight: 1, // avoids querying mayachain for block height
		},
	}
	ns := strconv.Itoa(time.Now().Nanosecond())
	ttypes.SetupConfigForTest()
	ctypes.Network = ctypes.TestNetwork
	c.Assert(os.Setenv("NET", "testnet"), IsNil)

	mayadir := filepath.Join(os.TempDir(), ns, ".mayacli")
	cfg := config.BifrostClientConfiguration{
		ChainID:         "mayachain",
		ChainHost:       "localhost",
		SignerName:      "bob",
		SignerPasswd:    "password",
		ChainHomeFolder: mayadir,
	}

	kb := cKeys.NewInMemory()
	_, _, err := kb.NewMnemonic(bob, cKeys.English, cmd.BASEChainHDPath, password, hd.Secp256k1)
	c.Assert(err, IsNil)
	mayaKeys := mayaclient.NewKeysWithKeybase(kb, cfg.SignerName, cfg.SignerPasswd)
	c.Assert(err, IsNil)

	s.server = httptest.NewServer(mockHttpResponses(c))
	cfg.ChainHost = s.server.Listener.Addr().String()
	s.bridge, err = mayaclient.NewMayachainBridge(cfg, s.m, mayaKeys)
	c.Assert(err, IsNil)
	s.cfg.RPCHost = s.server.Listener.Addr().String()
	s.client, err = NewClient(mayaKeys, s.cfg, nil, s.bridge, s.m)
	c.Assert(err, IsNil)
	c.Assert(s.client, NotNil)
}

func (s *DashSuite) TearDownTest(c *C) {
	s.server.Close()
}

func (s *DashSuite) TestGetBlock(c *C) {
	block, err := s.client.getBlock(1696761)
	c.Assert(err, IsNil)
	c.Assert(block.Hash, Equals, "0000000000000002cb2f0d030e9744cd6d93b913b80958b232e4f0b84d828ecf")
	c.Assert(block.Tx[0].Txid, Equals, "cc80c1b9768195613e67703c6aca3df45f0705f585da8498d1ffaef3b746352a")
	c.Assert(len(block.Tx), Equals, 6)
}

func (s *DashSuite) TestFetchTxs(c *C) {
	txs, err := s.client.FetchTxs(0)
	c.Assert(err, IsNil)
	c.Assert(txs.Chain, Equals, common.DASHChain)
	c.Assert(txs.Count, Equals, "2")
	c.Assert(txs.TxArray[0].BlockHeight, Equals, int64(1500000))
	c.Assert(txs.TxArray[0].Tx, Equals, "257cf06031508e0c98281dce36b8053025b4c940a9babf129ef09055871e4576")
	c.Assert(txs.TxArray[0].Sender, Equals, "XfzcD6HBSkYzW35Z5yaCnpMkdngp1d1uzA")
	c.Assert(txs.TxArray[0].To, Equals, "7bhSPoNLn5WYUCMrw7AFqN2dN1gSwffYnC")
	c.Assert(txs.TxArray[0].Coins.Equals(common.Coins{common.NewCoin(common.DASHAsset, cosmos.NewUint(101483085))}), Equals, true)
	c.Assert(txs.TxArray[0].Gas.Equals(common.Gas{common.NewCoin(common.DASHAsset, cosmos.NewUint(100000))}), Equals, true)
	c.Assert(len(txs.TxArray), Equals, 2)
}

func (s *DashSuite) TestGetSender(c *C) {
	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "9d340b9f72976b1a85679e0b617d46a442ae654fc964eaeea4928e56841c94cb",
				Vout: 0,
			},
		},
	}
	sender, err := s.client.getSender(&tx)
	c.Assert(err, IsNil)
	c.Assert(sender, Equals, "XtCRm5UTh6fzC2ekSKZv3xRc31VFp6rT3r")

	tx.Vin[0].Vout = 1
	sender, err = s.client.getSender(&tx)
	c.Assert(err, IsNil)
	c.Assert(sender, Equals, "XfzcD6HBSkYzW35Z5yaCnpMkdngp1d1uzA")
}

func (s *DashSuite) TestGetMemo(c *C) {
	tx := btcjson.TxRawResult{
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:       "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Hex:       "",
					ReqSigs:   0,
					Type:      "nulldata",
					Addresses: nil,
				},
			},
		},
	}
	memo, err := s.client.getMemo(&tx)
	c.Assert(err, IsNil)
	c.Assert(memo, Equals, "thorchain:consolidate")

	tx = btcjson.TxRawResult{
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 737761703a6574682e3078633534633135313236393646334541373935366264396144343130383138654563414443466666663a30786335346331353132363936463345413739353662643961443431",
					Type: "nulldata",
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 30383138654563414443466666663a3130303030303030303030",
					Type: "nulldata",
				},
			},
		},
	}
	memo, err = s.client.getMemo(&tx)
	c.Assert(err, IsNil)
	c.Assert(memo, Equals, "swap:eth.0xc54c1512696F3EA7956bd9aD410818eEcADCFfff:0xc54c1512696F3EA7956bd9aD410818eEcADCFfff:10000000000")

	tx = btcjson.TxRawResult{
		Vout: []btcjson.Vout{},
	}
	memo, err = s.client.getMemo(&tx)
	c.Assert(err, IsNil)
	c.Assert(memo, Equals, "")
}

func (s *DashSuite) TestIgnoreTx(c *C) {
	// valid tx that will NOT be ignored
	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.12345678,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:       "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
					Type:      "nulldata",
				},
			},
		},
	}
	ignored, _ := s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, false)

	// invalid tx missing Vout
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, true)

	// invalid tx missing vout[0].Value == no coins
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, true)

	// invalid tx missing vin[0].Txid means coinbase
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, true)

	// invalid tx missing vin
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, true)

	// invalid tx multiple vout[0].Addresses
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, true)

	// invalid tx > 2 vout with coins we only expect 2 max
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, true)

	// valid tx == 2 vout with coins, 1 to vault, 1 with change back to user
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, false)
	// memo at first output should not ignore
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, false)

	// memo in the middle , should not ignore
	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
					},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 0.1234565,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{
						"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6",
					},
				},
			},
		},
	}
	ignored, _ = s.client.ignoreTx(&tx)
	c.Assert(ignored, Equals, false)
}

func (s *DashSuite) TestGetGas(c *C) {
	// vin[0] returns value 0.05020000
	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5da79b6f016521a92b6ed45803f99950649efdeb495fd27e3e7d309ac611336c",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.01234567,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"XuhjXrjKBJumufDkPFP9heunnMrgXoHoqx"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm: "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
				},
			},
		},
	}
	gas, err := s.client.getGas(&tx)
	c.Assert(err, IsNil)
	c.Assert(gas.Equals(common.Gas{common.NewCoin(common.DASHAsset, cosmos.NewUint(3785433))}), Equals, true)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5da79b6f016521a92b6ed45803f99950649efdeb495fd27e3e7d309ac611336c",
				Vout: 0,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				Value: 0.00496556,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm: "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
				},
			},
		},
	}
	gas, err = s.client.getGas(&tx)
	c.Assert(err, IsNil)
	c.Assert(gas.Equals(common.Gas{common.NewCoin(common.DASHAsset, cosmos.NewUint(4328060))}), Equals, true)
}

func (s *DashSuite) TestGetChain(c *C) {
	chain := s.client.GetChain()
	c.Assert(chain, Equals, common.DASHChain)
}

func (s *DashSuite) TestGetAddress(c *C) {
	os.Setenv("NET", "mainnet")
	pubkey := common.PubKey("tmayapub1addwnpepqt7qug8vk9r3saw8n4r803ydj2g3dqwx0mvq5akhnze86fc536xcy7cau6l")
	addr := s.client.GetAddress(pubkey)
	c.Assert(addr, Equals, "XiBCCabtPW9w8EfADdZJZiAmYYaNaZQDsf")
}

func (s *DashSuite) TestGetHeight(c *C) {
	height, err := s.client.GetHeight()
	c.Assert(err, IsNil)
	c.Assert(height, Equals, int64(1906205))
}

func (s *DashSuite) TestGetAccount(c *C) {
	acct, err := s.client.GetAccount("tmayapub1addwnpepqt7qug8vk9r3saw8n4r803ydj2g3dqwx0mvq5akhnze86fc536xcy7cau6l", nil)
	c.Assert(err, IsNil)
	c.Assert(acct.AccountNumber, Equals, int64(0))
	c.Assert(acct.Sequence, Equals, int64(0))
	c.Assert(acct.Coins[0].Amount.Uint64(), Equals, uint64(10674375))

	acct1, err := s.client.GetAccount("", nil)
	c.Assert(err, NotNil)
	c.Assert(acct1.AccountNumber, Equals, int64(0))
	c.Assert(acct1.Sequence, Equals, int64(0))
	c.Assert(acct1.Coins, HasLen, 0)
}

func (s *DashSuite) TestOnObservedTxIn(c *C) {
	pkey := ttypes.GetRandomPubKey()
	txIn := types.TxIn{
		Count: "1",
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 1,
				Tx:          "31f8699ce9028e9cd37f8a6d58a79e614a96e3fdd0f58be5fc36d2d95484716f",
				Sender:      "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				To:          "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(123456789)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
	}
	blockMeta := utxo.NewBlockMeta("000000001ab8a8484eb89f04b87d90eb88e2cbb2829e84eb36b966dcb28af90b", 1, "00000000ffa57c95f4f226f751114e9b24fdf8dbe2dbc02a860da9320bebd63e")
	c.Assert(s.client.temporalStorage.SaveBlockMeta(blockMeta.Height, blockMeta), IsNil)
	s.client.OnObservedTxIn(txIn.TxArray[0], 1)
	blockMeta, err := s.client.temporalStorage.GetBlockMeta(1)
	c.Assert(err, IsNil)
	c.Assert(blockMeta, NotNil)

	txIn = types.TxIn{
		Count: "1",
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(123456)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
	}
	blockMeta = utxo.NewBlockMeta("000000001ab8a8484eb89f04b87d90eb88e2cbb2829e84eb36b966dcb28af90b", 2, "00000000ffa57c95f4f226f751114e9b24fdf8dbe2dbc02a860da9320bebd63e")
	c.Assert(s.client.temporalStorage.SaveBlockMeta(blockMeta.Height, blockMeta), IsNil)
	s.client.OnObservedTxIn(txIn.TxArray[0], 2)
	blockMeta, err = s.client.temporalStorage.GetBlockMeta(2)
	c.Assert(err, IsNil)
	c.Assert(blockMeta, NotNil)

	txIn = types.TxIn{
		Count: "2",
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 3,
				Tx:          "44ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(12345678)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
			{
				BlockHeight: 3,
				Tx:          "54ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(123456)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
	}
	blockMeta = utxo.NewBlockMeta("000000001ab8a8484eb89f04b87d90eb88e2cbb2829e84eb36b966dcb28af90b", 3, "00000000ffa57c95f4f226f751114e9b24fdf8dbe2dbc02a860da9320bebd63e")
	c.Assert(s.client.temporalStorage.SaveBlockMeta(blockMeta.Height, blockMeta), IsNil)
	for _, item := range txIn.TxArray {
		s.client.OnObservedTxIn(item, 3)
	}

	blockMeta, err = s.client.temporalStorage.GetBlockMeta(3)
	c.Assert(err, IsNil)
	c.Assert(blockMeta, NotNil)
}

func (s *DashSuite) TestConfirmationCountReady(c *C) {
	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain:    common.DASHChain,
		TxArray:  nil,
		Filtered: true,
		MemPool:  false,
	}), Equals, true)
	pkey := ttypes.GetRandomPubKey()
	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(123456)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
		Filtered: true,
		MemPool:  true,
	}), Equals, true)
	s.client.currentBlockHeight.Store(3)
	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(123456)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, true)
	s.client.currentBlockHeight.Store(7)
	c.Assert(s.client.ConfirmationCountReady(types.TxIn{
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(12345600000)),
				},
				Gas: common.Gas{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(40000)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 5,
	}), Equals, true)
}

func (s *DashSuite) TestGetConfirmationCount(c *C) {
	pkey := ttypes.GetRandomPubKey()
	// no tx in item , confirmation count should be 0
	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain:   common.DASHChain,
		TxArray: nil,
	}), Equals, int64(0))
	// mempool txin , confirmation count should be 0
	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(123456)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
		Filtered:             true,
		MemPool:              true,
		ConfirmationRequired: 0,
	}), Equals, int64(0))

	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(123456)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, int64(0))

	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(12345600)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 0,
	}), Equals, int64(0))

	c.Assert(s.client.GetConfirmationCount(types.TxIn{
		Chain: common.DASHChain,
		TxArray: []types.TxInItem{
			{
				BlockHeight: 2,
				Tx:          "24ed2d26fd5d4e0e8fa86633e40faf1bdfc8d1903b1cd02855286312d48818a2",
				Sender:      "bc1q0s4mg25tu6termrk8egltfyme4q7sg3h0e56p3",
				To:          "bc1q2gjc0rnhy4nrxvuklk6ptwkcs9kcr59mcl2q9j",
				Coins: common.Coins{
					common.NewCoin(common.DASHAsset, cosmos.NewUint(1234560000)),
				},
				Memo:                "MEMO",
				ObservedVaultPubKey: pkey,
			},
		},
		Filtered:             true,
		MemPool:              false,
		ConfirmationRequired: 7,
	}), Equals, int64(0))
}

func (s *DashSuite) TestGetOutput(c *C) {
	tx := btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	out, err := s.client.getOutput("tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6", &tx, false)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
		},
	}
	out, err = s.client.getOutput("tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6", &tx, false)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
		},
	}
	out, err = s.client.getOutput("tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6", &tx, false)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835"},
				},
			},
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	out, err = s.client.getOutput("tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6", &tx, false)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qj08ys4ct2hzzc2hcz6h2hgrvlmsjynaw43s835")
	c.Assert(out.Value, Equals, 1.49655603)

	tx = btcjson.TxRawResult{
		Vin: []btcjson.Vin{
			{
				Txid: "5b0876dcc027d2f0c671fc250460ee388df39697c3ff082007b6ddd9cb9a7513",
				Vout: 1,
			},
		},
		Vout: []btcjson.Vout{
			{
				Value: 1.49655603,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				Value: 0.00195384,
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Addresses: []string{"tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6"},
				},
			},
			{
				ScriptPubKey: btcjson.ScriptPubKeyResult{
					Asm:  "OP_RETURN 74686f72636861696e3a636f6e736f6c6964617465",
					Type: "nulldata",
				},
			},
		},
	}
	out, err = s.client.getOutput("tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6", &tx, true)
	c.Assert(err, IsNil)
	c.Assert(out.ScriptPubKey.Addresses[0], Equals, "tb1qkq7weysjn6ljc2ywmjmwp8ttcckg8yyxjdz5k6")
	c.Assert(out.Value, Equals, 1.49655603)
}
