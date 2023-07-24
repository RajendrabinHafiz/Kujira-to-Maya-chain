package mayachain

import (
	"github.com/blang/semver"
	"gitlab.com/mayachain/mayanode/common"
)

func GetSupportChains(version semver.Version) common.Chains {
	switch {
	case version.GTE(semver.MustParse("1.105.0")):
		return SUPPORT_CHAINS_V105
	default:
		return SUPPORT_CHAINS_V104
	}
}
