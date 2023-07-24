package types

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
)

var _ codec.ProtoMarshaler = &LiquidityProvider{}

// LiquidityProviders a list of liquidity providers
type LiquidityProviders []LiquidityProvider

// Valid check whether lp represent valid information
func (m *LiquidityProvider) Valid() error {
	if m.LastAddHeight == 0 {
		return errors.New("last add liquidity height cannot be empty")
	}
	if m.AssetAddress.IsEmpty() && m.CacaoAddress.IsEmpty() {
		return errors.New("asset address and rune address cannot be empty")
	}
	return nil
}

func (lp LiquidityProvider) GetAddress() common.Address {
	if !lp.CacaoAddress.IsEmpty() {
		return lp.CacaoAddress
	}
	return lp.AssetAddress
}

// Key return a string which can be used to identify lp
func (lp LiquidityProvider) Key() string {
	return fmt.Sprintf("%s/%s", lp.Asset.String(), lp.GetAddress().String())
}

// Deprecated: do not use
func (lp LiquidityProvider) IsLiquidityBondProvider() bool {
	return !lp.NodeBondAddress.Empty()
}

// Deprecated: do not use
func (lps LiquidityProviders) SetNodeAccount(na cosmos.AccAddress) {
	for i := range lps {
		lps[i].NodeBondAddress = na
	}
}

// Bond creates a new bond record for the given node or increase its units if it already exists
func (lp *LiquidityProvider) Bond(nodeAddr cosmos.AccAddress, units cosmos.Uint) {
	for i := range lp.BondedNodes {
		if lp.BondedNodes[i].NodeAddress.Equals(nodeAddr) {
			lp.BondedNodes[i].Units = lp.BondedNodes[i].Units.Add(units)
			return
		}
	}

	lp.BondedNodes = append(lp.BondedNodes, LPBondedNode{
		NodeAddress: nodeAddr,
		Units:       units,
	})
}

// Unbond removes a bond record for the given node or decrease its units if it already exists
func (lp *LiquidityProvider) Unbond(nodeAddr cosmos.AccAddress, units cosmos.Uint) {
	// Soft migration to new bond model
	if !lp.NodeBondAddress.Empty() {
		lp.NodeBondAddress = nil

		if lp.Units.GT(units) {
			lp.BondedNodes = []LPBondedNode{
				{
					NodeAddress: nodeAddr,
					Units:       common.SafeSub(lp.Units, units),
				},
			}
		}

		return
	}

	for i := range lp.BondedNodes {
		if lp.BondedNodes[i].NodeAddress.Equals(nodeAddr) {
			lp.BondedNodes[i].Units = common.SafeSub(lp.BondedNodes[i].Units, units)
			if lp.BondedNodes[i].Units.IsZero() {
				lp.BondedNodes = append(lp.BondedNodes[:i], lp.BondedNodes[i+1:]...)
			}
			return
		}
	}
}

// GetRemainingUnits returns the number of LP units that are not bonded to a node
func (lp *LiquidityProvider) GetRemainingUnits() cosmos.Uint {
	if !lp.NodeBondAddress.Empty() {
		return cosmos.ZeroUint()
	}

	bondedUnits := cosmos.ZeroUint()
	for _, bond := range lp.BondedNodes {
		bondedUnits = bondedUnits.Add(bond.Units)
	}

	return common.SafeSub(lp.Units, bondedUnits)
}

// GetUnitsBondedToNode returns the number of LP units that are bonded to the given node
func (lp LiquidityProvider) GetUnitsBondedToNode(nodeAddr cosmos.AccAddress) cosmos.Uint {
	if lp.NodeBondAddress.Equals(nodeAddr) {
		return lp.Units
	}

	for _, bond := range lp.BondedNodes {
		if bond.NodeAddress.Equals(nodeAddr) {
			return bond.Units
		}
	}

	return cosmos.ZeroUint()
}

// IsEmpty returns true when the LP is empty
func (bn LPBondedNode) IsEmpty() bool {
	return bn.NodeAddress.Empty()
}
