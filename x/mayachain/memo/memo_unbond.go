package mayachain

import (
	"fmt"

	"github.com/blang/semver"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
)

type UnbondMemo struct {
	MemoBase
	NodeAddress         cosmos.AccAddress
	BondProviderAddress cosmos.AccAddress
	Units               cosmos.Uint
}

func (m UnbondMemo) GetAccAddress() cosmos.AccAddress { return m.NodeAddress }

func NewUnbondMemo(asset common.Asset, addr, additional cosmos.AccAddress, units cosmos.Uint) UnbondMemo {
	return UnbondMemo{
		MemoBase: MemoBase{
			TxType: TxUnbond,
			Asset:  asset,
		},
		NodeAddress:         addr,
		BondProviderAddress: additional,
		Units:               units,
	}
}

func ParseUnbondMemo(version semver.Version, parts []string) (UnbondMemo, error) {
	switch {
	case version.GTE(semver.MustParse("1.105.0")):
		return ParseUnbondMemoV105(parts)
	case version.GTE(semver.MustParse("0.81.0")):
		return ParseUnbondMemoV81(parts)
	}
	return UnbondMemo{}, fmt.Errorf("invalid version(%s)", version.String())
}

func ParseUnbondMemoV105(parts []string) (UnbondMemo, error) {
	var err error
	var asset common.Asset
	units := cosmos.ZeroUint()
	additional := cosmos.AccAddress{}
	if len(parts) < 2 {
		return UnbondMemo{}, fmt.Errorf("not enough parameters")
	}

	if asset, err = common.NewAsset(parts[1]); err == nil {
		if len(parts) < 4 {
			return UnbondMemo{}, fmt.Errorf("not enough parameters")
		}

		units, err = cosmos.ParseUint(parts[2])
		if err != nil {
			return UnbondMemo{}, fmt.Errorf("%s is an invalid unbond units: %w", parts[2], err)
		}

		// Remove asset and units from parts
		parts = parts[2:]
	}
	addr, err := cosmos.AccAddressFromBech32(parts[1])
	if err != nil {
		return UnbondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[1], err)
	}
	if len(parts) >= 3 {
		additional, err = cosmos.AccAddressFromBech32(parts[2])
		if err != nil {
			return UnbondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[2], err)
		}
	}
	return NewUnbondMemo(asset, addr, additional, units), nil
}

func ParseUnbondMemoV81(parts []string) (UnbondMemo, error) {
	additional := cosmos.AccAddress{}
	if len(parts) < 2 {
		return UnbondMemo{}, fmt.Errorf("not enough parameters")
	}
	addr, err := cosmos.AccAddressFromBech32(parts[1])
	if err != nil {
		return UnbondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[1], err)
	}
	if len(parts) >= 3 {
		additional, err = cosmos.AccAddressFromBech32(parts[2])
		if err != nil {
			return UnbondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[2], err)
		}
	}
	return NewUnbondMemo(common.EmptyAsset, addr, additional, cosmos.ZeroUint()), nil
}
