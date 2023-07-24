package mayachain

import (
	"errors"
	"fmt"

	se "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	. "gopkg.in/check.v1"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
	"gitlab.com/mayachain/mayanode/x/mayachain/keeper"
	"gitlab.com/mayachain/mayanode/x/mayachain/types"
)

type HandlerUnBondSuite struct{}

var errReturnYgg = errors.New("returnYgg")

type BlankValidatorManager struct {
	ValidatorDummyMgr
}

func (vm BlankValidatorManager) BeginBlock(_ cosmos.Context, _ constants.ConstantValues, _ []string) error {
	return nil
}

func (vm BlankValidatorManager) EndBlock(_ cosmos.Context, _ Manager) []abci.ValidatorUpdate {
	return nil
}

func (vm BlankValidatorManager) RequestYggReturn(_ cosmos.Context, _ NodeAccount, _ Manager) error {
	return errReturnYgg
}

func (vm BlankValidatorManager) processRagnarok(_ cosmos.Context, _ Manager) error {
	return nil
}

func (vm BlankValidatorManager) NodeAccountPreflightCheck(ctx cosmos.Context, na NodeAccount, constAccessor constants.ConstantValues) (NodeStatus, error) {
	return NodeActive, nil
}

type BlankSlasherManager struct {
	DummySlasher
}

func (d BlankSlasherManager) BeginBlock(ctx cosmos.Context, req abci.RequestBeginBlock, constAccessor constants.ConstantValues) {
}

func (d BlankSlasherManager) HandleDoubleSign(ctx cosmos.Context, addr crypto.Address, infractionHeight int64, constAccessor constants.ConstantValues) error {
	return nil
}

func (d BlankSlasherManager) LackObserving(ctx cosmos.Context, constAccessor constants.ConstantValues) error {
	return nil
}

func (d BlankSlasherManager) LackSigning(ctx cosmos.Context, mgr Manager) error {
	return nil
}

func (d BlankSlasherManager) SlashVault(ctx cosmos.Context, vaultPK common.PubKey, coins common.Coins, mgr Manager) error {
	return nil
}

func (d BlankSlasherManager) SlashVaultToLP(ctx cosmos.Context, vaultPK common.PubKey, coins common.Coins, mgr Manager, subsidize bool) error {
	return nil
}

func (d BlankSlasherManager) SlashNodeAccountLP(ctx cosmos.Context, na NodeAccount, slash cosmos.Uint) (cosmos.Uint, []types.PoolAmt, error) {
	return cosmos.ZeroUint(), []types.PoolAmt{}, nil
}

func (d BlankSlasherManager) IncSlashPoints(ctx cosmos.Context, point int64, addresses ...cosmos.AccAddress) {
	for _, addr := range addresses {
		found := false
		for k := range d.pts {
			if k == addr.String() {
				d.pts[k] += point
				found = true
				break
			}
		}
		if !found {
			d.pts[addr.String()] = point
		}
	}
}

func (d BlankSlasherManager) DecSlashPoints(ctx cosmos.Context, point int64, addresses ...cosmos.AccAddress) {
	for _, addr := range addresses {
		found := false
		for k := range d.pts {
			if k == addr.String() {
				d.pts[k] -= point
				found = true
				break
			}
		}
		if !found {
			d.pts[addr.String()] = -point
		}
	}
}

type TestProviderBond struct {
	BondAddress common.Address
	Bond        cosmos.Uint
}

type TestUnBondKeeper struct {
	keeper.KVStoreDummy
	activeNodeAccount      NodeAccount
	failGetNodeAccount     NodeAccount
	notEmptyNodeAccount    NodeAccount
	jailNodeAccount        NodeAccount
	standbyNodeAccount     NodeAccount
	standbyNodeAccountBond cosmos.Uint
	currentPool            Pool
	lp                     LiquidityProvider
	bp                     BondProviders
	providerBond           []TestProviderBond

	vault Vault
}

