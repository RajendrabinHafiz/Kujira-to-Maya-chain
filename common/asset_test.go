package common

import (
	. "gopkg.in/check.v1"
)

type AssetSuite struct{}

var _ = Suite(&AssetSuite{})

func (s AssetSuite) TestAsset(c *C) {
	asset, err := NewAsset("maya.cacao")
	c.Assert(err, IsNil)
	c.Check(asset.Equals(BaseNative), Equals, true)
	c.Check(asset.IsBase(), Equals, true)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.Synth, Equals, false)
	c.Check(asset.String(), Equals, "MAYA.CACAO")

	asset, err = NewAsset("maya/cacao")
	c.Assert(err, IsNil)
	c.Check(asset.Equals(BaseNative), Equals, false)
	c.Check(asset.IsBase(), Equals, false)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.Synth, Equals, true)
	c.Check(asset.String(), Equals, "MAYA/CACAO")

	c.Check(asset.Chain.Equals(BASEChain), Equals, true)
	c.Check(asset.Symbol.Equals(Symbol("CACAO")), Equals, true)
	c.Check(asset.Ticker.Equals(Ticker("CACAO")), Equals, true)

	asset, err = NewAsset("BNB.SWIPE.B-DC0")
	c.Assert(err, IsNil)
	c.Check(asset.String(), Equals, "BNB.SWIPE.B-DC0")
	c.Check(asset.Chain.Equals(BNBChain), Equals, true)
	c.Check(asset.Symbol.Equals(Symbol("SWIPE.B-DC0")), Equals, true)
	c.Check(asset.Ticker.Equals(Ticker("SWIPE.B")), Equals, true)

	// parse without chain
	asset, err = NewAsset("cacao")
	c.Assert(err, IsNil)
	c.Check(asset.Equals(BaseNative), Equals, true)

	// ETH test
	asset, err = NewAsset("eth.knc")
	c.Assert(err, IsNil)
	c.Check(asset.Chain.Equals(ETHChain), Equals, true)
	c.Check(asset.Symbol.Equals(Symbol("KNC")), Equals, true)
	c.Check(asset.Ticker.Equals(Ticker("KNC")), Equals, true)
	asset, err = NewAsset("ETH.CACAO-0x3155ba85d5f96b2d030a4966af206230e46849cb")
	c.Assert(err, IsNil)

	// DASH test
	asset, err = NewAsset("dash.dash")
	c.Assert(err, IsNil)
	c.Check(asset.Chain.Equals(DASHChain), Equals, true)
	c.Check(asset.Equals(DASHAsset), Equals, true)
	c.Check(asset.IsBase(), Equals, false)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.String(), Equals, "DASH.DASH")

	// DOGE test
	asset, err = NewAsset("doge.doge")
	c.Assert(err, IsNil)
	c.Check(asset.Chain.Equals(DOGEChain), Equals, true)
	c.Check(asset.Equals(DOGEAsset), Equals, true)
	c.Check(asset.IsBase(), Equals, false)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.String(), Equals, "DOGE.DOGE")

	// BCH test
	asset, err = NewAsset("bch.bch")
	c.Assert(err, IsNil)
	c.Check(asset.Chain.Equals(BCHChain), Equals, true)
	c.Check(asset.Equals(BCHAsset), Equals, true)
	c.Check(asset.IsBase(), Equals, false)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.String(), Equals, "BCH.BCH")

	// LTC test
	asset, err = NewAsset("ltc.ltc")
	c.Assert(err, IsNil)
	c.Check(asset.Chain.Equals(LTCChain), Equals, true)
	c.Check(asset.Equals(LTCAsset), Equals, true)
	c.Check(asset.IsBase(), Equals, false)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.String(), Equals, "LTC.LTC")

	// btc/btc
	asset, err = NewAsset("btc/btc")
	c.Check(err, IsNil)
	c.Check(asset.Chain.Equals(BTCChain), Equals, true)
	c.Check(asset.Equals(BTCAsset), Equals, false)
	c.Check(asset.IsEmpty(), Equals, false)
	c.Check(asset.String(), Equals, "BTC/BTC")
}
