package mayachain

import (
	"encoding/json"
	"strconv"

	"github.com/blang/semver"

	abci "github.com/tendermint/tendermint/abci/types"
	. "gopkg.in/check.v1"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	ckeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	types2 "github.com/cosmos/cosmos-sdk/types"

	"gitlab.com/mayachain/mayanode/cmd"
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
	"gitlab.com/mayachain/mayanode/x/mayachain/keeper"
	"gitlab.com/mayachain/mayanode/x/mayachain/query"
	"gitlab.com/mayachain/mayanode/x/mayachain/types"
)

type QuerierSuite struct {
	kb      cosmos.KeybaseStore
	mgr     *Mgrs
	k       keeper.Keeper
	querier cosmos.Querier
	ctx     cosmos.Context
}

var _ = Suite(&QuerierSuite{})

type TestQuerierKeeper struct {
	keeper.KVStoreDummy
	txOut *TxOut
}

func (k *TestQuerierKeeper) GetTxOut(_ cosmos.Context, _ int64) (*TxOut, error) {
	return k.txOut, nil
}

func (s *QuerierSuite) SetUpTest(c *C) {
	kb := ckeys.NewInMemory()
	username := "mayachain"
	password := "password"

	_, _, err := kb.NewMnemonic(username, ckeys.English, cmd.BASEChainHDPath, password, hd.Secp256k1)
	c.Assert(err, IsNil)
	s.kb = cosmos.KeybaseStore{
		SignerName:   username,
		SignerPasswd: password,
		Keybase:      kb,
	}
	s.ctx, s.mgr = setupManagerForTest(c)
	s.k = s.mgr.Keeper()
	s.querier = NewQuerier(s.mgr, s.kb)
}

func (s *QuerierSuite) TestQueryKeysign(c *C) {
	ctx, _ := setupKeeperForTest(c)
	ctx = ctx.WithBlockHeight(12)

	pk := GetRandomPubKey()
	toAddr := GetRandomBNBAddress()
	txOut := NewTxOut(1)
	txOutItem := TxOutItem{
		Chain:       common.BNBChain,
		VaultPubKey: pk,
		ToAddress:   toAddr,
		InHash:      GetRandomTxHash(),
		Coin:        common.NewCoin(common.BNBAsset, cosmos.NewUint(100*common.One)),
	}
	txOut.TxArray = append(txOut.TxArray, txOutItem)
	keeper := &TestQuerierKeeper{
		txOut: txOut,
	}

	_, mgr := setupManagerForTest(c)
	mgr.K = keeper
	querier := NewQuerier(mgr, s.kb)

	path := []string{
		"keysign",
		"5",
		pk.String(),
	}
	res, err := querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)
	c.Assert(res, NotNil)
}

func (s *QuerierSuite) TestQueryPool(c *C) {
	ctx, mgr := setupManagerForTest(c)
	querier := NewQuerier(mgr, s.kb)
	path := []string{"pools"}

	pubKey := GetRandomPubKey()
	asgard := NewVault(ctx.BlockHeight(), ActiveVault, AsgardVault, pubKey, common.Chains{common.BNBChain}.Strings(), []ChainContract{})
	c.Assert(mgr.Keeper().SetVault(ctx, asgard), IsNil)

	poolBNB := NewPool()
	poolBNB.Asset = common.BNBAsset
	poolBNB.LPUnits = cosmos.NewUint(100)

	poolBTC := NewPool()
	poolBTC.Asset = common.BTCAsset
	poolBTC.LPUnits = cosmos.NewUint(0)

	err := mgr.Keeper().SetPool(ctx, poolBNB)
	c.Assert(err, IsNil)

	err = mgr.Keeper().SetPool(ctx, poolBTC)
	c.Assert(err, IsNil)

	res, err := querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	var out Pools

	err = json.Unmarshal(res, &out)
	c.Assert(err, IsNil)
	c.Assert(len(out), Equals, 1)

	poolBTC.LPUnits = cosmos.NewUint(100)
	err = mgr.Keeper().SetPool(ctx, poolBTC)
	c.Assert(err, IsNil)

	res, err = querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	err = json.Unmarshal(res, &out)
	c.Assert(err, IsNil)
	c.Assert(len(out), Equals, 2)

	result, err := s.querier(s.ctx, []string{query.QueryPool.Key, "BNB.BNB"}, abci.RequestQuery{})
	c.Assert(result, HasLen, 0)
	c.Assert(err, NotNil)
}

