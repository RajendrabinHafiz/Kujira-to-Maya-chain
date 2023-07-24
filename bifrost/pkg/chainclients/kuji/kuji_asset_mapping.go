package kuji

import "strings"

type KujiAssetMapping struct {
	KujiDenom       string
	KujiDecimals    int
	BASEChainSymbol string
}

// KujiAssetMappings maps a Kuji denom to a BASEChain symbol and provides the asset decimals
// CHANGEME: define assets that should be observed by BASEChain here. This also acts a whitelist.
var KujiAssetMappings = []KujiAssetMapping{
	{
		KujiDenom:       "ukujira",
		KujiDecimals:    6,
		BASEChainSymbol: "KUJI",
	},
}

func GetAssetByKujiDenom(denom string) (KujiAssetMapping, bool) {
	for _, asset := range KujiAssetMappings {
		if strings.EqualFold(asset.KujiDenom, denom) {
			return asset, true
		}
	}
	return KujiAssetMapping{}, false
}

func GetAssetByMayachainSymbol(symbol string) (KujiAssetMapping, bool) {
	for _, asset := range KujiAssetMappings {
		if strings.EqualFold(asset.BASEChainSymbol, symbol) {
			return asset, true
		}
	}
	return KujiAssetMapping{}, false
}
