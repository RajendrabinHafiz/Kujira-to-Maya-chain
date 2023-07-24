package mayachain

import (
	"fmt"
	"strconv"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
)

type ManageMAYANameMemo struct {
	MemoBase
	Name           string
	Chain          common.Chain
	Address        common.Address
	PreferredAsset common.Asset
	Expire         int64
	Owner          cosmos.AccAddress
}

func (m ManageMAYANameMemo) GetName() string            { return m.Name }
func (m ManageMAYANameMemo) GetChain() common.Chain     { return m.Chain }
func (m ManageMAYANameMemo) GetAddress() common.Address { return m.Address }
func (m ManageMAYANameMemo) GetBlockExpire() int64      { return m.Expire }

func NewManageMAYANameMemo(name string, chain common.Chain, addr common.Address, expire int64, asset common.Asset, owner cosmos.AccAddress) ManageMAYANameMemo {
	return ManageMAYANameMemo{
		MemoBase:       MemoBase{TxType: TxMAYAName},
		Name:           name,
		Chain:          chain,
		Address:        addr,
		PreferredAsset: asset,
		Expire:         expire,
		Owner:          owner,
	}
}

func ParseManageMAYANameMemo(parts []string) (ManageMAYANameMemo, error) {
	var err error
	var name string
	var owner cosmos.AccAddress
	preferredAsset := common.EmptyAsset
	expire := int64(0)

	if len(parts) < 4 {
		return ManageMAYANameMemo{}, fmt.Errorf("not enough parameters")
	}

	name = parts[1]
	chain, err := common.NewChain(parts[2])
	if err != nil {
		return ManageMAYANameMemo{}, err
	}

	addr, err := common.NewAddress(parts[3])
	if err != nil {
		return ManageMAYANameMemo{}, err
	}

	if len(parts) >= 5 {
		owner, err = cosmos.AccAddressFromBech32(parts[4])
		if err != nil {
			return ManageMAYANameMemo{}, err
		}
	}

	if len(parts) >= 6 {
		preferredAsset, err = common.NewAsset(parts[5])
		if err != nil {
			return ManageMAYANameMemo{}, err
		}
	}

	if len(parts) >= 7 {
		expire, err = strconv.ParseInt(parts[6], 10, 64)
		if err != nil {
			return ManageMAYANameMemo{}, err
		}
	}

	return NewManageMAYANameMemo(name, chain, addr, expire, preferredAsset, owner), nil
}
