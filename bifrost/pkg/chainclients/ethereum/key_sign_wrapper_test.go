package ethereum

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	ecommon "github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
	. "gopkg.in/check.v1"

	"gitlab.com/mayachain/mayanode/bifrost/mayaclient"
	"gitlab.com/mayachain/mayanode/bifrost/tss"
	"gitlab.com/mayachain/mayanode/cmd"
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/config"
	"gitlab.com/mayachain/mayanode/x/mayachain/types"
)

type ETHKeysignWrapperTestSuite struct {
	thorKeys *mayaclient.Keys
	wrapper  *keySignWrapper
}

var _ = Suite(
	&ETHKeysignWrapperTestSuite{},
)

// SetUpSuite setup the test conditions
func (s *ETHKeysignWrapperTestSuite) SetUpSuite(c *C) {
	cfg := config.BifrostClientConfiguration{
		ChainID:      "thorchain",
		SignerName:   "bob",
		SignerPasswd: "password",
	}

	kb := cKeys.NewInMemory()
	_, _, err := kb.NewMnemonic(cfg.SignerName, cKeys.English, cmd.BASEChainHDPath, cfg.SignerPasswd, hd.Secp256k1)
	c.Assert(err, IsNil)
	s.thorKeys = mayaclient.NewKeysWithKeybase(kb, cfg.SignerName, cfg.SignerPasswd)

	privateKey, err := s.thorKeys.GetPrivateKey()
	c.Assert(err, IsNil)
	temp, err := codec.ToTmPubKeyInterface(privateKey.PubKey())
	c.Assert(err, IsNil)
	pk, err := common.NewPubKeyFromCrypto(temp)
	c.Assert(err, IsNil)
	keyMgr := &tss.MockMayachainKeyManager{}
	ethPrivateKey, err := getETHPrivateKey(privateKey)
	c.Assert(err, IsNil)
	c.Assert(ethPrivateKey, NotNil)
	wrapper, err := newKeySignWrapper(ethPrivateKey, pk, keyMgr, big.NewInt(15))
	c.Assert(err, IsNil)
	c.Assert(wrapper, NotNil)
	s.wrapper = wrapper
}

func (s *ETHKeysignWrapperTestSuite) TestGetPrivKey(c *C) {
	c.Assert(s.wrapper.GetPrivKey(), NotNil)
}

func (s *ETHKeysignWrapperTestSuite) TestGetPubKey(c *C) {
	c.Assert(s.wrapper.GetPubKey(), NotNil)
}

func (s *ETHKeysignWrapperTestSuite) TestSign(c *C) {
	buf, err := s.wrapper.Sign(nil, types.GetRandomPubKey())
	c.Assert(err, NotNil)
	c.Assert(buf, IsNil)
	createdTx := etypes.NewTransaction(0, ecommon.HexToAddress("0x7d182d6a138eaa06f6f452bc3f8fc57e17d1e193"), big.NewInt(1), MaxContractGas, big.NewInt(1), []byte("whatever"))
	buf, err = s.wrapper.Sign(createdTx, common.EmptyPubKey)
	c.Assert(err, NotNil)
	c.Assert(buf, IsNil)
	_, err = s.wrapper.Sign(createdTx, s.wrapper.pubKey)
	c.Assert(err, IsNil)
	// test sign with TSS
	buf, err = s.wrapper.Sign(createdTx, types.GetRandomPubKey())
	c.Assert(err, NotNil)
	c.Assert(buf, IsNil)
}
