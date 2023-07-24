package types

import (
	"github.com/blang/semver"
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/liquiditypools"
)

func GetLiquidityPools(version semver.Version) common.Assets {
	switch {
	case version.GTE(semver.MustParse("1.105.0")):
		return liquiditypools.LiquidityPoolsV105
	default:
		return liquiditypools.LiquidityPoolsV104
	}
}
