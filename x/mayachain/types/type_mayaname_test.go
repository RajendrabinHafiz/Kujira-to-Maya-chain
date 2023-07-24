package types

import (
	. "gopkg.in/check.v1"

	"gitlab.com/mayachain/mayanode/common"
)

type MAYANameSuite struct{}

var _ = Suite(&MAYANameSuite{})

func (MAYANameSuite) TestMAYAName(c *C) {
	// happy path
	n := NewMAYAName("iamthewalrus", 0, []MAYANameAlias{{Chain: common.BASEChain, Address: GetRandomBaseAddress()}})
	c.Check(n.Valid(), IsNil)

	// unhappy path
	n1 := NewMAYAName("", 0, []MAYANameAlias{{Chain: common.BNBChain, Address: GetRandomBaseAddress()}})
	c.Check(n1.Valid(), NotNil)
	n2 := NewMAYAName("hello", 0, []MAYANameAlias{{Chain: common.EmptyChain, Address: GetRandomBaseAddress()}})
	c.Check(n2.Valid(), NotNil)
	n3 := NewMAYAName("hello", 0, []MAYANameAlias{{Chain: common.BASEChain, Address: common.Address("")}})
	c.Check(n3.Valid(), NotNil)
}
