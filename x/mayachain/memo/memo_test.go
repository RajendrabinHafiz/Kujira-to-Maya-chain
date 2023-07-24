package mayachain

import (
	"fmt"
	"strings"
	"testing"

	. "gopkg.in/check.v1"

	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	kv1 "gitlab.com/mayachain/mayanode/x/mayachain/keeper/v1"
	"gitlab.com/mayachain/mayanode/x/mayachain/types"
)

type MemoSuite struct{}

func TestPackage(t *testing.T) { TestingT(t) }

var _ = Suite(&MemoSuite{})

func (s *MemoSuite) SetUpSuite(c *C) {
	types.SetupConfigForTest()
}

func (s *MemoSuite) TestTxType(c *C) {
	for _, trans := range []TxType{TxAdd, TxWithdraw, TxSwap, TxOutbound, TxDonate, TxBond, TxUnbond, TxLeave} {
		tx, err := StringToTxType(trans.String())
		c.Assert(err, IsNil)
		c.Check(tx, Equals, trans)
		c.Check(tx.IsEmpty(), Equals, false)
	}
}

func (s *MemoSuite) TestParseWithAbbreviated(c *C) {
	ctx := cosmos.Context{}
	k := kv1.KVStore{}
	k.SetVersion(kv1.GetCurrentVersion())

	// happy paths
	memo, err := ParseMemoWithMAYANames(ctx, k, "d:"+common.BaseAsset().String())
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxDonate), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.IsInbound(), Equals, true)
	c.Check(memo.IsInternal(), Equals, false)
	c.Check(memo.IsOutbound(), Equals, false)

	memo, err = ParseMemoWithMAYANames(ctx, k, "+:"+common.BaseAsset().String())
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxAdd), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.IsInbound(), Equals, true)
	c.Check(memo.IsInternal(), Equals, false)
	c.Check(memo.IsOutbound(), Equals, false)

	_, err = ParseMemoWithMAYANames(ctx, k, "add:BTC.BTC:tbnb1yeuljgpkg2c2qvx3nlmgv7gvnyss6ye2u8rasf:xxxx")
	c.Assert(err, IsNil)

	memo, err = ParseMemoWithMAYANames(ctx, k, fmt.Sprintf("-:%s:25", common.BaseAsset().String()))
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxWithdraw), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetAmount().Uint64(), Equals, uint64(25), Commentf("%d", memo.GetAmount().Uint64()))
	c.Check(memo.IsInbound(), Equals, true)
	c.Check(memo.IsInternal(), Equals, false)
	c.Check(memo.IsOutbound(), Equals, false)

	memo, err = ParseMemoWithMAYANames(ctx, k, "=:"+common.BaseAsset().String()+":bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6:870000000")
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxSwap), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetDestination().String(), Equals, "bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Log(memo.GetSlipLimit().Uint64())
	c.Check(memo.GetSlipLimit().Equal(cosmos.NewUint(870000000)), Equals, true)
	c.Check(memo.IsInbound(), Equals, true)
	c.Check(memo.IsInternal(), Equals, false)
	c.Check(memo.IsOutbound(), Equals, false)

	memo, err = ParseMemoWithMAYANames(ctx, k, "=:"+common.BaseAsset().String()+":bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxSwap), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetDestination().String(), Equals, "bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(memo.GetSlipLimit().Uint64(), Equals, uint64(0))
	c.Check(memo.IsInbound(), Equals, true)

	memo, err = ParseMemoWithMAYANames(ctx, k, "=:"+common.BaseAsset().String()+":bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6:")
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxSwap), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetDestination().String(), Equals, "bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(memo.GetSlipLimit().Equal(cosmos.ZeroUint()), Equals, true)

	memo, err = ParseMemoWithMAYANames(ctx, k, "OUT:MUKVQILIHIAUSEOVAXBFEZAJKYHFJYHRUUYGQJZGFYBYVXCXYNEMUOAIQKFQLLCX")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxOutbound), Equals, true, Commentf("%s", memo.GetType()))
	c.Check(memo.IsOutbound(), Equals, true)
	c.Check(memo.IsInbound(), Equals, false)
	c.Check(memo.IsInternal(), Equals, false)

	memo, err = ParseMemoWithMAYANames(ctx, k, "REFUND:MUKVQILIHIAUSEOVAXBFEZAJKYHFJYHRUUYGQJZGFYBYVXCXYNEMUOAIQKFQLLCX")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxRefund), Equals, true)
	c.Check(memo.IsOutbound(), Equals, true)

	memo, err = ParseMemoWithMAYANames(ctx, k, "leave:whatever")
	c.Assert(err, NotNil)
	c.Check(memo.IsType(TxUnknown), Equals, true)

	addr := types.GetRandomBech32Addr()
	memo, err = ParseMemoWithMAYANames(ctx, k, fmt.Sprintf("leave:%s", addr.String()))
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxLeave), Equals, true)
	c.Check(memo.GetAccAddress().String(), Equals, addr.String())

	memo, err = ParseMemoWithMAYANames(ctx, k, "yggdrasil+:30")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxYggdrasilFund), Equals, true)
	c.Check(memo.IsInbound(), Equals, false)
	c.Check(memo.IsInternal(), Equals, true)
	memo, err = ParseMemoWithMAYANames(ctx, k, "yggdrasil-:30")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxYggdrasilReturn), Equals, true)
	c.Check(memo.IsInternal(), Equals, true)
	memo, err = ParseMemoWithMAYANames(ctx, k, "migrate:100")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxMigrate), Equals, true)
	c.Check(memo.IsInternal(), Equals, true)

	memo, err = ParseMemoWithMAYANames(ctx, k, "ragnarok:100")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxRagnarok), Equals, true)
	c.Check(memo.IsOutbound(), Equals, true)

	memo, err = ParseMemoWithMAYANames(ctx, k, "reserve")
	c.Check(err, IsNil)
	c.Check(memo.IsType(TxReserve), Equals, true)
	c.Check(memo.IsInbound(), Equals, true)
	c.Check(memo.IsInternal(), Equals, false)
	c.Check(memo.IsOutbound(), Equals, false)

	memo, err = ParseMemoWithMAYANames(ctx, k, "noop")
	c.Check(err, IsNil)
	c.Check(memo.IsType(TxNoOp), Equals, true)
	c.Check(memo.IsInbound(), Equals, true)
	c.Check(memo.IsInternal(), Equals, false)
	c.Check(memo.IsOutbound(), Equals, false)

	memo, err = ParseMemoWithMAYANames(ctx, k, "noop:novault")
	c.Check(err, IsNil)
	c.Check(memo.IsType(TxNoOp), Equals, true)
	c.Check(memo.IsInbound(), Equals, true)
	c.Check(memo.IsInternal(), Equals, false)
	c.Check(memo.IsOutbound(), Equals, false)

	addr1 := types.GetRandomBaseAddress()
	addr2 := types.GetRandomBaseAddress()
	memo, err = ParseMemoWithMAYANames(ctx, k, fmt.Sprintf("unbond:%s:%s", addr1.String(), addr2.String()))
	c.Check(err, IsNil)
	c.Check(memo.IsType(TxUnbond), Equals, true)
	c.Check(memo.GetAsset().String(), Equals, common.EmptyAsset.String())
	c.Check(memo.GetAccAddress().String(), Equals, addr1.String())

	// unhappy paths
	_, err = ParseMemoWithMAYANames(ctx, k, "")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "bogus")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "CREATE") // missing symbol
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "c:") // bad symbol
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "-:bnb") // withdraw basis points is optional
	c.Assert(err, IsNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "-:bnb:twenty-two") // bad amount
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "=:bnb:bad_DES:5.6") // bad destination
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, ">:bnb:bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6:five") // bad slip limit
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "!:key:val") // not enough arguments
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "!:bogus:key:value") // bogus admin command type
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "nextpool:whatever")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "migrate")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "switch")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "switch:")
	c.Assert(err, NotNil)
}

