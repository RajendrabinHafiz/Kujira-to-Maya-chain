package mayachain

import (
	"io/ioutil"

	"gitlab.com/mayachain/mayanode/common/cosmos"
	. "gopkg.in/check.v1"
)

type PreMAYANameTestSuite struct{}

var _ = Suite(&PreMAYANameTestSuite{})

func (s *PreMAYANameTestSuite) TestLoadingJson(c *C) {
	// use the mainnet preregister mayanames for test
	var err error
	preregisterMAYANames, err = ioutil.ReadFile("preregister_mayanames_stagenet.json")
	c.Assert(err, IsNil)

	ctx, _ := setupKeeperForTest(c)
	config := cosmos.GetConfig()
	config.SetBech32PrefixForAccount("smaya", "smayapub")
	names, err := getPreRegisterMAYANames(ctx, 100)
	c.Assert(err, IsNil)
	c.Check(names, HasLen, 1, Commentf("%d", len(names)))
}
