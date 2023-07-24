package mayachain

import (
	"errors"
	"fmt"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
)

func refundTxV47(ctx cosmos.Context, tx ObservedTx, mgr Manager, refundCode uint32, refundReason, nativeRuneModuleName string) error {
	// If THORNode recognize one of the coins, and therefore able to refund
	// withholding fees, refund all coins.

	addEvent := func(refundCoins common.Coins) error {
		eventRefund := NewEventRefund(refundCode, refundReason, tx.Tx, common.NewFee(common.Coins{}, cosmos.ZeroUint()))
		if len(refundCoins) > 0 {
			// create a new TX based on the coins thorchain refund , some of the coins thorchain doesn't refund
			// coin thorchain doesn't have pool with , likely airdrop
			newTx := common.NewTx(tx.Tx.ID, tx.Tx.FromAddress, tx.Tx.ToAddress, tx.Tx.Coins, tx.Tx.Gas, tx.Tx.Memo)

			// all the coins in tx.Tx should belongs to the same chain
			transactionFee := mgr.GasMgr().GetFee(ctx, tx.Tx.Chain, common.BaseAsset())
			fee := getFee(tx.Tx.Coins, refundCoins, transactionFee)
			eventRefund = NewEventRefund(refundCode, refundReason, newTx, fee)
		}
		if err := mgr.EventMgr().EmitEvent(ctx, eventRefund); err != nil {
			return fmt.Errorf("fail to emit refund event: %w", err)
		}
		return nil
	}

	// for BASEChain transactions, create the event before we txout. For other
	// chains, do it after. The reason for this is we need to make sure the
	// first event (refund) is created, before we create the outbound events
	// (second). Because its BASEChain, its safe to assume all the coins are
	// safe to send back. Where as for external coins, we cannot make this
	// assumption (ie coins we don't have pools for and therefore, don't know
	// the value of it relative to rune)
	if tx.Tx.Chain.Equals(common.BASEChain) {
		if err := addEvent(tx.Tx.Coins); err != nil {
			return err
		}
	}
	refundCoins := make(common.Coins, 0)
	for _, coin := range tx.Tx.Coins {
		if coin.Asset.IsBase() && coin.Asset.GetChain().Equals(common.ETHChain) {
			continue
		}
		pool, err := mgr.Keeper().GetPool(ctx, coin.Asset)
		if err != nil {
			return fmt.Errorf("fail to get pool: %w", err)
		}

		if coin.Asset.IsBase() || !pool.BalanceCacao.IsZero() {
			toi := TxOutItem{
				Chain:       coin.Asset.GetChain(),
				InHash:      tx.Tx.ID,
				ToAddress:   tx.Tx.FromAddress,
				VaultPubKey: tx.ObservedPubKey,
				Coin:        coin,
				Memo:        NewRefundMemo(tx.Tx.ID).String(),
				ModuleName:  nativeRuneModuleName,
			}

			success, err := mgr.TxOutStore().TryAddTxOutItem(ctx, mgr, toi, cosmos.ZeroUint())
			if err != nil {
				ctx.Logger().Error("fail to prepare outbund tx", "error", err)
				// concatenate the refund failure to refundReason
				refundReason = fmt.Sprintf("%s; fail to refund (%s): %s", refundReason, toi.Coin.String(), err)
			}
			if success {
				refundCoins = append(refundCoins, toi.Coin)
			}
		}
		// Zombie coins are just dropped.
	}
	if !tx.Tx.Chain.Equals(common.BASEChain) {
		if err := addEvent(refundCoins); err != nil {
			return err
		}
	}

	return nil
}