func (s *MemoSuite) TestParse(c *C) {
	ctx := cosmos.Context{}
	k := kv1.KVStore{}
	k.SetVersion(types.GetCurrentVersion())

	// happy paths
	memo, err := ParseMemoWithMAYANames(ctx, k, "d:"+common.BaseAsset().String())
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxDonate), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.String(), Equals, "DONATE:"+common.BaseAsset().String())

	memo, err = ParseMemoWithMAYANames(ctx, k, "ADD:"+common.BaseAsset().String())
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxAdd), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.String(), Equals, "")

	_, err = ParseMemoWithMAYANames(ctx, k, "ADD:BTC.BTC")
	c.Assert(err, IsNil)
	memo, err = ParseMemoWithMAYANames(ctx, k, "ADD:BTC.BTC:bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej")
	c.Assert(err, IsNil)
	c.Check(memo.GetDestination().String(), Equals, "bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej")
	c.Check(memo.IsType(TxAdd), Equals, true, Commentf("MEMO: %+v", memo))

	_, err = ParseMemoWithMAYANames(ctx, k, "ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:tmaya176xrckly4p7efq7fshhcuc2kax3dyxu9hlzwfw:1000")
	c.Assert(err, IsNil)

	memo, err = ParseMemoWithMAYANames(ctx, k, "WITHDRAW:"+common.BaseAsset().String()+":25")
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxWithdraw), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetAmount().Equal(cosmos.NewUint(25)), Equals, true, Commentf("%d", memo.GetAmount().Uint64()))

	memo, err = ParseMemoWithMAYANames(ctx, k, "SWAP:"+common.BaseAsset().String()+":bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6:870000000")
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxSwap), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetDestination().String(), Equals, "bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Log(memo.GetSlipLimit().String())
	c.Check(memo.GetSlipLimit().Equal(cosmos.NewUint(870000000)), Equals, true)

	memo, err = ParseMemoWithMAYANames(ctx, k, "SWAP:"+common.BaseAsset().String()+":bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxSwap), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetDestination().String(), Equals, "bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(memo.GetSlipLimit().Uint64(), Equals, uint64(0))

	memo, err = ParseMemoWithMAYANames(ctx, k, "SWAP:"+common.BaseAsset().String()+":bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6:")
	c.Assert(err, IsNil)
	c.Check(memo.GetAsset().String(), Equals, common.BaseAsset().String())
	c.Check(memo.IsType(TxSwap), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetDestination().String(), Equals, "bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6")
	c.Check(memo.GetSlipLimit().Uint64(), Equals, uint64(0))

	whiteListAddr := types.GetRandomBech32Addr()
	bondProvider := types.GetRandomBech32Addr()
	memo, err = ParseMemoWithMAYANames(ctx, k, fmt.Sprintf("BOND:%s:1000:%s:%s", common.BNBAsset.String(), whiteListAddr, bondProvider))
	c.Assert(err, IsNil)
	c.Assert(memo.IsType(TxBond), Equals, true)
	c.Assert(memo.GetAsset().String(), Equals, common.BNBAsset.String())
	c.Assert(memo.GetAmount().Equal(cosmos.NewUint(1000)), Equals, true)
	c.Assert(memo.GetAccAddress().String(), Equals, whiteListAddr.String())
	mem, err := ParseBondMemo(k.GetVersion(), []string{"BOND", common.BTCAsset.String(), "1000", whiteListAddr.String(), bondProvider.String()})
	c.Assert(err, IsNil)
	c.Assert(mem.Asset.String(), Equals, common.BTCAsset.String())
	c.Assert(mem.Units.Equal(cosmos.NewUint(1000)), Equals, true)
	c.Assert(mem.NodeAddress.String(), Equals, whiteListAddr.String())
	c.Assert(mem.BondProviderAddress.String(), Equals, bondProvider.String())
	c.Assert(mem.NodeOperatorFee, Equals, int64(-1))
	// Bond as invite
	mem, err = ParseBondMemo(k.GetVersion(), []string{"BOND", whiteListAddr.String(), bondProvider.String()})
	c.Assert(err, IsNil)
	c.Assert(mem.Asset.String(), Equals, common.EmptyAsset.String())
	c.Assert(mem.Units.Equal(cosmos.ZeroUint()), Equals, true)
	c.Assert(mem.NodeAddress.String(), Equals, whiteListAddr.String())
	c.Assert(mem.BondProviderAddress.String(), Equals, bondProvider.String())
	c.Assert(mem.NodeOperatorFee, Equals, int64(-1))
	mem, err = ParseBondMemo(k.GetVersion(), []string{"BOND", whiteListAddr.String(), bondProvider.String(), "1000"})
	c.Assert(err, IsNil)
	c.Assert(mem.Asset.String(), Equals, common.EmptyAsset.String())
	c.Assert(mem.Units.Equal(cosmos.NewUint(0)), Equals, true)
	c.Assert(mem.BondProviderAddress.String(), Equals, bondProvider.String())
	c.Assert(mem.NodeOperatorFee, Equals, int64(1000))
	mem, err = ParseBondMemo(k.GetVersion(), []string{"BOND", common.ETHAsset.String(), "1000", whiteListAddr.String(), bondProvider.String(), "1000"})
	c.Assert(err, IsNil)
	c.Assert(mem.BondProviderAddress.String(), Equals, bondProvider.String())
	c.Assert(mem.NodeOperatorFee, Equals, int64(1000))

	memo, err = ParseMemoWithMAYANames(ctx, k, "leave:"+types.GetRandomBech32Addr().String())
	c.Assert(err, IsNil)
	c.Assert(memo.IsType(TxLeave), Equals, true)

	memo, err = ParseMemoWithMAYANames(ctx, k, "unbond:"+whiteListAddr.String())
	c.Assert(err, IsNil)
	c.Assert(memo.IsType(TxUnbond), Equals, true)
	c.Assert(memo.GetAccAddress().String(), Equals, whiteListAddr.String())
	unbondMemo, err := ParseUnbondMemo(k.GetVersion(), []string{"UNBOND", whiteListAddr.String(), bondProvider.String()})
	c.Assert(err, IsNil)
	c.Assert(unbondMemo.BondProviderAddress.String(), Equals, bondProvider.String())

	memo, err = ParseMemoWithMAYANames(ctx, k, "migrate:100")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxMigrate), Equals, true)
	c.Check(memo.GetBlockHeight(), Equals, int64(100))
	c.Check(memo.String(), Equals, "MIGRATE:100")

	txID := types.GetRandomTxHash()
	memo, err = ParseMemoWithMAYANames(ctx, k, "OUT:"+txID.String())
	c.Check(err, IsNil)
	c.Check(memo.IsOutbound(), Equals, true)
	c.Check(memo.GetTxID(), Equals, txID)
	c.Check(memo.String(), Equals, "OUT:"+txID.String())

	refundMemo := "REFUND:" + txID.String()
	memo, err = ParseMemoWithMAYANames(ctx, k, refundMemo)
	c.Check(err, IsNil)
	c.Check(memo.GetTxID(), Equals, txID)
	c.Check(memo.String(), Equals, refundMemo)

	yggFundMemo := "YGGDRASIL+:100"
	memo, err = ParseMemoWithMAYANames(ctx, k, yggFundMemo)
	c.Check(err, IsNil)
	c.Check(memo.GetBlockHeight(), Equals, int64(100))
	c.Check(memo.String(), Equals, yggFundMemo)

	yggReturnMemo := "YGGDRASIL-:100"
	memo, err = ParseMemoWithMAYANames(ctx, k, yggReturnMemo)
	c.Check(err, IsNil)
	c.Check(memo.GetBlockHeight(), Equals, int64(100))
	c.Check(memo.String(), Equals, yggReturnMemo)

	ragnarokMemo := "RAGNAROK:1024"
	memo, err = ParseMemoWithMAYANames(ctx, k, ragnarokMemo)
	c.Check(err, IsNil)
	c.Check(memo.IsType(TxRagnarok), Equals, true)
	c.Check(memo.GetBlockHeight(), Equals, int64(1024))
	c.Check(memo.String(), Equals, ragnarokMemo)

	baseMemo := MemoBase{}
	c.Check(baseMemo.String(), Equals, "")
	c.Check(baseMemo.GetAmount().Uint64(), Equals, cosmos.ZeroUint().Uint64())
	c.Check(baseMemo.GetDestination(), Equals, common.NoAddress)
	c.Check(baseMemo.GetSlipLimit().Uint64(), Equals, cosmos.ZeroUint().Uint64())
	c.Check(baseMemo.GetTxID(), Equals, common.TxID(""))
	c.Check(baseMemo.GetAccAddress().Empty(), Equals, true)
	c.Check(baseMemo.IsEmpty(), Equals, true)
	c.Check(baseMemo.GetBlockHeight(), Equals, int64(0))

	// unhappy paths
	_, err = ParseMemoWithMAYANames(ctx, k, "")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "bogus")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "CREATE") // missing symbol
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "CREATE:") // bad symbol
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "withdraw") // not enough parameters
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "withdraw:bnb") // withdraw basis points is optional
	c.Assert(err, IsNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "withdraw:bnb:twenty-two") // bad amount
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "swap") // not enough parameters
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "swap:bnb:PROVIDER-1:5.6") // bad destination
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "swap:bnb:bad_DES:5.6") // bad destination
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "swap:bnb:bnb1lejrrtta9cgr49fuh7ktu3sddhe0ff7wenlpn6:five") // bad slip limit
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "admin:key:val") // not enough arguments
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "admin:bogus:key:value") // bogus admin command type
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "migrate:abc")
	c.Assert(err, NotNil)

	_, err = ParseMemoWithMAYANames(ctx, k, "withdraw:A")
	c.Assert(err, IsNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "leave")
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "out") // not enough parameter
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "bond") // not enough parameter
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "refund") // not enough parameter
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "yggdrasil+") // not enough parameter
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "yggdrasil+:A") // invalid block height
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "yggdrasil-") // not enough parameter
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "yggdrasil-:B") // invalid block height
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "ragnarok") // not enough parameter
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "ragnarok:what") // not enough parameter
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "bond:what") // invalid address
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "switch:what") // invalid address
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "whatever") // not support
	c.Assert(err, NotNil)
	_, err = ParseMemoWithMAYANames(ctx, k, "unbond") // not enough parameter
	c.Assert(err, NotNil)
}