func (k *TestUnBondKeeper) GetNodeAccount(_ cosmos.Context, addr cosmos.AccAddress) (NodeAccount, error) {
	if k.standbyNodeAccount.NodeAddress.Equals(addr) {
		return k.standbyNodeAccount, nil
	}
	if k.activeNodeAccount.NodeAddress.Equals(addr) {
		return k.activeNodeAccount, nil
	}
	if k.failGetNodeAccount.NodeAddress.Equals(addr) {
		return NodeAccount{}, fmt.Errorf("you asked for this error")
	}
	if k.notEmptyNodeAccount.NodeAddress.Equals(addr) {
		return k.notEmptyNodeAccount, nil
	}
	if k.jailNodeAccount.NodeAddress.Equals(addr) {
		return k.jailNodeAccount, nil
	}
	return NodeAccount{}, nil
}

func (k *TestUnBondKeeper) GetVault(ctx cosmos.Context, pk common.PubKey) (Vault, error) {
	if k.vault.PubKey.Equals(pk) {
		return k.vault, nil
	}
	return k.KVStoreDummy.GetVault(ctx, pk)
}

func (k *TestUnBondKeeper) VaultExists(ctx cosmos.Context, pkey common.PubKey) bool {
	return k.vault.PubKey.Equals(pkey)
}

func (k *TestUnBondKeeper) GetNodeAccountJail(ctx cosmos.Context, addr cosmos.AccAddress) (Jail, error) {
	if k.jailNodeAccount.NodeAddress.Equals(addr) {
		return Jail{
			NodeAddress:   addr,
			ReleaseHeight: ctx.BlockHeight() + 100,
			Reason:        "bad boy",
		}, nil
	}
	return Jail{}, nil
}

func (k *TestUnBondKeeper) GetAsgardVaultsByStatus(_ cosmos.Context, status VaultStatus) (Vaults, error) {
	if status == k.vault.Status {
		return Vaults{k.vault}, nil
	}
	return nil, nil
}

func (k *TestUnBondKeeper) GetMostSecure(_ cosmos.Context, vaults Vaults, _ int64) Vault {
	if len(vaults) == 0 {
		return Vault{}
	}
	return vaults[0]
}

func (k *TestUnBondKeeper) GetPool(_ cosmos.Context, _ common.Asset) (Pool, error) {
	if k.currentPool.Asset.IsEmpty() {
		return Pool{}, errKaboom
	}

	return k.currentPool, nil
}

func (k *TestUnBondKeeper) CalcNodeLiquidityBond(ctx cosmos.Context, na NodeAccount) (cosmos.Uint, error) {
	return k.standbyNodeAccountBond, nil
}

func (k *TestUnBondKeeper) GetLiquidityProviderByAssets(ctx cosmos.Context, assets common.Assets, addr common.Address) (LiquidityProviders, error) {
	return LiquidityProviders{k.lp}, nil
}

func (k *TestUnBondKeeper) GetLiquidityProvider(ctx cosmos.Context, asset common.Asset, addr common.Address) (LiquidityProvider, error) {
	return k.lp, nil
}

func (k *TestUnBondKeeper) CalcLPLiquidityBond(ctx cosmos.Context, bondAddr common.Address, nodeAddr cosmos.AccAddress) (cosmos.Uint, error) {
	for _, p := range k.providerBond {
		if p.BondAddress == bondAddr {
			// return double the bond, because we want to include both asset and cacao
			return p.Bond.MulUint64(2), nil
		}
	}
	return cosmos.ZeroUint(), fmt.Errorf("BondProvider not found")
}

func (k *TestUnBondKeeper) SetBondProviders(ctx cosmos.Context, bp BondProviders) error {
	k.bp = bp
	return nil
}

func (k *TestUnBondKeeper) GetBondProviders(ctx cosmos.Context, acc cosmos.AccAddress) (BondProviders, error) {
	return k.bp, nil
}

func (k *TestUnBondKeeper) SetVault(ctx cosmos.Context, vault Vault) error {
	k.vault = vault
	return nil
}

func (k *TestUnBondKeeper) SetNodeAccount(ctx cosmos.Context, na NodeAccount) error {
	k.standbyNodeAccount = na
	return nil
}

func (k *TestUnBondKeeper) DeleteVault(ctx cosmos.Context, pk common.PubKey) error {
	return nil
}

