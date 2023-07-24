package mayaclient

import (
	"net/http"
	"net/http/httptest"

	. "gopkg.in/check.v1"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/x/mayachain"
)

type PoolManagerTestSuite struct {
	server *httptest.Server
	bridge *mayachainBridge
}

var _ = Suite(&PoolManagerTestSuite{})

func (s *PoolManagerTestSuite) SetUpSuite(c *C) {
	mayachain.SetupConfigForTest()
	cfg, _, kb := SetupMayachainForTest(c)
	s.server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.RequestURI == PoolsEndpoint {
			httpTestHandler(c, rw, "../../test/fixtures/endpoints/pools/pools.json")
		}
	}))
	cfg.ChainHost = s.server.Listener.Addr().String()
	cfg.ChainRPC = s.server.Listener.Addr().String()
	bridge, err := NewMayachainBridge(cfg, GetMetricForTest(c), NewKeysWithKeybase(kb, cfg.SignerName, cfg.SignerPasswd))
	c.Assert(err, IsNil)
	var ok bool
	s.bridge, ok = bridge.(*mayachainBridge)
	c.Assert(ok, Equals, true)
	s.bridge.httpClient.RetryMax = 1 // fail fast
	c.Assert(err, IsNil)
	c.Assert(s.bridge, NotNil)
}

func (s *PoolManagerTestSuite) TestGetPrice(c *C) {
	poolMgr := NewPoolMgr(s.bridge)
	c.Assert(poolMgr, NotNil)
	value, err := poolMgr.GetValue(common.BNBAsset, common.ETHAsset, cosmos.NewUint(1000))
	c.Assert(err, NotNil)
	c.Assert(value.IsZero(), Equals, true)
	asset, err := common.NewAsset("ETH.TKN-0X3B7FA4DD21C6F9BA3CA375217EAD7CAB9D6BF483")
	c.Assert(err, IsNil)
	value, err = poolMgr.GetValue(asset, common.ETHAsset, cosmos.NewUint(1000))
	c.Assert(err, IsNil)
	c.Assert(value.IsZero(), Equals, false)
	c.Assert(value.String(), Equals, "564")
}
