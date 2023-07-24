package mayachain

import (
	"fmt"

	"github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
)

func (h AddLiquidityHandler) addLiquidityV96(ctx cosmos.Context,
	asset common.Asset,
	addCacaoAmount, addAssetAmount cosmos.Uint,
	runeAddr, assetAddr common.Address,
	tx common.Tx,
	stage bool,
	constAccessor constants.ConstantValues,
	tier int64,
) error {
	ctx.Logger().Info("liquidity provision", "asset", asset, "rune amount", addCacaoAmount, "asset amount", addAssetAmount)
	if err := h.validateAddLiquidityMessage(ctx, h.mgr.Keeper(), asset, tx, runeAddr, assetAddr); err != nil {
		return fmt.Errorf("add liquidity message fail validation: %w", err)
	}

	pool, err := h.mgr.Keeper().GetPool(ctx, asset)
	if err != nil {
		return ErrInternal(err, fmt.Sprintf("fail to get pool(%s)", asset))
	}
	synthSupply := h.mgr.Keeper().GetTotalSupply(ctx, pool.Asset.GetSyntheticAsset())
	originalUnits := pool.CalcUnits(h.mgr.GetVersion(), synthSupply)

	// if THORNode have no balance, set the default pool status
	if originalUnits.IsZero() {
		defaultPoolStatus := PoolAvailable.String()
		// if the pools is for gas asset on the chain, automatically enable it
		if !pool.Asset.Equals(pool.Asset.GetChain().GetGasAsset()) && !isLiquidityAuction(ctx, h.mgr.Keeper()) {
			defaultPoolStatus = constAccessor.GetStringValue(constants.DefaultPoolStatus)
		}
		pool.Status = GetPoolStatus(defaultPoolStatus)
	}

	fetchAddr := runeAddr
	if fetchAddr.IsEmpty() {
		fetchAddr = assetAddr
	}
	su, err := h.mgr.Keeper().GetLiquidityProvider(ctx, asset, fetchAddr)
	if err != nil {
		return ErrInternal(err, "fail to get liquidity provider")
	}

	su.LastAddHeight = ctx.BlockHeight()
	if su.Units.IsZero() {
		if su.PendingTxID.IsEmpty() {
			if su.CacaoAddress.IsEmpty() {
				su.CacaoAddress = runeAddr
			}
			if su.AssetAddress.IsEmpty() {
				su.AssetAddress = assetAddr
			}
		}

		if asset.IsVaultAsset() {
			// new SU, by default, places the thor address to the rune address,
			// but here we want it to be on the asset address only
			su.AssetAddress = assetAddr
			su.CacaoAddress = common.NoAddress // no rune to add/withdraw
		} else {
			// ensure input addresses match LP position addresses
			if !runeAddr.Equals(su.CacaoAddress) {
				return errAddLiquidityMismatchAddr
			}
			if !assetAddr.Equals(su.AssetAddress) {
				return errAddLiquidityMismatchAddr
			}
		}
	}

	if !assetAddr.IsEmpty() && !su.AssetAddress.Equals(assetAddr) && !asset.IsVaultAsset() {
		// mismatch of asset addresses from what is known to the address
		// given. Refund it.
		return errAddLiquidityMismatchAddr
	}

	// get tx hashes
	cacaoTxID := tx.ID
	assetTxID := tx.ID
	if addCacaoAmount.IsZero() {
		cacaoTxID = su.PendingTxID
	} else {
		assetTxID = su.PendingTxID
	}

	pendingCacaoAmt := su.PendingCacao.Add(addCacaoAmount)
	pendingAssetAmt := su.PendingAsset.Add(addAssetAmount)

	// if we have an asset address and no asset amount, put the rune pending
	if stage && pendingAssetAmt.IsZero() {
		pool.PendingInboundCacao = pool.PendingInboundCacao.Add(addCacaoAmount)
		su.PendingCacao = pendingCacaoAmt
		su.PendingTxID = tx.ID
		h.mgr.Keeper().SetLiquidityProvider(ctx, su)
		if err := h.mgr.Keeper().SetPool(ctx, pool); err != nil {
			ctx.Logger().Error("fail to save pool pending inbound rune", "error", err)
		}

		// add pending liquidity event
		evt := NewEventPendingLiquidity(pool.Asset, AddPendingLiquidity, su.CacaoAddress, addCacaoAmount, su.AssetAddress, cosmos.ZeroUint(), tx.ID, common.TxID(""))
		if err := h.mgr.EventMgr().EmitEvent(ctx, evt); err != nil {
			return ErrInternal(err, "fail to emit partial add liquidity event")
		}
		return nil
	}

	// if we have a rune address and no rune asset, put the asset in pending
	if stage && pendingCacaoAmt.IsZero() {
		pool.PendingInboundAsset = pool.PendingInboundAsset.Add(addAssetAmount)
		su.PendingAsset = pendingAssetAmt
		su.PendingTxID = tx.ID
		h.mgr.Keeper().SetLiquidityProvider(ctx, su)
		if err := h.mgr.Keeper().SetPool(ctx, pool); err != nil {
			ctx.Logger().Error("fail to save pool pending inbound asset", "error", err)
		}

		// Set Liquidity Auction Tier
		if isLiquidityAuction(ctx, h.mgr.Keeper()) {
			tier1 := constAccessor.GetInt64Value(constants.WithdrawTier1)
			tier3 := constAccessor.GetInt64Value(constants.WithdrawTier3)
			oldTier, err := h.mgr.Keeper().GetLiquidityAuctionTier(ctx, su.CacaoAddress)
			if err != nil {
				return ErrInternal(err, "fail to get liquidity auction tier")
			}

			if oldTier != 0 && tier > oldTier {
				tier = oldTier
			} else if tier < tier1 || tier > tier3 {
				tier = tier3
			}

			err = h.mgr.Keeper().SetLiquidityAuctionTier(ctx, runeAddr, tier)
			if err != nil {
				ctx.Logger().Error("fail to set liquidity auction tier", "error", err)
			}
		}

		evt := NewEventPendingLiquidity(pool.Asset, AddPendingLiquidity, su.CacaoAddress, cosmos.ZeroUint(), su.AssetAddress, addAssetAmount, common.TxID(""), tx.ID)
		if err := h.mgr.EventMgr().EmitEvent(ctx, evt); err != nil {
			return ErrInternal(err, "fail to emit partial add liquidity event")
		}
		return nil
	}

	pool.PendingInboundCacao = common.SafeSub(pool.PendingInboundCacao, su.PendingCacao)
	pool.PendingInboundAsset = common.SafeSub(pool.PendingInboundAsset, su.PendingAsset)
	su.PendingAsset = cosmos.ZeroUint()
	su.PendingCacao = cosmos.ZeroUint()
	su.PendingTxID = ""

	ctx.Logger().Info("pre add liquidity", "pool", pool.Asset, "rune", pool.BalanceCacao, "asset", pool.BalanceAsset, "LP units", pool.LPUnits, "synth units", pool.SynthUnits)
	ctx.Logger().Info("adding liquidity", "rune", addCacaoAmount, "asset", addAssetAmount)

	balanceCacao := pool.BalanceCacao
	balanceAsset := pool.BalanceAsset

	oldPoolUnits := pool.GetPoolUnits()
	var newPoolUnits, liquidityUnits cosmos.Uint
	if asset.IsVaultAsset() {
		pendingCacaoAmt = cosmos.ZeroUint() // sanity check
		newPoolUnits, liquidityUnits = calculateVaultUnitsV1(oldPoolUnits, balanceAsset, pendingAssetAmt)
	} else {
		newPoolUnits, liquidityUnits, err = calculatePoolUnitsV1(oldPoolUnits, balanceCacao, balanceAsset, pendingCacaoAmt, pendingAssetAmt)
		if err != nil {
			return ErrInternal(err, "fail to calculate pool unit")
		}
	}
	ctx.Logger().Info("current pool status", "pool units", newPoolUnits, "liquidity units", liquidityUnits)
	poolRune := balanceCacao.Add(pendingCacaoAmt)
	poolAsset := balanceAsset.Add(pendingAssetAmt)
	pool.LPUnits = pool.LPUnits.Add(liquidityUnits)
	pool.BalanceCacao = poolRune
	pool.BalanceAsset = poolAsset
	ctx.Logger().Info("post add liquidity", "pool", pool.Asset, "rune", pool.BalanceCacao, "asset", pool.BalanceAsset, "LP units", pool.LPUnits, "synth units", pool.SynthUnits, "add liquidity units", liquidityUnits)
	if (pool.BalanceCacao.IsZero() && !asset.IsVaultAsset()) || pool.BalanceAsset.IsZero() {
		return ErrInternal(err, "pool cannot have zero rune or asset balance")
	}
	if err := h.mgr.Keeper().SetPool(ctx, pool); err != nil {
		return ErrInternal(err, "fail to save pool")
	}
	if originalUnits.IsZero() && !pool.GetPoolUnits().IsZero() {
		poolEvent := NewEventPool(pool.Asset, pool.Status)
		if err := h.mgr.EventMgr().EmitEvent(ctx, poolEvent); err != nil {
			ctx.Logger().Error("fail to emit pool event", "error", err)
		}
	}

	su.Units = su.Units.Add(liquidityUnits)
	if pool.Status == PoolAvailable {
		if su.AssetDepositValue.IsZero() && su.CacaoDepositValue.IsZero() {
			su.CacaoDepositValue = common.GetSafeShare(su.Units, pool.GetPoolUnits(), pool.BalanceCacao)
			su.AssetDepositValue = common.GetSafeShare(su.Units, pool.GetPoolUnits(), pool.BalanceAsset)
		} else {
			su.CacaoDepositValue = su.CacaoDepositValue.Add(common.GetSafeShare(liquidityUnits, pool.GetPoolUnits(), pool.BalanceCacao))
			su.AssetDepositValue = su.AssetDepositValue.Add(common.GetSafeShare(liquidityUnits, pool.GetPoolUnits(), pool.BalanceAsset))
		}
		// Recalculate bond provided to node if lp is bonder
		liquidityPools := GetLiquidityPools(h.mgr.GetVersion())
		if liquidityPools.Contains(pool.Asset) && !su.NodeBondAddress.Empty() {
			na, err := h.mgr.Keeper().GetNodeAccount(ctx, su.NodeBondAddress)
			if err != nil {
				ctx.Logger().Error("fail to get bonded node account of LP %s", su.CacaoAddress)
			}

			bp, err := h.mgr.Keeper().GetBondProviders(ctx, su.NodeBondAddress)
			if err != nil {
				ctx.Logger().Error("fail to get bonded providers for node account of LP %s", su.CacaoAddress)
			}

			addedUnits := su.Units.Sub(originalUnits)
			addedBond := common.GetSafeShare(addedUnits, pool.LPUnits, pool.BalanceCacao)
			bondEvent := NewEventBond(addedBond, BondPaid, tx)

			from, err := su.CacaoAddress.AccAddress()
			if err != nil {
				ctx.Logger().Error("fail to get lp account", "error", err)
			}

			if bp.Has(from) {
				bp.BondLiquidity(from)
			}

			if err := h.mgr.EventMgr().EmitEvent(ctx, bondEvent); err != nil {
				ctx.Logger().Error("fail to emit bond event", "error", err)
			}

			if err := h.mgr.Keeper().SetNodeAccount(ctx, na); err != nil {
				return ErrInternal(err, fmt.Sprintf("fail to save node account(%s)", na.String()))
			}

			if err := h.mgr.Keeper().SetBondProviders(ctx, bp); err != nil {
				return ErrInternal(err, fmt.Sprintf("fail to save bond providers(%s)", bp.NodeAddress.String()))
			}
		}
	}

	h.mgr.Keeper().SetLiquidityProvider(ctx, su)

	evt := NewEventAddLiquidity(asset, liquidityUnits, su.CacaoAddress, pendingCacaoAmt, pendingAssetAmt, cacaoTxID, assetTxID, su.AssetAddress)
	if err := h.mgr.EventMgr().EmitEvent(ctx, evt); err != nil {
		return ErrInternal(err, "fail to emit add liquidity event")
	}

	// if its the POL is adding, track rune added
	polAddress, err := h.mgr.Keeper().GetModuleAddress(ReserveName)
	if err != nil {
		return err
	}

	if polAddress.Equals(su.CacaoAddress) {
		pol, err := h.mgr.Keeper().GetPOL(ctx)
		if err != nil {
			return err
		}
		pol.CacaoDeposited = pol.CacaoDeposited.Add(pendingCacaoAmt)

		if err := h.mgr.Keeper().SetPOL(ctx, pol); err != nil {
			return err
		}

		ctx.Logger().Info("POL deposit", "pool", pool.Asset, "rune", pendingCacaoAmt)
		telemetry.IncrCounterWithLabels(
			[]string{"mayanode", "pol", "pool", "rune_deposited"},
			telem(pendingCacaoAmt),
			[]metrics.Label{telemetry.NewLabel("pool", pool.Asset.String())},
		)
	}
	return nil
}
