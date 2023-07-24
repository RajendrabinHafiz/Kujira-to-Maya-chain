package mayachain

import (
	"gitlab.com/mayachain/mayanode/common/cosmos"
)

type DummyPoolManager struct{}

func NewDummyPoolManager() *DummyPoolManager {
	return &DummyPoolManager{}
}

func (m *DummyPoolManager) EndBlock(ctx cosmos.Context, mgr Manager) error {
	return nil
}
