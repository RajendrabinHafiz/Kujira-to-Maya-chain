package mayachain

import (
	. "gopkg.in/check.v1"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/constants"
	"gitlab.com/mayachain/mayanode/x/mayachain/keeper"
)

type HandlerManageMAYANameSuite struct{}

var _ = Suite(&HandlerManageMAYANameSuite{})

type KeeperManageMAYANameTest struct {
	keeper.Keeper
}

func NewKeeperManageMAYANameTest(k keeper.Keeper) KeeperManageMAYANameTest {
	return KeeperManageMAYANameTest{Keeper: k}
}

func (s *HandlerManageMAYANameSuite) TestValidator(c *C) {
	ctx, mgr := setupManagerForTest(c)

	handler := NewManageMAYANameHandler(mgr)
	coin := common.NewCoin(common.BaseAsset(), cosmos.NewUint(1001*common.One))
	addr := GetRandomBaseAddress()
	acc, _ := addr.AccAddress()
	name := NewMAYAName("hello", 50, []MAYANameAlias{{Chain: common.BASEChain, Address: addr}})
	mgr.Keeper().SetMAYAName(ctx, name)

	// happy path
	msg := NewMsgManageMAYAName("I-am_the_99th_walrus+", common.BASEChain, addr, coin, 0, common.BNBAsset, acc, acc)
	c.Assert(handler.validate(ctx, *msg), IsNil)

	// fail: name is too long
	msg.Name = "this_name_is_way_too_long_to_be_a_valid_name"
	c.Assert(handler.validate(ctx, *msg), NotNil)

	// fail: bad characters
	msg.Name = "i am the walrus"
	c.Assert(handler.validate(ctx, *msg), NotNil)

	// fail: bad attempt to inflate expire block height
	msg.Name = "hello"
	msg.ExpireBlockHeight = 100
	c.Assert(handler.validate(ctx, *msg), NotNil)

	// fail: bad auth
	msg.ExpireBlockHeight = 0
	msg.Signer = GetRandomBech32Addr()
	c.Assert(handler.validate(ctx, *msg), NotNil)

	// fail: not enough funds for new MAYAName
	msg.Name = "bang"
	msg.Coin.Amount = cosmos.ZeroUint()
	c.Assert(handler.validate(ctx, *msg), NotNil)
}

func (s *HandlerManageMAYANameSuite) TestHandler(c *C) {
	ver := GetCurrentVersion()
	constAccessor := constants.GetConstantValues(ver)
	feePerBlock := constAccessor.GetInt64Value(constants.TNSFeePerBlock)
	registrationFee := constAccessor.GetInt64Value(constants.TNSRegisterFee)
	ctx, mgr := setupManagerForTest(c)

	blocksPerYear := mgr.GetConstants().GetInt64Value(constants.BlocksPerYear)
	handler := NewManageMAYANameHandler(mgr)
	coin := common.NewCoin(common.BaseAsset(), cosmos.NewUint(10000*common.One))
	addr := GetRandomBaseAddress()
	acc, _ := addr.AccAddress()
	tnName := "hello"

	// add rune to addr for gas
	FundAccount(c, ctx, mgr.Keeper(), acc, 10*common.One)

	// happy path, register new name
	msg := NewMsgManageMAYAName(tnName, common.BASEChain, addr, coin, 0, common.BaseAsset(), acc, acc)
	_, err := handler.handle(ctx, *msg)
	c.Assert(err, IsNil)
	name, err := mgr.Keeper().GetMAYAName(ctx, tnName)
	c.Assert(err, IsNil)
	c.Check(name.Owner.Empty(), Equals, false)
	c.Check(name.ExpireBlockHeight, Equals, ctx.BlockHeight()+blocksPerYear+(int64(coin.Amount.Uint64())-registrationFee)/feePerBlock)

	// happy path, set alt chain address
	bnbAddr := GetRandomBNBAddress()
	msg = NewMsgManageMAYAName(tnName, common.BNBChain, bnbAddr, coin, 0, common.BaseAsset(), acc, acc)
	_, err = handler.handle(ctx, *msg)
	c.Assert(err, IsNil)
	name, err = mgr.Keeper().GetMAYAName(ctx, tnName)
	c.Assert(err, IsNil)
	c.Check(name.GetAlias(common.BNBChain).Equals(bnbAddr), Equals, true)

	// happy path, update alt chain address
	bnbAddr = GetRandomBNBAddress()
	msg = NewMsgManageMAYAName(tnName, common.BNBChain, bnbAddr, coin, 0, common.BaseAsset(), acc, acc)
	_, err = handler.handle(ctx, *msg)
	c.Assert(err, IsNil)
	name, err = mgr.Keeper().GetMAYAName(ctx, tnName)
	c.Assert(err, IsNil)
	c.Check(name.GetAlias(common.BNBChain).Equals(bnbAddr), Equals, true)

	// happy path, release mayaname back into the wild
	msg = NewMsgManageMAYAName(tnName, common.BASEChain, addr, common.NewCoin(common.BaseAsset(), cosmos.ZeroUint()), 1, common.BaseAsset(), acc, acc)
	_, err = handler.handle(ctx, *msg)
	c.Assert(err, IsNil)
	name, err = mgr.Keeper().GetMAYAName(ctx, tnName)
	c.Assert(err, IsNil)
	c.Check(name.Owner.Empty(), Equals, true)
	c.Check(name.ExpireBlockHeight, Equals, int64(0))
}
