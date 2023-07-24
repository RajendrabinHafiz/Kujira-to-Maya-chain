package keeperv1

import (
	"fmt"

	"gitlab.com/mayachain/mayanode/common/cosmos"
)

func (k KVStore) setMAYAName(ctx cosmos.Context, key string, record MAYAName) {
	store := ctx.KVStore(k.storeKey)
	buf := k.cdc.MustMarshal(&record)
	if buf == nil {
		store.Delete([]byte(key))
	} else {
		store.Set([]byte(key), buf)
	}
}

func (k KVStore) getMAYAName(ctx cosmos.Context, key string, record *MAYAName) (bool, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return false, nil
	}

	bz := store.Get([]byte(key))
	if err := k.cdc.Unmarshal(bz, record); err != nil {
		return true, dbError(ctx, fmt.Sprintf("Unmarshal kvstore: (%T) %s", record, key), err)
	}

	return true, nil
}

// GetMAYANameIterator only iterate MAYANames
func (k KVStore) GetMAYANameIterator(ctx cosmos.Context) cosmos.Iterator {
	return k.getIterator(ctx, prefixMAYAName)
}

// SetMAYAName save the MAYAName object to store
func (k KVStore) SetMAYAName(ctx cosmos.Context, name MAYAName) {
	k.setMAYAName(ctx, k.GetKey(ctx, prefixMAYAName, name.Key()), name)
}

// MAYANameExists check whether the given name exists
func (k KVStore) MAYANameExists(ctx cosmos.Context, name string) bool {
	record := MAYAName{
		Name: name,
	}
	if k.has(ctx, k.GetKey(ctx, prefixMAYAName, record.Key())) {
		record, _ = k.GetMAYAName(ctx, name)
		return record.ExpireBlockHeight >= ctx.BlockHeight()
	}
	return false
}

// GetMAYAName get MAYAName with the given pubkey from data store
func (k KVStore) GetMAYAName(ctx cosmos.Context, name string) (MAYAName, error) {
	record := MAYAName{
		Name: name,
	}
	ok, err := k.getMAYAName(ctx, k.GetKey(ctx, prefixMAYAName, record.Key()), &record)
	if !ok {
		return record, fmt.Errorf("MAYAName doesn't exist: %s", name)
	}
	if record.ExpireBlockHeight < ctx.BlockHeight() {
		return MAYAName{Name: name}, nil
	}
	return record, err
}

// DeleteMAYAName remove the given MAYAName from data store
func (k KVStore) DeleteMAYAName(ctx cosmos.Context, name string) error {
	n := MAYAName{Name: name}
	k.del(ctx, k.GetKey(ctx, prefixMAYAName, n.Key()))
	return nil
}
