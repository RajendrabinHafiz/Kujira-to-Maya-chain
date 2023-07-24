package mayachain

import (
	"fmt"
	"strconv"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"

	"github.com/blang/semver"
)

type BondMemo struct {
	MemoBase
	NodeAddress         cosmos.AccAddress
	BondProviderAddress cosmos.AccAddress
	NodeOperatorFee     int64
	Units               cosmos.Uint
}

func (m BondMemo) GetAmount() cosmos.Uint           { return m.Units }
func (m BondMemo) GetAccAddress() cosmos.AccAddress { return m.NodeAddress }

func NewBondMemo(asset common.Asset, addr, additional cosmos.AccAddress, units cosmos.Uint, operatorFee int64) BondMemo {
	return BondMemo{
		MemoBase: MemoBase{
			TxType: TxBond,
			Asset:  asset,
		},
		NodeAddress:         addr,
		BondProviderAddress: additional,
		NodeOperatorFee:     operatorFee,
		Units:               units,
	}
}

func ParseBondMemo(version semver.Version, parts []string) (BondMemo, error) {
	switch {
	case version.GTE(semver.MustParse("1.105.0")):
		return ParseBondMemoV105(parts)
	case version.GTE(semver.MustParse("1.88.0")):
		return ParseBondMemoV88(parts)
	case version.GTE(semver.MustParse("0.81.0")):
		return ParseBondMemoV81(parts)
	default:
		return BondMemo{}, fmt.Errorf("invalid version(%s)", version.String())
	}
}

func ParseBondMemoV105(parts []string) (BondMemo, error) {
	var err error
	var asset common.Asset
	units := cosmos.ZeroUint()
	additional := cosmos.AccAddress{}
	var operatorFee int64 = -1
	if len(parts) < 3 {
		return BondMemo{}, fmt.Errorf("not enough parameters")
	}

	if asset, err = common.NewAsset(parts[1]); err == nil {
		if len(parts) < 4 {
			return BondMemo{}, fmt.Errorf("not enough parameters")
		}

		units, err = cosmos.ParseUint(parts[2])
		if err != nil {
			return BondMemo{}, fmt.Errorf("%s is an invalid bond units: %w", parts[2], err)
		}

		// Remove asset and units from parts
		parts = parts[2:]
	}
	addr, err := cosmos.AccAddressFromBech32(parts[1])
	if err != nil {
		return BondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[3], err)
	}
	if len(parts) >= 3 {
		additional, err = cosmos.AccAddressFromBech32(parts[2])
		if err != nil {
			return BondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[4], err)
		}
	}
	if len(parts) >= 4 {
		operatorFee, err = strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return BondMemo{}, fmt.Errorf("%s invalid operator fee: %w", parts[5], err)
		}
	}
	mem := NewBondMemo(asset, addr, additional, units, operatorFee)
	return mem, nil
}

func ParseBondMemoV88(parts []string) (BondMemo, error) {
	additional := cosmos.AccAddress{}
	var operatorFee int64 = -1
	if len(parts) < 2 {
		return BondMemo{}, fmt.Errorf("not enough parameters")
	}
	addr, err := cosmos.AccAddressFromBech32(parts[1])
	if err != nil {
		return BondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[1], err)
	}
	if len(parts) == 3 || len(parts) == 4 {
		additional, err = cosmos.AccAddressFromBech32(parts[2])
		if err != nil {
			return BondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[2], err)
		}
	}
	if len(parts) == 4 {
		operatorFee, err = strconv.ParseInt(parts[3], 10, 64)
		if err != nil {
			return BondMemo{}, fmt.Errorf("%s invalid operator fee: %w", parts[3], err)
		}
	}
	return NewBondMemo(common.EmptyAsset, addr, additional, cosmos.ZeroUint(), operatorFee), nil
}

func ParseBondMemoV81(parts []string) (BondMemo, error) {
	additional := cosmos.AccAddress{}
	if len(parts) < 2 {
		return BondMemo{}, fmt.Errorf("not enough parameters")
	}
	addr, err := cosmos.AccAddressFromBech32(parts[1])
	if err != nil {
		return BondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[1], err)
	}
	if len(parts) >= 3 {
		additional, err = cosmos.AccAddressFromBech32(parts[2])
		if err != nil {
			return BondMemo{}, fmt.Errorf("%s is an invalid thorchain address: %w", parts[2], err)
		}
	}
	return NewBondMemo(common.EmptyAsset, addr, additional, cosmos.ZeroUint(), -1), nil
}