func (s *QuerierSuite) TestVaultss(c *C) {
	ctx, mgr := setupManagerForTest(c)
	querier := NewQuerier(mgr, s.kb)
	path := []string{"pools"}

	pubKey := GetRandomPubKey()
	asgard := NewVault(ctx.BlockHeight(), ActiveVault, AsgardVault, pubKey, common.Chains{common.BNBChain}.Strings(), nil)
	c.Assert(mgr.Keeper().SetVault(ctx, asgard), IsNil)

	poolBNB := NewPool()
	poolBNB.Asset = common.BNBAsset
	poolBNB.LPUnits = cosmos.NewUint(100)

	poolBTC := NewPool()
	poolBTC.Asset = common.BTCAsset
	poolBTC.LPUnits = cosmos.NewUint(0)

	err := mgr.Keeper().SetPool(ctx, poolBNB)
	c.Assert(err, IsNil)

	err = mgr.Keeper().SetPool(ctx, poolBTC)
	c.Assert(err, IsNil)

	res, err := querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	var out Pools
	err = json.Unmarshal(res, &out)
	c.Assert(err, IsNil)
	c.Assert(len(out), Equals, 1)

	poolBTC.LPUnits = cosmos.NewUint(100)
	err = mgr.Keeper().SetPool(ctx, poolBTC)
	c.Assert(err, IsNil)

	res, err = querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	err = json.Unmarshal(res, &out)
	c.Assert(err, IsNil)
	c.Assert(len(out), Equals, 2)

	result, err := s.querier(s.ctx, []string{query.QueryPool.Key, "BNB.BNB"}, abci.RequestQuery{})
	c.Assert(result, HasLen, 0)
	c.Assert(err, NotNil)
}

func (s *QuerierSuite) TestBucket(c *C) {
	ctx, mgr := setupManagerForTest(c)
	querier := NewQuerier(mgr, s.kb)
	path := []string{"buckets"}

	poolBNB := NewPool()
	poolBNB.Asset = common.BNBAsset.GetSyntheticAsset()
	poolBNB.LPUnits = cosmos.NewUint(100)

	poolBTC := NewPool()
	poolBTC.Asset = common.BTCAsset
	poolBTC.LPUnits = cosmos.NewUint(1000)

	err := mgr.Keeper().SetPool(ctx, poolBNB)
	c.Assert(err, IsNil)

	err = mgr.Keeper().SetPool(ctx, poolBTC)
	c.Assert(err, IsNil)

	res, err := querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	var out Pools
	err = json.Unmarshal(res, &out)
	c.Assert(err, IsNil)
	c.Assert(len(out), Equals, 1)
}