func (s *MemoSuite) TestParseWithdrawPairAddress(c *C) {
	ctx := cosmos.Context{}

	memo, err := ParseMemoWithMAYANames(ctx, nil, "withdraw:"+common.BNBAsset.String()+":25:"+common.BNBAsset.String()+":tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxWithdraw), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetAsset().String(), Equals, common.BNBAsset.String())
	c.Check(memo.GetAmount().Equal(cosmos.NewUint(25)), Equals, true, Commentf("%d", memo.GetAmount().Uint64()))

	memo, err = ParseMemoWithMAYANames(ctx, nil, "-:"+common.BTCAsset.String()+":9999:"+common.BTCAsset.String()+":bc1qwqdg6squsna38e46795at95yu9atm8azzmyvckulcc7kytlcckxswvvzej")
	c.Assert(err, IsNil)
	c.Check(memo.IsType(TxWithdraw), Equals, true, Commentf("MEMO: %+v", memo))
	c.Check(memo.GetAsset().String(), Equals, common.BTCAsset.String())
	c.Check(memo.GetAmount().Equal(cosmos.NewUint(9999)), Equals, true, Commentf("%d", memo.GetAmount().Uint64()))

	_, err = ParseMemoWithMAYANames(ctx, nil, "-:"+common.BTCAsset.String()+":10000:"+common.BTCAsset.String()+":wrongaddress")
	c.Assert(err.Error(), Equals, "address format not supported: wrongaddress")
}

func (s *MemoSuite) TestParseAddLiquidityMemoTier(c *C) {
	ctx := cosmos.Context{}
	k := kv1.KVStore{}
	k.SetVersion(types.GetCurrentVersion())

	asset, err := common.NewAsset("BNB.BNB")
	c.Assert(err, IsNil)

	// happy paths
	parts := strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:tmaya176xrckly4p7efq7fshhcuc2kax3dyxu9hlzwfw:1000:TIER1", ":")
	addMemo, err := ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(1))

	parts = strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:TIER1", ":")
	addMemo, err = ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(1))

	parts = strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:TIER2", ":")
	addMemo, err = ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(2))

	parts = strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:TIER3", ":")
	addMemo, err = ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(3))

	// if tear is not between 1 and 3 or is not set, tier by default is set to 3
	parts = strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:tmaya176xrckly4p7efq7fshhcuc2kax3dyxu9hlzwfw:1000:TIER4", ":")
	addMemo, err = ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(3))

	parts = strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:TIER4", ":")
	addMemo, err = ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(3))

	parts = strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll:tmaya176xrckly4p7efq7fshhcuc2kax3dyxu9hlzwfw:1000", ":")
	addMemo, err = ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(3))

	parts = strings.Split("ADD:BNB.BNB:tbnb18f55frcvknxvcpx2vvpfedvw4l8eutuhca3lll", ":")
	addMemo, err = ParseAddLiquidityMemo(ctx, k, asset, parts)
	c.Assert(err, IsNil)
	c.Assert(addMemo.Tier, Equals, int64(3))
}
