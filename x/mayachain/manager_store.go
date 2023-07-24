package mayachain

import (
	"fmt"

	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
)

// StoreManager define the method as the entry point for store upgrade
type StoreManager interface {
	Iterator(_ cosmos.Context) error
}

// StoreMgr implement StoreManager interface
type StoreMgr struct {
	mgr *Mgrs
}

// newStoreMgr create a new instance of StoreMgr
func newStoreMgr(mgr *Mgrs) *StoreMgr {
	return &StoreMgr{
		mgr: mgr,
	}
}

// Iterator implement StoreManager interface decide whether it need to upgrade store
func (smgr *StoreMgr) Iterator(ctx cosmos.Context) error {
	version := smgr.mgr.GetVersion()

	if version.Major > constants.SWVersion.Major || version.Minor > constants.SWVersion.Minor {
		return fmt.Errorf("out of date software: have %s, network running %s", constants.SWVersion, version)
	}

	storeVer := smgr.mgr.Keeper().GetStoreVersion(ctx)
	if storeVer < 0 {
		return fmt.Errorf("unable to get store version: %d", storeVer)
	}
	if uint64(storeVer) < version.Minor {
		for i := uint64(storeVer + 1); i <= version.Minor; i++ {
			if err := smgr.migrate(ctx, i); err != nil {
				return err
			}
		}
	} else {
		ctx.Logger().Debug("No store migration needed")
	}
	return nil
}

func (smgr *StoreMgr) migrate(ctx cosmos.Context, i uint64) error {
	ctx.Logger().Info("Migrating store to new version", "version", i)
	// add the logic to migrate store here when it is needed

	switch i {
	case 96:
		migrateStoreV96(ctx, smgr.mgr)
	case 102:
		migrateStoreV102(ctx, smgr.mgr)
	case 104:
		migrateStoreV104(ctx, smgr.mgr)
	case 105:
		migrateStoreV105(ctx, smgr.mgr)
	}

	smgr.mgr.Keeper().SetStoreVersion(ctx, int64(i))
	return nil
}
