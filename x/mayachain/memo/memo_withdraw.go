package mayachain

import (
	"fmt"

	"gitlab.com/mayachain/mayanode/common"
	cosmos "gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/x/mayachain/keeper"
	"gitlab.com/mayachain/mayanode/x/mayachain/types"
)

type WithdrawLiquidityMemo struct {
	MemoBase
	Amount          cosmos.Uint
	WithdrawalAsset common.Asset
	PairAddress     common.Address
}

func (m WithdrawLiquidityMemo) GetAmount() cosmos.Uint           { return m.Amount }
func (m WithdrawLiquidityMemo) GetWithdrawalAsset() common.Asset { return m.WithdrawalAsset }
func (m WithdrawLiquidityMemo) GetPairAddress() common.Address   { return m.PairAddress }

func NewWithdrawLiquidityMemo(asset common.Asset, amt cosmos.Uint, withdrawalAsset common.Asset, pairAddress common.Address) WithdrawLiquidityMemo {
	return WithdrawLiquidityMemo{
		MemoBase:        MemoBase{TxType: TxWithdraw, Asset: asset},
		Amount:          amt,
		WithdrawalAsset: withdrawalAsset,
		PairAddress:     pairAddress,
	}
}

func ParseWithdrawLiquidityMemo(ctx cosmos.Context, keeper keeper.Keeper, asset common.Asset, parts []string) (WithdrawLiquidityMemo, error) {
	var err error
	if len(parts) < 2 {
		return WithdrawLiquidityMemo{}, fmt.Errorf("not enough parameters")
	}
	withdrawalBasisPts := cosmos.ZeroUint()
	withdrawalAsset := common.EmptyAsset
	pairAddress := common.NoAddress
	if len(parts) > 2 {
		withdrawalBasisPts, err = cosmos.ParseUint(parts[2])
		if err != nil {
			return WithdrawLiquidityMemo{}, err
		}
		if withdrawalBasisPts.IsZero() || withdrawalBasisPts.GT(cosmos.NewUint(types.MaxWithdrawBasisPoints)) {
			return WithdrawLiquidityMemo{}, fmt.Errorf("withdraw amount %s is invalid", parts[2])
		}
	}
	if len(parts) > 3 {
		withdrawalAsset, err = common.NewAsset(parts[3])
		if err != nil {
			return WithdrawLiquidityMemo{}, err
		}
	}
	if len(parts) > 4 {
		if keeper == nil {
			pairAddress, err = common.NewAddress(parts[4])
			if err != nil {
				return WithdrawLiquidityMemo{}, err
			}
		} else {
			pairAddress, err = FetchAddress(ctx, keeper, parts[4], asset.Chain)
			if err != nil {
				return WithdrawLiquidityMemo{}, err
			}
		}
	}
	return NewWithdrawLiquidityMemo(asset, withdrawalBasisPts, withdrawalAsset, pairAddress), nil
}
