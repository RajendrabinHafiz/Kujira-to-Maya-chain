package types

const (
	// ModuleName name of BASEChain module
	ModuleName = "mayachain"
	// DefaultCodespace is the same as ModuleName
	DefaultCodespace = ModuleName
	// ReserveName the module account name to keep reserve
	ReserveName = "reserve"
	// AsgardName the module account name to keep asgard fund
	AsgardName = "asgard"
	// BondName the name of account used to store bond
	BondName = "bond"
	// MayaFund the module account name to keep the maya fund
	MayaFund = "maya_fund"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName
	// RouterKey used in the RPC query
	RouterKey = ModuleName // this was defined in your key.go file
)