func (s *QuerierSuite) TestQueryNodeAccounts(c *C) {
	ctx, keeper := setupKeeperForTest(c)

	_, mgr := setupManagerForTest(c)
	querier := NewQuerier(mgr, s.kb)
	path := []string{"nodes"}

	nodeAccount := GetRandomValidatorNode(NodeActive)
	bp := NewBondProviders(nodeAccount.NodeAddress)
	acc, err := nodeAccount.BondAddress.AccAddress()
	c.Assert(err, IsNil)
	bp.Providers = append(bp.Providers, NewBondProvider(acc))
	bp.Providers[0].Bonded = true
	SetupLiquidityBondForTest(c, ctx, keeper, common.BNBAsset, nodeAccount.BondAddress, nodeAccount, cosmos.NewUint(1000*common.One))
	c.Assert(keeper.SetBondProviders(ctx, bp), IsNil)
	c.Assert(keeper.SetNodeAccount(ctx, nodeAccount), IsNil)
	vault := GetRandomVault()
	vault.Status = ActiveVault
	vault.BlockHeight = 1
	c.Assert(keeper.SetVault(ctx, vault), IsNil)
	res, err := querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	var out types.NodeAccounts
	err1 := json.Unmarshal(res, &out)
	c.Assert(err1, IsNil)
	c.Assert(len(out), Equals, 1)

	nodeAccount2 := GetRandomValidatorNode(NodeActive)
	bp = NewBondProviders(nodeAccount2.NodeAddress)
	acc, err = nodeAccount2.BondAddress.AccAddress()
	c.Assert(err, IsNil)
	bp.Providers = append(bp.Providers, NewBondProvider(acc))
	bp.Providers[0].Bonded = true
	SetupLiquidityBondForTest(c, ctx, keeper, common.BNBAsset, nodeAccount2.BondAddress, nodeAccount2, cosmos.NewUint(3000*common.One))
	c.Assert(keeper.SetBondProviders(ctx, bp), IsNil)
	c.Assert(keeper.SetNodeAccount(ctx, nodeAccount2), IsNil)

	/* Check Bond-weighted rewards estimation works*/
	var nodeAccountResp []QueryNodeAccount

	// Add bond rewards + set min bond for bond-weighted system
	network, _ := keeper.GetNetwork(ctx)
	network.BondRewardRune = cosmos.NewUint(common.One * 1000)
	c.Assert(keeper.SetNetwork(ctx, network), IsNil)
	keeper.SetMimir(ctx, "MinimumBondInCacao", common.One*100000)

	res, err = querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	err1 = json.Unmarshal(res, &nodeAccountResp)
	c.Assert(err1, IsNil)
	c.Assert(len(nodeAccountResp), Equals, 2)

	for _, node := range nodeAccountResp {
		if node.NodeAddress.Equals(nodeAccount.NodeAddress) {
			// First node has 25% of total bond, gets 25% of rewards
			c.Assert(node.Reward.Uint64(), Equals, cosmos.NewUint(common.One*250).Uint64(), Commentf("expected %s, got %s", cosmos.NewUint(250*common.One), node.Reward))
			continue
		} else if node.NodeAddress.Equals(nodeAccount2.NodeAddress) {
			// Second node has 75% of total bond, gets 75% of rewards
			c.Assert(node.Reward.Uint64(), Equals, cosmos.NewUint(common.One*750).Uint64(), Commentf("expected %s, got %s", cosmos.NewUint(750*common.One), node.Reward))
			continue
		}

		c.Fail()
	}

	/* Check querier only returns nodes with bond */
	c.Assert(keeper.SetNodeAccount(ctx, nodeAccount2), IsNil)

	res, err = querier(ctx, path, abci.RequestQuery{})
	c.Assert(err, IsNil)

	err1 = json.Unmarshal(res, &out)
	c.Assert(err1, IsNil)
	c.Assert(len(out), Equals, 2)
}

func (s *QuerierSuite) TestQuerierRagnarokInProgress(c *C) {
	req := abci.RequestQuery{
		Data:   nil,
		Path:   query.QueryRagnarok.Key,
		Height: s.ctx.BlockHeight(),
		Prove:  false,
	}
	// test ragnarok
	result, err := s.querier(s.ctx, []string{query.QueryRagnarok.Key}, req)
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var ragnarok bool
	c.Assert(json.Unmarshal(result, &ragnarok), IsNil)
	c.Assert(ragnarok, Equals, false)
}

