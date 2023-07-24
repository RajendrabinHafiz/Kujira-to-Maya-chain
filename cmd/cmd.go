//go:build !testnet && !mocknet && !stagenet
// +build !testnet,!mocknet,!stagenet

package cmd

const (
	Bech32PrefixAccAddr         = "maya"
	Bech32PrefixAccPub          = "mayapub"
	Bech32PrefixValAddr         = "mayav"
	Bech32PrefixValPub          = "mayavpub"
	Bech32PrefixConsAddr        = "mayac"
	Bech32PrefixConsPub         = "mayacpub"
	DenomRegex                  = `[a-zA-Z][a-zA-Z0-9:\\/\\\-\\_\\.]{2,127}`
	BASEChainCoinType    uint32 = 931
	BASEChainCoinPurpose uint32 = 44
	BASEChainHDPath      string = `m/44'/931'/0'/0/0`
)
