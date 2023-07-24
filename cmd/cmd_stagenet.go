//go:build stagenet
// +build stagenet

package cmd

const (
	Bech32PrefixAccAddr         = "smaya"
	Bech32PrefixAccPub          = "smayapub"
	Bech32PrefixValAddr         = "smayav"
	Bech32PrefixValPub          = "smayavpub"
	Bech32PrefixConsAddr        = "smayac"
	Bech32PrefixConsPub         = "smayacpub"
	DenomRegex                  = `[a-zA-Z][a-zA-Z0-9:\\/\\\-\\_\\.]{2,127}`
	BASEChainCoinType    uint32 = 931
	BASEChainCoinPurpose uint32 = 44
	BASEChainHDPath      string = `m/44'/931'/0'/0/0`
)
