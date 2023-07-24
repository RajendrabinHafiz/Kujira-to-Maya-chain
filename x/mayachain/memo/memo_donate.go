package mayachain

import (
	"fmt"

	"gitlab.com/mayachain/mayanode/common"
)

type DonateMemo struct{ MemoBase }

func (m DonateMemo) String() string {
	return fmt.Sprintf("DONATE:%s", m.Asset)
}

func NewDonateMemo(asset common.Asset) DonateMemo {
	return DonateMemo{
		MemoBase: MemoBase{TxType: TxDonate, Asset: asset},
	}
}