var _ = Suite(&HandlerUnBondSuite{})

func (HandlerUnBondSuite) TestUnBondHandler_Run(c *C) {
	ctx, k1 := setupKeeperForTest(c)
	// happy path
	activeNodeAccount := GetRandomValidatorNode(NodeActive)
	standbyNodeAccount := GetRandomValidatorNode(NodeStandby)
	c.Assert(k1.SetNodeAccount(ctx, activeNodeAccount), IsNil)
	c.Assert(k1.SetNodeAccount(ctx, standbyNodeAccount), IsNil)
	vault := NewVault(12, ActiveVault, YggdrasilVault, standbyNodeAccount.PubKeySet.Secp256k1, nil, []ChainContract{})
	c.Assert(k1.SetVault(ctx, vault), IsNil)
	vault = NewVault(12, ActiveVault, AsgardVault, GetRandomPubKey(), nil, []ChainContract{})
	vault.Coins = common.Coins{
		common.NewCoin(common.BaseAsset(), cosmos.NewUint(10000*common.One)),
	}
	c.Assert(k1.SetVault(ctx, vault), IsNil)
	standbyNodeAccount.Reward = cosmos.NewUint(2 * common.One)

	mgr := NewDummyMgrWithKeeper(k1)
	mgr.slasher = BlankSlasherManager{}
	handler := NewUnBondHandler(mgr)

	txIn := common.NewTx(
		GetRandomTxHash(),
		standbyNodeAccount.BondAddress,
		GetRandomBNBAddress(),
		common.Coins{
			common.NewCoin(common.BaseAsset(), cosmos.NewUint(uint64(0))),
		},
		BNBGasFeeSingleton,
		"unbond me please",
	)
	na, _ := k1.GetNodeAccount(ctx, standbyNodeAccount.NodeAddress)
	msg := NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, standbyNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint())
	_, err := handler.Run(ctx, msg)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	c.Assert(na.Reward.Equal(cosmos.ZeroUint()), Equals, true, Commentf("%d", na.Reward.Uint64()))

	bondAddr := GetRandomBaseAddress()
	bondAcc, err := bondAddr.AccAddress()
	c.Assert(err, IsNil)
	k := &TestUnBondKeeper{
		activeNodeAccount:      activeNodeAccount,
		failGetNodeAccount:     GetRandomValidatorNode(NodeActive),
		notEmptyNodeAccount:    standbyNodeAccount,
		jailNodeAccount:        GetRandomValidatorNode(NodeStandby),
		currentPool:            Pool{},
		standbyNodeAccountBond: cosmos.NewUint(100 * common.One),
		lp: LiquidityProvider{
			Asset:        common.BNBAsset,
			CacaoAddress: bondAddr,
			AssetAddress: GetRandomBNBAddress(),
			Units:        cosmos.NewUint(200 * common.One),
			BondedNodes: []LPBondedNode{
				{
					NodeAddress: activeNodeAccount.NodeAddress,
					Units:       cosmos.NewUint(100 * common.One),
				},
				{
					NodeAddress: standbyNodeAccount.NodeAddress,
					Units:       cosmos.NewUint(100 * common.One),
				},
			},
		},
	}
	bp := NewBondProviders(standbyNodeAccount.NodeAddress)
	p := NewBondProvider(bondAcc)
	p.Bonded = true
	bp.Providers = append(bp.Providers, p)
	k.bp = bp
	mgr = NewDummyMgrWithKeeper(k)
	mgr.validatorMgr = BlankValidatorManager{}
	handler = NewUnBondHandler(mgr)

	// simulate fail to get node account
	msg = NewMsgUnBond(txIn, k.failGetNodeAccount.NodeAddress, GetRandomBNBAddress(), nil, activeNodeAccount.NodeAddress, common.BNBAsset, cosmos.NewUint(1))
	_, err = handler.Run(ctx, msg)
	c.Assert(errors.Is(err, errInternal), Equals, true)

	// simulate vault with funds
	k.vault = Vault{
		Type: YggdrasilVault,
		Coins: common.Coins{
			common.NewCoin(common.BNBAsset, cosmos.NewUint(1)),
		},
		PubKey: standbyNodeAccount.PubKeySet.Secp256k1,
	}
	msg = NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, common.Address(standbyNodeAccount.NodeAddress.String()), nil, standbyNodeAccount.NodeAddress, common.BNBAsset, cosmos.NewUint(1))
	_, err = handler.Run(ctx, msg)
	c.Assert(errors.Is(err, errReturnYgg), Equals, true, Commentf("%s", err))

	// simulate fail to get vault
	k.vault = GetRandomVault()
	msg = NewMsgUnBond(txIn, activeNodeAccount.NodeAddress, GetRandomBNBAddress(), nil, activeNodeAccount.NodeAddress, common.BNBAsset, cosmos.OneUint())
	result, err := handler.Run(ctx, msg)
	c.Assert(err, NotNil)
	c.Assert(result, IsNil)

	// simulate vault is not yggdrasil

	k.vault = Vault{
		Type:   AsgardVault,
		PubKey: standbyNodeAccount.PubKeySet.Secp256k1,
	}

	msg = NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, GetRandomBNBAddress(), nil, standbyNodeAccount.NodeAddress, common.BNBAsset, cosmos.OneUint())
	result, err = handler.Run(ctx, msg)
	c.Assert(err, NotNil)
	c.Assert(result, IsNil)

	// simulate jail nodeAccount can't unbound
	msg = NewMsgUnBond(txIn, k.jailNodeAccount.NodeAddress, GetRandomBNBAddress(), nil, k.jailNodeAccount.NodeAddress, common.BNBAsset, cosmos.OneUint())
	result, err = handler.Run(ctx, msg)
	c.Assert(err, NotNil)
	c.Assert(result, IsNil)

	// invalid message should cause error
	result, err = handler.Run(ctx, NewMsgMimir("whatever", 1, GetRandomBech32Addr()))
	c.Assert(err, NotNil)
	c.Assert(result, IsNil)
}

