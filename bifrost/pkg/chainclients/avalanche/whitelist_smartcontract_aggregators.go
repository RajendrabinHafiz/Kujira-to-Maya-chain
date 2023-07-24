package avalanche

import (
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/x/mayachain/aggregators"
)

func LatestAggregatorContracts() []common.Address {
	addrs := []common.Address{}
	for _, agg := range aggregators.DexAggregators(common.LatestVersion) {
		if agg.Chain.Equals(common.AVAXChain) {
			addrs = append(addrs, common.Address(agg.Address))
		}
	}
	return addrs
}
