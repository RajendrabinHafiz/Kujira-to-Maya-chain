package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
)

type Assets []Asset

var (
	// EmptyAsset empty asset, not valid
	EmptyAsset = Asset{Chain: EmptyChain, Symbol: "", Ticker: "", Synth: false}
	// RUNEAsset RUNE
	RUNEAsset = Asset{Chain: THORChain, Symbol: "RUNE", Ticker: "RUNE", Synth: false}
	// ATOMAsset ATOM
	ATOMAsset = Asset{Chain: GAIAChain, Symbol: "ATOM", Ticker: "ATOM", Synth: false}
	// KUJIAsset KUJI
	KUJIAsset = Asset{Chain: KUJIChain, Symbol: "KUJI", Ticker: "KUJI", Synth: false}
	// BNBAsset BNB
	BNBAsset = Asset{Chain: BNBChain, Symbol: "BNB", Ticker: "BNB", Synth: false}
	// BTCAsset BTC
	BTCAsset = Asset{Chain: BTCChain, Symbol: "BTC", Ticker: "BTC", Synth: false}
	// LTCAsset BTC
	LTCAsset = Asset{Chain: LTCChain, Symbol: "LTC", Ticker: "LTC", Synth: false}
	// BCHAsset BCH
	BCHAsset = Asset{Chain: BCHChain, Symbol: "BCH", Ticker: "BCH", Synth: false}
	// DASHAsset DASH
	DASHAsset = Asset{Chain: DASHChain, Symbol: "DASH", Ticker: "DASH", Synth: false}
	// DOGEAsset DOGE
	DOGEAsset = Asset{Chain: DOGEChain, Symbol: "DOGE", Ticker: "DOGE", Synth: false}
	// ETHAsset ETH
	ETHAsset = Asset{Chain: ETHChain, Symbol: "ETH", Ticker: "ETH", Synth: false}
	// USDTAsset ETH
	USDTAsset = Asset{Chain: ETHChain, Symbol: "USDT-0xdAC17F958D2ee523a2206206994597C13D831ec7", Ticker: "ETH", Synth: false}
	// AVAXAsset AVAX
	AVAXAsset = Asset{Chain: AVAXChain, Symbol: "AVAX", Ticker: "AVAX", Synth: false}
	// Rune67CAsset RUNE on Binance test net
	// BaseNative CACAO on mayachain
	BaseNative = Asset{Chain: BASEChain, Symbol: "CACAO", Ticker: "CACAO", Synth: false}
	MayaNative = Asset{Chain: BASEChain, Symbol: "MAYA", Ticker: "MAYA", Synth: false}
)

// NewAsset parse the given input into Asset object
func NewAsset(input string) (Asset, error) {
	var err error
	var asset Asset
	var sym string
	var parts []string
	if strings.Count(input, "/") > 0 {
		parts = strings.SplitN(input, "/", 2)
		asset.Synth = true
	} else {
		parts = strings.SplitN(input, ".", 2)
		asset.Synth = false
	}
	if len(parts) == 1 {
		asset.Chain = BASEChain
		sym = parts[0]
	} else {
		asset.Chain, err = NewChain(parts[0])
		if err != nil {
			return EmptyAsset, err
		}
		sym = parts[1]
	}

	asset.Symbol, err = NewSymbol(sym)
	if err != nil {
		return EmptyAsset, err
	}

	parts = strings.SplitN(sym, "-", 2)
	asset.Ticker, err = NewTicker(parts[0])
	if err != nil {
		return EmptyAsset, err
	}

	return asset, nil
}

// Equals determinate whether two assets are equivalent
func (a Asset) Equals(a2 Asset) bool {
	return a.Chain.Equals(a2.Chain) && a.Symbol.Equals(a2.Symbol) && a.Ticker.Equals(a2.Ticker) && a.Synth == a2.Synth
}

func (a Asset) GetChain() Chain {
	if a.Synth {
		return BASEChain
	}
	return a.Chain
}

