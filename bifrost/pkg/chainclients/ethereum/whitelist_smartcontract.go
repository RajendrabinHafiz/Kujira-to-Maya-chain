//go:build !testnet && !mocknet && !stagenet
// +build !testnet,!mocknet,!stagenet

package ethereum

import (
	"gitlab.com/mayachain/mayanode/common"
)

var whitelistSmartContractAddress = []common.Address{}