func refundBondV92(ctx cosmos.Context, tx common.Tx, acc cosmos.AccAddress, nodeAcc *NodeAccount, mgr Manager) error {
	if nodeAcc.Status == NodeActive {
		ctx.Logger().Info("node still active, cannot refund bond", "node address", nodeAcc.NodeAddress, "node pub key", nodeAcc.PubKeySet.Secp256k1)
		return nil
	}

	// ensures nodes don't return bond while being churned into the network
	// (removing their bond last second)
	if nodeAcc.Status == NodeReady {
		ctx.Logger().Info("node ready, cannot refund bond", "node address", nodeAcc.NodeAddress, "node pub key", nodeAcc.PubKeySet.Secp256k1)
		return nil
	}

	nodeBond, err := mgr.Keeper().CalcNodeLiquidityBond(ctx, *nodeAcc)
	if err != nil {
		return ErrInternal(err, fmt.Sprintf("fail to get node liquidity bond (%s)", nodeAcc.NodeAddress))
	}

	ygg := Vault{}
	if mgr.Keeper().VaultExists(ctx, nodeAcc.PubKeySet.Secp256k1) {
		var err error
		ygg, err = mgr.Keeper().GetVault(ctx, nodeAcc.PubKeySet.Secp256k1)
		if err != nil {
			return err
		}
		if !ygg.IsYggdrasil() {
			return errors.New("this is not a Yggdrasil vault")
		}
	}

	bp, err := mgr.Keeper().GetBondProviders(ctx, nodeAcc.NodeAddress)
	if err != nil {
		return ErrInternal(err, fmt.Sprintf("fail to get bond providers(%s)", nodeAcc.NodeAddress))
	}

	// Calculate total value (in rune) the Yggdrasil pool has
	yggRune, err := getTotalYggValueInRune(ctx, mgr.Keeper(), ygg)
	if err != nil {
		return fmt.Errorf("fail to get total ygg value in RUNE: %w", err)
	}

	if nodeBond.LT(yggRune) {
		ctx.Logger().Error("Node Account left with more funds in their Yggdrasil vault than their bond's value", "address", nodeAcc.NodeAddress, "ygg-value", yggRune, "bond", nodeBond)
	}
	// slash yggdrasil remains
	penaltyPts := fetchConfigInt64(ctx, mgr, constants.SlashPenalty)
	slashRune := common.GetUncappedShare(cosmos.NewUint(uint64(penaltyPts)), cosmos.NewUint(10_000), yggRune)
	if slashRune.GT(nodeBond) {
		slashRune = nodeBond
	}

	slashedAmount, _, err := mgr.Slasher().SlashNodeAccountLP(ctx, *nodeAcc, slashRune)
	if err != nil {
		return ErrInternal(err, "fail to slash node account")
	}

	if slashedAmount.LT(slashRune) {
		ctx.Logger().Error("slashed amount is less than slash rune", "slashed amount", slashedAmount, "slash rune", slashRune)
	}

	provider := bp.Get(acc)

	nodeBond, err = mgr.Keeper().CalcNodeLiquidityBond(ctx, *nodeAcc)
	if err != nil {
		return ErrInternal(err, fmt.Sprintf("fail to get node liquidity bond (%s)", nodeAcc.NodeAddress))
	}

	providerBond, err := mgr.Keeper().CalcLPLiquidityBond(ctx, common.Address(provider.BondAddress.String()), nodeAcc.NodeAddress)
	if err != nil {
		return ErrInternal(err, "fail to get bond provider liquidity")
	}

	if !provider.IsEmpty() && !nodeBond.IsZero() && !providerBond.IsZero() {

		bp.Unbond(provider.BondAddress)

		toAddress, err := common.NewAddress(provider.BondAddress.String())
		if err != nil {
			return fmt.Errorf("fail to parse bond address: %w", err)
		}

		// calculate rewards for bond provider
		//  Rewards * (ProviderBond / NodeBond)
		if !nodeAcc.Reward.IsZero() {
			bondRewards := common.GetSafeShare(providerBond, nodeBond, nodeAcc.Reward)
			nodeAcc.Reward = common.SafeSub(nodeAcc.Reward, bondRewards)

			// refund bond
			txOutItem := TxOutItem{
				Chain:      common.BaseAsset().Chain,
				ToAddress:  toAddress,
				InHash:     tx.ID,
				Coin:       common.NewCoin(common.BaseAsset(), bondRewards),
				ModuleName: BondName,
			}
			_, err = mgr.TxOutStore().TryAddTxOutItem(ctx, mgr, txOutItem, cosmos.ZeroUint())
			if err != nil {
				return fmt.Errorf("fail to add outbound tx: %w", err)
			}

			bondEvent := NewEventBond(bondRewards, BondReturned, tx)
			if err := mgr.EventMgr().EmitEvent(ctx, bondEvent); err != nil {
				ctx.Logger().Error("fail to emit bond event", "error", err)
			}
		}
	}

	if nodeAcc.RequestedToLeave {
		// when node already request to leave , it can't come back , here means the node already unbond
		// so set the node to disabled status
		nodeAcc.UpdateStatus(NodeDisabled, ctx.BlockHeight())
	}
	if err := mgr.Keeper().SetNodeAccount(ctx, *nodeAcc); err != nil {
		ctx.Logger().Error(fmt.Sprintf("fail to save node account(%s)", nodeAcc), "error", err)
		return err
	}
	if err := mgr.Keeper().SetBondProviders(ctx, bp); err != nil {
		return ErrInternal(err, fmt.Sprintf("fail to save bond providers(%s)", bp.NodeAddress.String()))
	}

	if err := subsidizePoolsWithSlashBond(ctx, ygg.Coins, ygg, yggRune, slashedAmount, mgr); err != nil {
		ctx.Logger().Error("fail to subsidize pools with slash bond", "error", err)
	}
	// at this point , all coins in yggdrasil vault has been accounted for , and node already been slashed
	ygg.SubFunds(ygg.Coins)
	if err := mgr.Keeper().SetVault(ctx, ygg); err != nil {
		ctx.Logger().Error("fail to save yggdrasil vault", "error", err)
		return err
	}

	if err := mgr.Keeper().DeleteVault(ctx, ygg.PubKey); err != nil {
		return err
	}

	// Output bond events for the slashed and returned bond.
	if !slashRune.IsZero() {
		fakeTx := common.Tx{}
		fakeTx.ID = common.BlankTxID
		fakeTx.FromAddress = nodeAcc.BondAddress
		bondEvent := NewEventBond(slashRune, BondCost, fakeTx)
		if err := mgr.EventMgr().EmitEvent(ctx, bondEvent); err != nil {
			ctx.Logger().Error("fail to emit bond event", "error", err)
		}
	}
	return nil
}
