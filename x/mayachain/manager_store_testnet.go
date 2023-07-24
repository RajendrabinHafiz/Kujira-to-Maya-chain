//go:build testnet
// +build testnet

package mayachain

import (
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
)

// migrateStoreV86 remove all LTC asset from the retiring vault
func migrateStoreV86(ctx cosmos.Context, mgr *Mgrs) {
	defer func() {
		if err := recover(); err != nil {
			ctx.Logger().Error("fail to migrate store to v86", "error", err)
		}
	}()
	vaults, err := mgr.Keeper().GetAsgardVaultsByStatus(ctx, RetiringVault)
	if err != nil {
		ctx.Logger().Error("fail to get retiring asgard vaults", "error", err)
		return
	}
	for _, v := range vaults {
		ltcCoin := v.GetCoin(common.LTCAsset)
		v.SubFunds(common.NewCoins(ltcCoin))
		if err := mgr.Keeper().SetVault(ctx, v); err != nil {
			ctx.Logger().Error("fail to save vault", "error", err)
		}
	}
}

func migrateStoreV88(ctx cosmos.Context, mgr Manager)  {}
func migrateStoreV90(ctx cosmos.Context, mgr Manager)  {}
func migrateStoreV104(ctx cosmos.Context, mgr Manager) {}
func migrateStoreV105(ctx cosmos.Context, mgr Manager) {}