func (s *QuerierSuite) TestQueryLiquidityProviders(c *C) {
	req := abci.RequestQuery{
		Data:   nil,
		Path:   query.QueryLiquidityProviders.Key,
		Height: s.ctx.BlockHeight(),
		Prove:  false,
	}
	// test liquidity providers
	result, err := s.querier(s.ctx, []string{query.QueryLiquidityProviders.Key, "BNB.BNB"}, req)
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	s.k.SetLiquidityProvider(s.ctx, LiquidityProvider{
		Asset:              common.BNBAsset,
		CacaoAddress:       GetRandomBNBAddress(),
		AssetAddress:       GetRandomBNBAddress(),
		LastAddHeight:      1024,
		LastWithdrawHeight: 0,
		Units:              cosmos.NewUint(10),
	})
	result, err = s.querier(s.ctx, []string{query.QueryLiquidityProviders.Key, "BNB.BNB"}, req)
	c.Assert(err, IsNil)
	var lps LiquidityProviders
	c.Assert(json.Unmarshal(result, &lps), IsNil)
	c.Assert(lps, HasLen, 1)
}

func (s *QuerierSuite) TestQueryTxInVoter(c *C) {
	req := abci.RequestQuery{
		Data:   nil,
		Path:   query.QueryTxVoter.Key,
		Height: s.ctx.BlockHeight(),
		Prove:  false,
	}
	tx := GetRandomTx()
	// test getTxInVoter
	result, err := s.querier(s.ctx, []string{query.QueryTxVoter.Key, tx.ID.String()}, req)
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)
	observedTxInVote := NewObservedTxVoter(tx.ID, []ObservedTx{NewObservedTx(tx, s.ctx.BlockHeight(), GetRandomPubKey(), s.ctx.BlockHeight())})
	s.k.SetObservedTxInVoter(s.ctx, observedTxInVote)
	result, err = s.querier(s.ctx, []string{query.QueryTxVoter.Key, tx.ID.String()}, req)
	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	var voter ObservedTxVoter
	c.Assert(json.Unmarshal(result, &voter), IsNil)
	c.Assert(voter.Valid(), IsNil)
}

func (s *QuerierSuite) TestQueryTx(c *C) {
	req := abci.RequestQuery{
		Data:   nil,
		Path:   query.QueryTx.Key,
		Height: s.ctx.BlockHeight(),
		Prove:  false,
	}
	tx := GetRandomTx()
	// test get tx in
	result, err := s.querier(s.ctx, []string{query.QueryTx.Key, tx.ID.String()}, req)
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)
	nodeAccount := GetRandomValidatorNode(NodeActive)
	c.Assert(s.k.SetNodeAccount(s.ctx, nodeAccount), IsNil)
	voter, err := s.k.GetObservedTxInVoter(s.ctx, tx.ID)
	c.Assert(err, IsNil)
	voter.Add(NewObservedTx(tx, s.ctx.BlockHeight(), nodeAccount.PubKeySet.Secp256k1, s.ctx.BlockHeight()), nodeAccount.NodeAddress)
	s.k.SetObservedTxInVoter(s.ctx, voter)
	result, err = s.querier(s.ctx, []string{query.QueryTx.Key, tx.ID.String()}, req)
	c.Assert(err, IsNil)
	var newTx struct {
		ObservedTx    `json:"observed_tx"`
		KeygenMetrics types.TssKeysignMetric `json:"keysign_metric,omitempty"`
	}
	c.Assert(json.Unmarshal(result, &newTx), IsNil)
	c.Assert(newTx.Valid(), IsNil)
}

func (s *QuerierSuite) TestQueryKeyGen(c *C) {
	req := abci.RequestQuery{
		Data:   nil,
		Path:   query.QueryKeygensPubkey.Key,
		Height: s.ctx.BlockHeight(),
		Prove:  false,
	}

	result, err := s.querier(s.ctx, []string{
		query.QueryKeygensPubkey.Key,
		"whatever",
	}, req)

	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryKeygensPubkey.Key,
		"10000",
	}, req)

	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryKeygensPubkey.Key,
		strconv.FormatInt(s.ctx.BlockHeight(), 10),
	}, req)
	c.Assert(err, IsNil)
	c.Assert(result, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryKeygensPubkey.Key,
		strconv.FormatInt(s.ctx.BlockHeight(), 10),
		GetRandomPubKey().String(),
	}, req)
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
}