func (HandlerUnBondSuite) TestUnBondHandlerFailValidation(c *C) {
	ctx, k := setupKeeperForTest(c)
	activeNodeAccount := GetRandomValidatorNode(NodeActive)
	providerAddr := GetRandomBech32Addr()
	c.Assert(k.SetNodeAccount(ctx, activeNodeAccount), IsNil)
	handler := NewUnBondHandler(NewDummyMgrWithKeeper(k))
	txIn := common.NewTx(
		GetRandomTxHash(),
		activeNodeAccount.BondAddress,
		GetRandomBNBAddress(),
		common.Coins{
			common.NewCoin(common.BaseAsset(), cosmos.NewUint(uint64(1))),
		},
		BNBGasFeeSingleton,
		"unbond it",
	)
	txInNoTxID := txIn
	txInNoTxID.ID = ""
	testCases := []struct {
		name        string
		msg         *MsgUnBond
		expectedErr error
	}{
		{
			name:        "empty node address",
			msg:         NewMsgUnBond(txIn, cosmos.AccAddress{}, activeNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint()),
			expectedErr: se.ErrInvalidAddress,
		},
		{
			name:        "empty bond address",
			msg:         NewMsgUnBond(txIn, GetRandomValidatorNode(NodeStandby).NodeAddress, common.Address(""), nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint()),
			expectedErr: se.ErrInvalidAddress,
		},
		{
			name:        "empty request hash",
			msg:         NewMsgUnBond(txInNoTxID, GetRandomValidatorNode(NodeStandby).NodeAddress, activeNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint()),
			expectedErr: se.ErrUnknownRequest,
		},
		{
			name:        "empty signer",
			msg:         NewMsgUnBond(txIn, GetRandomValidatorNode(NodeStandby).NodeAddress, activeNodeAccount.BondAddress, nil, cosmos.AccAddress{}, common.EmptyAsset, cosmos.ZeroUint()),
			expectedErr: se.ErrInvalidAddress,
		},
		{
			name:        "account shouldn't be active",
			msg:         NewMsgUnBond(txIn, activeNodeAccount.NodeAddress, activeNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint()),
			expectedErr: se.ErrUnknownRequest,
		},
		{
			name:        "empty asset with non zero amount",
			msg:         NewMsgUnBond(txIn, GetRandomValidatorNode(NodeStandby).NodeAddress, activeNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.NewUint(1)),
			expectedErr: se.ErrUnknownRequest,
		},
		{
			name:        "non empty provider address with asset",
			msg:         NewMsgUnBond(txIn, GetRandomValidatorNode(NodeStandby).NodeAddress, activeNodeAccount.BondAddress, providerAddr, activeNodeAccount.NodeAddress, common.BNBAsset, cosmos.ZeroUint()),
			expectedErr: se.ErrUnknownRequest,
		},
		{
			name:        "request not from original bond address should not be accepted",
			msg:         NewMsgUnBond(GetRandomTx(), GetRandomBech32Addr(), activeNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint()),
			expectedErr: se.ErrUnauthorized,
		},
	}
	for _, item := range testCases {
		c.Log(item.name)
		_, err := handler.Run(ctx, item.msg)

		c.Check(errors.Is(err, item.expectedErr), Equals, true, Commentf("name: %s, %s", item.name, err))
	}
}

