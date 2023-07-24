package gaia

import "strings"

type CosmosAssetMapping struct {
	CosmosDenom     string
	CosmosDecimals  int
	BASEChainSymbol string
}

// CosmosAssetMappings maps a Cosmos denom to a BASEChain symbol and provides the asset decimals
// CHANGEME: define assets that should be observed by BASEChain here. This also acts a whitelist.
var CosmosAssetMappings = []CosmosAssetMapping{
	{
		CosmosDenom:     "uatom",
		CosmosDecimals:  6,
		BASEChainSymbol: "ATOM",
	},
}

func GetAssetByCosmosDenom(denom string) (CosmosAssetMapping, bool) {
	for _, asset := range CosmosAssetMappings {
		if strings.EqualFold(asset.CosmosDenom, denom) {
			return asset, true
		}
	}
	return CosmosAssetMapping{}, false
}

func GetAssetByThorchainSymbol(symbol string) (CosmosAssetMapping, bool) {
	for _, asset := range CosmosAssetMappings {
		if strings.EqualFold(asset.BASEChainSymbol, symbol) {
			return asset, true
		}
	}
	return CosmosAssetMapping{}, false
}
