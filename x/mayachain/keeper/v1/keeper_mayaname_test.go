package keeperv1

import (
	"gitlab.com/mayachain/mayanode/common"
	. "gopkg.in/check.v1"
)

type KeeperMAYANameSuite struct{}

var _ = Suite(&KeeperMAYANameSuite{})

func (s *KeeperMAYANameSuite) TestMAYAName(c *C) {
	ctx, k := setupKeeperForTest(c)
	var err error
	ref := "helloworld"

	ok := k.MAYANameExists(ctx, ref)
	c.Assert(ok, Equals, false)

	thorAddr := GetRandomBaseAddress()
	bnbAddr := GetRandomBNBAddress()
	name := NewMAYAName(ref, 50, []MAYANameAlias{{Chain: common.BASEChain, Address: thorAddr}, {Chain: common.BNBChain, Address: bnbAddr}})
	k.SetMAYAName(ctx, name)

	ok = k.MAYANameExists(ctx, ref)
	c.Assert(ok, Equals, true)
	ok = k.MAYANameExists(ctx, "bogus")
	c.Assert(ok, Equals, false)

	name, err = k.GetMAYAName(ctx, ref)
	c.Assert(err, IsNil)
	c.Assert(name.GetAlias(common.BASEChain).Equals(thorAddr), Equals, true)
	c.Assert(name.GetAlias(common.BNBChain).Equals(bnbAddr), Equals, true)
}