func (HandlerUnBondSuite) TestUnBondHanlder_retiringvault(c *C) {
	ctx, k1 := setupKeeperForTest(c)
	// happy path
	activeNodeAccount := GetRandomValidatorNode(NodeActive)
	standbyNodeAccount := GetRandomValidatorNode(NodeStandby)
	c.Assert(k1.SetNodeAccount(ctx, activeNodeAccount), IsNil)
	c.Assert(k1.SetNodeAccount(ctx, standbyNodeAccount), IsNil)
	vault := NewVault(12, ActiveVault, YggdrasilVault, standbyNodeAccount.PubKeySet.Secp256k1, []string{
		common.BNBChain.String(), common.BTCChain.String(), common.ETHChain.String(), common.LTCChain.String(), common.BCHChain.String(),
	}, []ChainContract{})
	c.Assert(k1.SetVault(ctx, vault), IsNil)
	vault = NewVault(12, ActiveVault, AsgardVault, GetRandomPubKey(), nil, []ChainContract{})
	vault.Coins = common.Coins{
		common.NewCoin(common.BaseAsset(), cosmos.NewUint(10000*common.One)),
	}
	c.Assert(k1.SetVault(ctx, vault), IsNil)
	retiringVault := NewVault(12, RetiringVault, AsgardVault, GetRandomPubKey(), []string{
		common.BNBChain.String(), common.BTCChain.String(), common.ETHChain.String(), common.LTCChain.String(), common.BCHChain.String(),
	}, []ChainContract{})
	retiringVault.Membership = []string{
		activeNodeAccount.PubKeySet.Secp256k1.String(),
		standbyNodeAccount.PubKeySet.Secp256k1.String(),
	}
	retiringVault.Coins = common.Coins{
		common.NewCoin(common.BaseAsset(), cosmos.NewUint(10000*common.One)),
	}
	c.Assert(k1.SetVault(ctx, retiringVault), IsNil)
	handler := NewUnBondHandler(NewDummyMgrWithKeeper(k1))
	txIn := common.NewTx(
		GetRandomTxHash(),
		standbyNodeAccount.BondAddress,
		GetRandomBNBAddress(),
		common.Coins{
			common.NewCoin(common.BaseAsset(), cosmos.NewUint(uint64(1))),
		},
		BNBGasFeeSingleton,
		"unbond me please",
	)
	msg := NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, standbyNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint())
	_, err := handler.Run(ctx, msg)
	c.Assert(err, NotNil)
}