func (s *QuerierSuite) TestQueryQueue(c *C) {
	result, err := s.querier(s.ctx, []string{
		query.QueryQueue.Key,
		strconv.FormatInt(s.ctx.BlockHeight(), 10),
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var q QueryQueue
	c.Assert(json.Unmarshal(result, &q), IsNil)
}

func (s *QuerierSuite) TestQueryHeights(c *C) {
	result, err := s.querier(s.ctx, []string{
		query.QueryHeights.Key,
		strconv.FormatInt(s.ctx.BlockHeight(), 10),
	}, abci.RequestQuery{})
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryHeights.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var q []QueryResLastBlockHeights
	c.Assert(json.Unmarshal(result, &q), IsNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryHeights.Key,
		"BTC",
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	c.Assert(json.Unmarshal(result, &q), IsNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryChainHeights.Key,
		"BTC",
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	c.Assert(json.Unmarshal(result, &q), IsNil)
}

func (s *QuerierSuite) TestQueryConstantValues(c *C) {
	result, err := s.querier(s.ctx, []string{
		query.QueryConstantValues.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
}

func (s *QuerierSuite) TestQueryMimir(c *C) {
	s.k.SetMimir(s.ctx, "hello", 111)
	result, err := s.querier(s.ctx, []string{
		query.QueryMimirValues.Key,
	}, abci.RequestQuery{})
	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	var m map[string]int64
	c.Assert(json.Unmarshal(result, &m), IsNil)
}

func (s *QuerierSuite) TestQueryBan(c *C) {
	result, err := s.querier(s.ctx, []string{
		query.QueryBan.Key,
	}, abci.RequestQuery{})
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryBan.Key,
		"Whatever",
	}, abci.RequestQuery{})
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryBan.Key,
		GetRandomBech32Addr().String(),
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
}

func (s *QuerierSuite) TestQueryNodeAccount(c *C) {
	result, err := s.querier(s.ctx, []string{
		query.QueryNode.Key,
	}, abci.RequestQuery{})
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryNode.Key,
		"Whatever",
	}, abci.RequestQuery{})
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	na := GetRandomValidatorNode(NodeActive)
	bp := NewBondProviders(na.NodeAddress)
	acc, err := na.BondAddress.AccAddress()
	c.Assert(err, IsNil)
	bp.Providers = append(bp.Providers, NewBondProvider(acc))
	bp.Providers[0].Bonded = true
	SetupLiquidityBondForTest(c, s.ctx, s.k, common.BNBAsset, na.BondAddress, na, cosmos.NewUint(1000*common.One))
	c.Assert(s.k.SetBondProviders(s.ctx, bp), IsNil)
	c.Assert(s.k.SetNodeAccount(s.ctx, na), IsNil)
	vault := GetRandomVault()
	vault.Status = ActiveVault
	vault.BlockHeight = 1
	c.Assert(s.k.SetVault(s.ctx, vault), IsNil)
	result, err = s.querier(s.ctx, []string{
		query.QueryNode.Key,
		na.NodeAddress.String(),
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r QueryNodeAccount
	c.Assert(json.Unmarshal(result, &r), IsNil)

	/* Check bond-weighted rewards estimation works */
	// Add another node with 75% of the bond
	nodeAccount2 := GetRandomValidatorNode(NodeActive)
	bp = NewBondProviders(nodeAccount2.NodeAddress)
	acc, err = nodeAccount2.BondAddress.AccAddress()
	c.Assert(err, IsNil)
	bp.Providers = append(bp.Providers, NewBondProvider(acc))
	bp.Providers[0].Bonded = true
	SetupLiquidityBondForTest(c, s.ctx, s.k, common.BNBAsset, nodeAccount2.BondAddress, nodeAccount2, cosmos.NewUint(3000*common.One))
	c.Assert(s.k.SetBondProviders(s.ctx, bp), IsNil)
	c.Assert(s.k.SetNodeAccount(s.ctx, nodeAccount2), IsNil)

	// Add bond rewards + set min bond for bond-weighted system
	network, _ := s.k.GetNetwork(s.ctx)
	network.BondRewardRune = cosmos.NewUint(common.One * 1000)
	c.Assert(s.k.SetNetwork(s.ctx, network), IsNil)
	s.k.SetMimir(s.ctx, "MinimumBondInCacao", common.One*100000)

	// Get first node
	result, err = s.querier(s.ctx, []string{
		query.QueryNode.Key,
		na.NodeAddress.String(),
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r2 QueryNodeAccount
	c.Assert(json.Unmarshal(result, &r2), IsNil)

	// First node has 25% of bond, should have 25% of the rewards
	c.Assert(r2.Bond.Uint64(), Equals, cosmos.NewUint(common.One*2000).Uint64(), Commentf("expected %s, got %s", cosmos.NewUint(2000*common.One).String(), r2.Bond.String()))
	c.Assert(r2.Reward.Uint64(), Equals, cosmos.NewUint(common.One*250).Uint64(), Commentf("expected %s, got %s", cosmos.NewUint(250*common.One).String(), r2.Reward.String()))

	// Get second node
	result, err = s.querier(s.ctx, []string{
		query.QueryNode.Key,
		nodeAccount2.NodeAddress.String(),
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r3 QueryNodeAccount
	c.Assert(json.Unmarshal(result, &r3), IsNil)

	// Second node has 75% of bond, should have 75% of the rewards
	c.Assert(r3.Bond.Uint64(), Equals, cosmos.NewUint(common.One*6000).Uint64(), Commentf("expected %s, got %s", cosmos.NewUint(6000*common.One).String(), r3.Bond.String()))
	c.Assert(r3.Reward.Uint64(), Equals, cosmos.NewUint(common.One*750).Uint64(), Commentf("expected %s, got %s", cosmos.NewUint(750*common.One).String(), r3.Reward.String()))
}

func (s *QuerierSuite) TestQueryPoolAddresses(c *C) {
	na := GetRandomValidatorNode(NodeActive)
	c.Assert(s.k.SetNodeAccount(s.ctx, na), IsNil)
	result, err := s.querier(s.ctx, []string{
		query.QueryInboundAddresses.Key,
		na.NodeAddress.String(),
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)

	var resp struct {
		Current []struct {
			Chain   common.Chain   `json:"chain"`
			PubKey  common.PubKey  `json:"pub_key"`
			Address common.Address `json:"address"`
			Halted  bool           `json:"halted"`
		} `json:"current"`
	}
	c.Assert(json.Unmarshal(result, &resp), IsNil)
}

func (s *QuerierSuite) TestQueryKeysignArrayPubKey(c *C) {
	na := GetRandomValidatorNode(NodeActive)
	c.Assert(s.k.SetNodeAccount(s.ctx, na), IsNil)
	result, err := s.querier(s.ctx, []string{
		query.QueryKeysignArrayPubkey.Key,
	}, abci.RequestQuery{})
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryKeysignArrayPubkey.Key,
		"asdf",
	}, abci.RequestQuery{})
	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	result, err = s.querier(s.ctx, []string{
		query.QueryKeysignArrayPubkey.Key,
		strconv.FormatInt(s.ctx.BlockHeight(), 10),
	}, abci.RequestQuery{})
	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	var r QueryKeysign
	c.Assert(json.Unmarshal(result, &r), IsNil)
}

func (s *QuerierSuite) TestQueryNetwork(c *C) {
	result, err := s.querier(s.ctx, []string{
		query.QueryNetwork.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r Network
	c.Assert(json.Unmarshal(result, &r), IsNil)
}

func (s *QuerierSuite) TestQueryAsgardVault(c *C) {
	c.Assert(s.k.SetVault(s.ctx, GetRandomVault()), IsNil)
	result, err := s.querier(s.ctx, []string{
		query.QueryVaultsAsgard.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r Vaults
	c.Assert(json.Unmarshal(result, &r), IsNil)
}

func (s *QuerierSuite) TestQueryYggdrasilVault(c *C) {
	vault := GetRandomVault()
	vault.Type = YggdrasilVault
	vault.AddFunds(common.Coins{
		common.NewCoin(common.BNBAsset, cosmos.NewUint(common.One*100)),
	})
	c.Assert(s.k.SetVault(s.ctx, vault), IsNil)
	result, err := s.querier(s.ctx, []string{
		query.QueryVaultsYggdrasil.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r []QueryYggdrasilVaults
	c.Assert(json.Unmarshal(result, &r), IsNil)
}

func (s *QuerierSuite) TestQueryVaultPubKeys(c *C) {
	node := GetRandomValidatorNode(NodeActive)
	c.Assert(s.k.SetNodeAccount(s.ctx, node), IsNil)
	vault := GetRandomVault()
	vault.PubKey = node.PubKeySet.Secp256k1
	vault.Type = YggdrasilVault
	vault.AddFunds(common.Coins{
		common.NewCoin(common.BNBAsset, cosmos.NewUint(common.One*100)),
	})
	vault.Routers = []types.ChainContract{
		{
			Chain:  "ETH",
			Router: "0xE65e9d372F8cAcc7b6dfcd4af6507851Ed31bb44",
		},
	}
	c.Assert(s.k.SetVault(s.ctx, vault), IsNil)
	vault1 := GetRandomVault()
	vault1.Routers = vault.Routers
	c.Assert(s.k.SetVault(s.ctx, vault1), IsNil)
	result, err := s.querier(s.ctx, []string{
		query.QueryVaultPubkeys.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r types.QueryVaultsPubKeys
	c.Assert(json.Unmarshal(result, &r), IsNil)
}

func (s *QuerierSuite) TestQueryBalanceModule(c *C) {
	c.Assert(s.k.SetVault(s.ctx, GetRandomVault()), IsNil)
	result, err := s.querier(s.ctx, []string{
		query.QueryBalanceModule.Key,
		"asgard",
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r struct {
		Name    string            `json:"name"`
		Address cosmos.AccAddress `json:"address"`
		Coins   types2.Coins      `json:"coins"`
	}
	c.Assert(json.Unmarshal(result, &r), IsNil)
}

func (s *QuerierSuite) TestQueryVault(c *C) {
	vault := GetRandomVault()

	// Not enough argument
	result, err := s.querier(s.ctx, []string{
		query.QueryVault.Key,
		"BNB",
	}, abci.RequestQuery{})

	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	c.Assert(s.k.SetVault(s.ctx, vault), IsNil)
	result, err = s.querier(s.ctx, []string{
		query.QueryVault.Key,
		vault.PubKey.String(),
	}, abci.RequestQuery{})
	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	var returnVault Vault
	c.Assert(json.Unmarshal(result, &returnVault), IsNil)
	c.Assert(vault.PubKey.Equals(returnVault.PubKey), Equals, true)
	c.Assert(vault.Type, Equals, returnVault.Type)
	c.Assert(vault.Status, Equals, returnVault.Status)
	c.Assert(vault.BlockHeight, Equals, returnVault.BlockHeight)
}

func (s *QuerierSuite) TestQueryVersion(c *C) {
	result, err := s.querier(s.ctx, []string{
		query.QueryVersion.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	var r types.QueryVersion
	c.Assert(json.Unmarshal(result, &r), IsNil)

	verComputed := s.k.GetLowestActiveVersion(s.ctx)
	c.Assert(r.Current.String(), Equals, verComputed.String(),
		Commentf("query should return same version as computed"))

	// override the version computed in BeginBlock
	s.k.SetVersionWithCtx(s.ctx, semver.MustParse("4.5.6"))

	result, err = s.querier(s.ctx, []string{
		query.QueryVersion.Key,
	}, abci.RequestQuery{})
	c.Assert(result, NotNil)
	c.Assert(err, IsNil)
	c.Assert(json.Unmarshal(result, &r), IsNil)
	c.Assert(r.Current.String(), Equals, "4.5.6",
		Commentf("query should use stored version"))
}

func (s *QuerierSuite) TestQueryLiquidityAuctionTier(c *C) {
	// Not enough argument
	result, err := s.querier(s.ctx, []string{
		query.QueryLiquidityAuctionTier.Key,
		"BNB.BNB",
	}, abci.RequestQuery{})

	c.Assert(result, IsNil)
	c.Assert(err, NotNil)

	// liquidity auction hasn't passed
	address := GetRandomBaseAddress()
	lp := types.LiquidityProvider{
		Asset:                     common.BNBAsset,
		CacaoAddress:              address,
		AssetAddress:              GetRandomBNBAddress(),
		Units:                     cosmos.NewUint(100000),
		WithdrawCounter:           cosmos.NewUint(0),
		LastWithdrawCounterHeight: 0,
	}
	s.k.SetLiquidityProvider(s.ctx, lp)

	result, err = s.querier(s.ctx, []string{
		query.QueryLiquidityAuctionTier.Key,
		common.BNBAsset.String(),
		address.String(),
	}, abci.RequestQuery{})

	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	var returnLATier struct {
		Address                common.Address          `json:"address"`
		Tier                   int64                   `json:"tier"`
		LiquidityProvider      types.LiquidityProvider `json:"liquidity_provider"`
		WithdrawLimitStopBlock int64                   `json:"withdraw_limit_stop_block"`
	}
	c.Assert(json.Unmarshal(result, &returnLATier), IsNil)
	c.Assert(lp.Asset.Equals(common.BNBAsset), Equals, true)
	c.Assert(int64(0), Equals, returnLATier.Tier)
	c.Assert(lp.CacaoAddress, Equals, returnLATier.Address)
	c.Assert(lp.WithdrawCounter.Uint64(), Equals, returnLATier.LiquidityProvider.WithdrawCounter.Uint64())
	c.Assert(lp.LastWithdrawCounterHeight, Equals, returnLATier.LiquidityProvider.LastWithdrawCounterHeight)
	c.Assert(lp.Units.Uint64(), Equals, returnLATier.LiquidityProvider.Units.Uint64())
	c.Assert(returnLATier.WithdrawLimitStopBlock, Equals, int64(0))

	// liquidity auction
	s.k.SetMimir(s.ctx, constants.LiquidityAuction.String(), 20)
	lp = types.LiquidityProvider{
		Asset:                     common.BNBAsset,
		CacaoAddress:              address,
		AssetAddress:              GetRandomBNBAddress(),
		Units:                     cosmos.NewUint(100000),
		WithdrawCounter:           cosmos.NewUint(15),
		LastWithdrawCounterHeight: 19,
	}
	s.k.SetLiquidityProvider(s.ctx, lp)
	c.Assert(s.k.SetLiquidityAuctionTier(s.ctx, address, 1), IsNil)
	s.ctx = s.ctx.WithBlockHeight(21)

	result, err = s.querier(s.ctx, []string{
		query.QueryLiquidityAuctionTier.Key,
		common.BNBAsset.String(),
		address.String(),
	}, abci.RequestQuery{})

	c.Assert(err, IsNil)
	c.Assert(result, NotNil)
	c.Assert(json.Unmarshal(result, &returnLATier), IsNil)
	c.Assert(lp.Asset.Equals(common.BNBAsset), Equals, true)
	c.Assert(int64(1), Equals, returnLATier.Tier)
	c.Assert(lp.CacaoAddress, Equals, returnLATier.Address)
	c.Assert(lp.WithdrawCounter.Uint64(), Equals, returnLATier.LiquidityProvider.WithdrawCounter.Uint64())
	c.Assert(lp.LastWithdrawCounterHeight, Equals, returnLATier.LiquidityProvider.LastWithdrawCounterHeight)
	c.Assert(lp.Units.Uint64(), Equals, returnLATier.LiquidityProvider.Units.Uint64())
	c.Assert(returnLATier.WithdrawLimitStopBlock, Equals, int64(220))
}
