//go:build stagenet
// +build stagenet

package mayachain

import (
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
)

func importPreRegistrationMAYANames(ctx cosmos.Context, mgr Manager) error {
	oneYear := fetchConfigInt64(ctx, mgr, constants.BlocksPerYear)
	names, err := getPreRegisterMAYANames(ctx, ctx.BlockHeight()+oneYear)
	if err != nil {
		return err
	}

	for _, name := range names {
		mgr.Keeper().SetMAYAName(ctx, name)
	}
	return nil
}

func migrateStoreV96(ctx cosmos.Context, mgr Manager) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v88", "error", err)
		}
	}()

	err := importPreRegistrationMAYANames(ctx, mgr)
	if err != nil {
		ctx.Logger().Error("fail to migrate store to v88", "error", err)
	}
}

func migrateStoreV102(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v102", "error", err)
		}
	}()

	if err := mgr.BeginBlock(ctx); err != nil {
		ctx.Logger().Error("fail to initialise block", "error", err)
		return
	}

	// Remove stuck txs
	hashes := []string{
		"0300C0123B05D5D4C229A0F1BBB51C3D14A328D5B60CDD45303CB121B7051FCF",
		"5C64F88DE868703CE2B5929827E563DDF28EA7BA718B4228BC66EEADE621F9BD",
	}
	removeTransactions(ctx, mgr, hashes...)

	// Rebalance asgard vs real balance for RUNE
	vaultPubkey, err := common.NewPubKey("smayapub1addwnpepqfwjxqah2qsgl97aqjjpn6ghpt54yc8z28t78zgur4uk6t29ccx22gvwrr0")
	if err != nil {
		ctx.Logger().Error("fail to get vault pubkey", "error", err)
		return
	}
	vault, err := mgr.Keeper().GetVault(ctx, vaultPubkey)
	if err != nil {
		ctx.Logger().Error("fail to get vault", "error", err)
		return
	}

	vault.AddFunds(common.NewCoins(common.NewCoin(common.RUNEAsset, cosmos.NewUint(12_73000000))))

	if err := mgr.Keeper().SetVault(ctx, vault); err != nil {
		ctx.Logger().Error("fail to set vault", "error", err)
		return
	}

	// Remove retiring vault balance
	vaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, RetiringVault)
	if err != nil {
		ctx.Logger().Error("fail to get retiring asgard vaults", "error", err)
		return
	}
	for _, v := range vaults {
		runeAsset := v.GetCoin(common.RUNEAsset)
		v.SubFunds(common.NewCoins(runeAsset))
		if err := mgr.Keeper().SetVault(ctx, v); err != nil {
			ctx.Logger().Error("fail to save vault", "error", err)
		}
	}

	// Refund unobserved txs with no memo
	type RefundAccount struct {
		Address string
		TxID    string
		Chain   common.Chain
		Asset   string
		Amount  uint64
	}

	refundAccounts := []RefundAccount{
		{
			Address: "thor18z343fsdlav47chtkyp0aawqt6sgxsh3v96x6j",
			TxID:    "0000000000000000000000000000000000000000000000000000000000000000",
			Amount:  1_00000000,
			Chain:   common.THORChain,
			Asset:   "THOR.RUNE",
		},
	}

	txOutStore := mgr.TxOutStore()
	for _, refundAccount := range refundAccounts {
		addr, err := common.NewAddress(refundAccount.Address)
		if err != nil {
			ctx.Logger().Error("fail to parse address", "error", err)
			return
		}
		token, err := common.NewAsset(refundAccount.Asset)
		if err != nil {
			ctx.Logger().Error("fail to parse asset", "error", err)
			return
		}
		refundTxID, err := common.NewTxID(refundAccount.TxID)
		if err != nil {
			ctx.Logger().Error("fail to parse transaction id", "error", err)
			return
		}

		refund := TxOutItem{
			Chain:     refundAccount.Chain,
			ToAddress: addr,
			Coin:      common.NewCoin(token, cosmos.NewUint(refundAccount.Amount)),
			Memo:      NewRefundMemo(refundTxID).String(),
			InHash:    refundTxID,
		}
		_, err = txOutStore.TryAddTxOutItem(ctx, mgr, refund, cosmos.ZeroUint())
		if err != nil {
			ctx.Logger().Error("fail to schedule refund transaction", "error", err)
			return
		}
	}

	// Add LPs from unobserved txs
	lps := []struct {
		MayaAddress string
		ThorAddress string
		Amount      cosmos.Uint
	}{
		{
			MayaAddress: "smaya18z343fsdlav47chtkyp0aawqt6sgxsh3ctcu6u",
			ThorAddress: "thor18z343fsdlav47chtkyp0aawqt6sgxsh3v96x6j",
			Amount:      cosmos.NewUint(1_00000000),
		},

		{
			MayaAddress: "smaya1x0jkvqdh2hlpeztd5zyyk70n3efx6mhuecefny",
			ThorAddress: "thor1x0jkvqdh2hlpeztd5zyyk70n3efx6mhudkmnn2",
			Amount:      cosmos.NewUint(1),
		},
	}

	for _, sender := range lps {
		address, err := common.NewAddress(sender.MayaAddress)
		if err != nil {
			ctx.Logger().Error("fail to parse address", "error", err)
			continue
		}

		lp, err := mgr.Keeper().GetLiquidityProvider(ctx, common.RUNEAsset, address)
		if err != nil {
			ctx.Logger().Error("fail to get liquidity provider", "error", err)
			continue
		}

		lp.PendingAsset = lp.PendingAsset.Add(sender.Amount)
		lp.LastAddHeight = ctx.BlockHeight()
		lp.PendingTxID = common.BlankTxID

		if lp.AssetAddress.IsEmpty() {
			thorAdd, err := common.NewAddress(sender.ThorAddress)
			if err != nil {
				ctx.Logger().Error("fail to parse address", "error", err)
				continue
			}
			lp.AssetAddress = thorAdd
		}
	}
	// Return cacao to reserve from Itzamna and BTC pool
	// Mint cacao
	toMint := common.NewCoin(common.BaseAsset(), cosmos.NewUint(9_900_000_000_00000000))
	mgr.Keeper().MintToModule(ctx, ModuleName, toMint)
	if err = mgr.Keeper().SendFromModuleToModule(ctx, ModuleName, ReserveName, common.NewCoins(toMint)); err != nil {
		ctx.Logger().Error("fail to send cacao to reserve", "error", err)
		return
	}

	// TO RESERVE TXS
	// 360209731069904 bond to reserve
	// Update network
	network, err := mgr.Keeper().GetNetwork(ctx)
	if err != nil {
		ctx.Logger().Error("fail to get network", "error", err)
		return
	}
	network.BondRewardRune = cosmos.ZeroUint()
	if err := mgr.Keeper().SetNetwork(ctx, network); err != nil {
		ctx.Logger().Error("fail to set network", "error", err)
		return
	}

	bondToReserve := common.NewCoin(common.BaseAsset(), cosmos.NewUint(24000000))
	if err := mgr.Keeper().SendFromModuleToModule(ctx, BondName, ReserveName, common.NewCoins(bondToReserve)); err != nil {
		ctx.Logger().Error("fail to send bond to reserve", "error", err)
		return
	}

	for _, asset := range []common.Asset{common.BTCAsset, common.ETHAsset, common.RUNEAsset} {
		pool, err := mgr.Keeper().GetPool(ctx, asset)
		if err != nil {
			ctx.Logger().Error("fail to get pool", "error", err)
			return
		}
		switch asset {
		case common.BTCAsset:
			pool.BalanceCacao = pool.BalanceCacao.SubUint64(6)
		}
		if err := mgr.Keeper().SetPool(ctx, pool); err != nil {
			ctx.Logger().Error("fail to set pool", "error", err)
			return
		}
	}

	// Sum of all the above will be sent
	asgardToReserve := common.NewCoin(common.BaseAsset(), cosmos.NewUint(6))
	if err := mgr.Keeper().SendFromModuleToModule(ctx, AsgardName, ReserveName, common.NewCoins(asgardToReserve)); err != nil {
		ctx.Logger().Error("fail to send asgard to reserve", "error", err)
		return
	}

	// 164293529917265 de itzamna a reserve
	itzamnaToReserve := common.NewCoin(common.BaseAsset(), cosmos.NewUint(2998_76000000))
	itzamnaAcc, err := cosmos.AccAddressFromBech32("smaya18z343fsdlav47chtkyp0aawqt6sgxsh3ctcu6u")
	if err != nil {
		ctx.Logger().Error("fail to parse address", "error", err)
		return
	}

	if err := mgr.Keeper().SendFromAccountToModule(ctx, itzamnaAcc, ReserveName, common.NewCoins(itzamnaToReserve)); err != nil {
		ctx.Logger().Error("fail to send itzamna to reserve", "error", err)
		return
	}

	// FROM RESERVE TXS
	// 8_910_000_000_00000000 from reserve to itzamna
	// 815_38461482 from reserve to USDT asgard
	reserveToItzamna := common.NewCoin(common.BaseAsset(), cosmos.NewUint(8_910_000_020_00000000))
	if err := mgr.Keeper().SendFromModuleToAccount(ctx, ReserveName, itzamnaAcc, common.NewCoins(reserveToItzamna)); err != nil {
		ctx.Logger().Error("fail to send reserve to itzamna", "error", err)
		return
	}

	usdtPool, err := mgr.Keeper().GetPool(ctx, common.USDTAsset)
	if err != nil {
		ctx.Logger().Error("fail to get pool", "error", err)
		return
	}

	usdtPool.BalanceCacao = usdtPool.BalanceCacao.Add(cosmos.NewUint(12))
	if err := mgr.Keeper().SetPool(ctx, usdtPool); err != nil {
		ctx.Logger().Error("fail to set pool", "error", err)
		return
	}

	reserveToAsgard := common.NewCoin(common.BaseAsset(), cosmos.NewUint(12))
	if err := mgr.Keeper().SendFromModuleToModule(ctx, ReserveName, AsgardName, common.NewCoins(reserveToAsgard)); err != nil {
		ctx.Logger().Error("fail to send reserve to asgard", "error", err)
		return
	}

	// Remove Slash points from genesis nodes
	for _, genesis := range GenesisNodes {
		acc, err := cosmos.AccAddressFromBech32(genesis)
		if err != nil {
			ctx.Logger().Error("fail to parse address", "error", err)
			continue
		}

		mgr.Keeper().ResetNodeAccountSlashPoints(ctx, acc)
	}

	for _, node := range []string{"smaya1g5pgvndmtpejhrnkwem3y5tpznkuhd3cearctd", "smaya1jltkeg0g56jegjwfld90d2g9fnd2kuwpnp265k", "smaya1mapzy6qfswyfjc6f8g08uj30vng74aqqqpethg", "smaya1nay0rxpjl2gk3nw7gmj3at50nc2xq3fnsgrt0w"} {
		acc, err := cosmos.AccAddressFromBech32(node)
		if err != nil {
			ctx.Logger().Error("fail to parse address", "error", err)
			continue
		}

		mgr.Keeper().ResetNodeAccountSlashPoints(ctx, acc)
	}
}

func migrateStoreV104(ctx cosmos.Context, mgr Manager) {}

func migrateStoreV105(ctx cosmos.Context, mgr Manager) {}