func (HandlerUnBondSuite) TestBondProviders_Validate(c *C) {
	ctx, k := setupKeeperForTest(c)
	activeNodeAccount := GetRandomValidatorNode(NodeActive)
	standbyNodeAccount := GetRandomValidatorNode(NodeStandby)
	c.Assert(k.SetNodeAccount(ctx, activeNodeAccount), IsNil)
	c.Assert(k.SetNodeAccount(ctx, standbyNodeAccount), IsNil)
	txIn := GetRandomTx()
	txIn.Coins = common.NewCoins(common.NewCoin(common.BaseAsset(), cosmos.NewUint(100*common.One)))
	handler := NewUnBondHandler(NewDummyMgrWithKeeper(k))

	// happy path
	msg := NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, standbyNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint())
	err := handler.validate(ctx, *msg)
	c.Assert(err, IsNil)

	// cannot unbond an active node
	msg = NewMsgUnBond(txIn, activeNodeAccount.NodeAddress, activeNodeAccount.BondAddress, nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint())
	err = handler.validate(ctx, *msg)
	c.Assert(err, NotNil)

	// test unbonding a bond provider
	bp := NewBondProviders(standbyNodeAccount.NodeAddress)
	p := NewBondProvider(GetRandomBech32Addr())
	bp.Providers = []BondProvider{p}
	c.Assert(k.SetBondProviders(ctx, bp), IsNil)

	msg = NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, common.Address(p.BondAddress.String()), nil, activeNodeAccount.NodeAddress, common.EmptyAsset, cosmos.ZeroUint())
	err = handler.validate(ctx, *msg)
	c.Assert(err, IsNil)
}