// Get layer1 asset version
func (a Asset) GetLayer1Asset() Asset {
	if !a.IsSyntheticAsset() {
		return a
	}
	return Asset{
		Chain:  a.Chain,
		Symbol: a.Symbol,
		Ticker: a.Ticker,
		Synth:  false,
	}
}

// Get synthetic asset of asset
func (a Asset) GetSyntheticAsset() Asset {
	if a.IsSyntheticAsset() {
		return a
	}
	return Asset{
		Chain:  a.Chain,
		Symbol: a.Symbol,
		Ticker: a.Ticker,
		Synth:  true,
	}
}

// Check if asset is a pegged asset
func (a Asset) IsSyntheticAsset() bool {
	return a.Synth
}

func (a Asset) IsVaultAsset() bool {
	return a.IsSyntheticAsset()
}

// Native return native asset, only relevant on THORChain
func (a Asset) Native() string {
	if a.IsBase() {
		return "cacao"
	}
	if a.Equals(MayaNative) {
		return "maya"
	}
	return strings.ToLower(a.String())
}

// IsEmpty will be true when any of the field is empty, chain,symbol or ticker
func (a Asset) IsEmpty() bool {
	return a.Chain.IsEmpty() || a.Symbol.IsEmpty() || a.Ticker.IsEmpty()
}

// String implement fmt.Stringer , return the string representation of Asset
func (a Asset) String() string {
	div := "."
	if a.Synth {
		div = "/"
	}
	return fmt.Sprintf("%s%s%s", a.Chain.String(), div, a.Symbol.String())
}

// IsGasAsset check whether asset is base asset used to pay for gas
func (a Asset) IsGasAsset() bool {
	gasAsset := a.GetChain().GetGasAsset()
	if gasAsset.IsEmpty() {
		return false
	}
	return a.Equals(gasAsset)
}

// IsCacao is a helper function ,return true only when the asset represent RUNE
func (a Asset) IsBase() bool {
	return a.Equals(BaseNative)
}

// IsNativeRune is a helper function, return true only when the asset represent NATIVE RUNE
func (a Asset) IsNativeBase() bool {
	return a.IsBase() && a.Chain.IsBASEChain()
}

// IsNative is a helper function, returns true when the asset is a native
// asset to THORChain (ie rune, a synth, etc)
func (a Asset) IsNative() bool {
	return a.GetChain().IsBASEChain()
}

// IsBNB is a helper function, return true only when the asset represent BNB
func (a Asset) IsBNB() bool {
	return a.Equals(BNBAsset)
}

// MarshalJSON implement Marshaler interface
func (a Asset) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// UnmarshalJSON implement Unmarshaler interface
func (a *Asset) UnmarshalJSON(data []byte) error {
	var err error
	var assetStr string
	if err := json.Unmarshal(data, &assetStr); err != nil {
		return err
	}
	if assetStr == "." {
		*a = EmptyAsset
		return nil
	}
	*a, err = NewAsset(assetStr)
	return err
}

// MarshalJSONPB implement jsonpb.Marshaler
func (a Asset) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return a.MarshalJSON()
}

// UnmarshalJSONPB implement jsonpb.Unmarshaler
func (a *Asset) UnmarshalJSONPB(unmarshal *jsonpb.Unmarshaler, content []byte) error {
	return a.UnmarshalJSON(content)
}

// Contains checks if the array contains the specified element
func (as *Assets) Contains(a Asset) bool {
	for _, asset := range *as {
		if asset.Equals(a) {
			return true
		}
	}
	return false
}

// BaseAsset return RUNE Asset depends on different environment
func BaseAsset() Asset {
	return BaseNative
}

// Replace pool name "." with a "-" for Mimir key checking.
func (a Asset) MimirString() string {
	return a.Chain.String() + "-" + a.Symbol.String()
}

// GetAsset returns true if the asset exists in the list of assets
func ContainsAsset(asset Asset, assets []Asset) bool {
	for _, a := range assets {
		if a.Equals(asset) {
			return true
		}
	}
	return false
}