func (HandlerUnBondSuite) TestBondProviders_Handler(c *C) {
	ctx, _ := setupKeeperForTest(c)
	standbyNodeAccount := GetRandomValidatorNode(NodeStandby)
	standbyNodeAccount.Reward = cosmos.NewUint(1 * common.One)

	vaultStandby := GetRandomVault()
	vaultStandby.Type = YggdrasilVault
	vaultStandby.PubKey = standbyNodeAccount.PubKeySet.Secp256k1
	txIn := GetRandomTx()
	txIn.Coins = common.NewCoins(common.NewCoin(common.BaseAsset(), cosmos.NewUint(0)))
	runeAddr := GetRandomBaseAddress()
	bnbAddr := GetRandomBNBAddress()
	bp := NewBondProviders(standbyNodeAccount.NodeAddress)
	acc, err := standbyNodeAccount.BondAddress.AccAddress()
	c.Assert(err, IsNil)
	b := NewBondProvider(acc)
	b.Bonded = true
	bp.Providers = append(bp.Providers, b)
	bond := []TestProviderBond{}
	bond = append(bond, TestProviderBond{
		BondAddress: standbyNodeAccount.BondAddress,
		Bond:        cosmos.NewUint(100 * common.One),
	})
	pool := NewPool()
	pool.Asset = common.BNBAsset
	pool.BalanceAsset = bond[0].Bond
	pool.BalanceCacao = bond[0].Bond
	pool.LPUnits = bond[0].Bond
	pool.Status = PoolAvailable
	k := &TestUnBondKeeper{
		standbyNodeAccount:     standbyNodeAccount,
		standbyNodeAccountBond: cosmos.NewUint(100 * common.One),
		providerBond:           bond,
		currentPool:            pool,
		lp: LiquidityProvider{
			Asset:        common.BNBAsset,
			CacaoAddress: runeAddr,
			AssetAddress: bnbAddr,
			Units:        bond[0].Bond,
			BondedNodes: []LPBondedNode{
				{
					NodeAddress: standbyNodeAccount.NodeAddress,
					Units:       bond[0].Bond,
				},
			},
		},
		vault: vaultStandby,
		bp:    bp,
	}

	mgr := NewDummyMgrWithKeeper(k)
	mgr.slasher = BlankSlasherManager{}
	handler := NewUnBondHandler(mgr)

	// // happy path
	c.Check(bp.Get(standbyNodeAccount.NodeAddress).Bonded, Equals, true)
	msg := NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, standbyNodeAccount.BondAddress, nil, standbyNodeAccount.NodeAddress, common.BNBAsset, bond[0].Bond)
	err = handler.handle(ctx, *msg)
	c.Assert(err, IsNil)
	na, _ := handler.mgr.Keeper().GetNodeAccount(ctx, standbyNodeAccount.NodeAddress)
	c.Check(na.Reward.Uint64(), Equals, uint64(0), Commentf("%d", standbyNodeAccount.Reward.Uint64()))
	bp, _ = handler.mgr.Keeper().GetBondProviders(ctx, standbyNodeAccount.NodeAddress)
	c.Check(len(bp.Providers), Equals, 1)
	c.Check(bp.Get(standbyNodeAccount.NodeAddress).Bonded, Equals, false)

	// node operator unbonds/removes bond provider
	k.standbyNodeAccount.Reward = cosmos.NewUint(100 * common.One)
	p := NewBondProvider(GetRandomBech32Addr())
	p.Bonded = true
	bp.Providers = append(bp.Providers, p)
	k.bp = bp
	k.providerBond = append(k.providerBond, TestProviderBond{
		BondAddress: common.Address(p.BondAddress.String()),
		Bond:        cosmos.NewUint(50 * common.One),
	})
	k.lp = LiquidityProvider{
		Asset:        common.BNBAsset,
		CacaoAddress: common.Address(p.BondAddress.String()),
		AssetAddress: GetRandomBNBAddress(),
		Units:        k.providerBond[1].Bond,
		BondedNodes: []LPBondedNode{
			{
				NodeAddress: standbyNodeAccount.NodeAddress,
				Units:       k.providerBond[1].Bond,
			},
		},
	}
	msg = NewMsgUnBond(txIn, k.standbyNodeAccount.NodeAddress, k.standbyNodeAccount.BondAddress, p.BondAddress, k.standbyNodeAccount.NodeAddress, common.BNBAsset, k.providerBond[1].Bond)
	err = handler.handle(ctx, *msg)
	c.Assert(err, IsNil)
	na, _ = handler.mgr.Keeper().GetNodeAccount(ctx, standbyNodeAccount.NodeAddress)
	c.Check(na.Reward.Uint64(), Equals, uint64(0), Commentf("expected %d got %d", uint64(0), na.Reward.Uint64()))
	bp, _ = handler.mgr.Keeper().GetBondProviders(ctx, standbyNodeAccount.NodeAddress)
	c.Check(bp.Has(p.BondAddress), Equals, true)
	c.Check(bp.Get(p.BondAddress).Bonded, Equals, false)

	// bond provider unbonds 75% of their bond themselves
	k.standbyNodeAccount.Reward = cosmos.NewUint(100 * common.One)
	k.standbyNodeAccountBond = cosmos.NewUint(120 * common.One)
	p2 := NewBondProvider(GetRandomBech32Addr())
	p2.Bonded = true
	bp.Providers = append(bp.Providers, p2)
	k.bp = bp
	k.providerBond = append(k.providerBond, TestProviderBond{
		BondAddress: common.Address(p2.BondAddress.String()),
		Bond:        cosmos.NewUint(60 * common.One),
	})
	k.lp = LiquidityProvider{
		Asset:        common.BNBAsset,
		CacaoAddress: common.Address(p2.BondAddress.String()),
		AssetAddress: GetRandomBNBAddress(),
		Units:        k.providerBond[2].Bond,
		BondedNodes: []LPBondedNode{
			{
				NodeAddress: standbyNodeAccount.NodeAddress,
				Units:       k.providerBond[2].Bond,
			},
		},
	}
	msg = NewMsgUnBond(txIn, standbyNodeAccount.NodeAddress, common.Address(p2.BondAddress.String()), nil, standbyNodeAccount.NodeAddress, common.BNBAsset, cosmos.NewUint(45*common.One))
	err = handler.handle(ctx, *msg)
	c.Assert(err, IsNil)
	na, _ = handler.mgr.Keeper().GetNodeAccount(ctx, standbyNodeAccount.NodeAddress)
	c.Check(na.Reward.Uint64(), Equals, cosmos.NewUint(25*common.One).Uint64(), Commentf("%d", na.Reward.Uint64()))
	bp, _ = handler.mgr.Keeper().GetBondProviders(ctx, standbyNodeAccount.NodeAddress)
	c.Check(bp.Has(p2.BondAddress), Equals, true)
	c.Check(bp.Get(p2.BondAddress).Bonded, Equals, true)
}
