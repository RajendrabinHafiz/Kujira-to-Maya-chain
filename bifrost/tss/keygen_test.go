package tss

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	. "gopkg.in/check.v1"

	"gitlab.com/mayachain/mayanode/bifrost/mayaclient"
	"gitlab.com/mayachain/mayanode/cmd"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	"gitlab.com/mayachain/mayanode/config"
	"gitlab.com/mayachain/mayanode/x/mayachain"
)

func TestTSSKeyGen(t *testing.T) { TestingT(t) }

type KeyGenTestSuite struct{}

var _ = Suite(&KeyGenTestSuite{})

func (*KeyGenTestSuite) SetUpSuite(c *C) {
	mayachain.SetupConfigForTest()
}

const (
	signerNameForTest     = `jack`
	signerPasswordForTest = `password`
)

func (*KeyGenTestSuite) setupKeysForTest(c *C) string {
	ns := strconv.Itoa(time.Now().Nanosecond())
	thorcliDir := filepath.Join(os.TempDir(), ns, ".thorcli")
	c.Logf("thorcliDir:%s", thorcliDir)
	buf := bytes.NewBufferString(signerPasswordForTest)
	// the library used by keyring is using ReadLine , which expect a new line
	buf.WriteByte('\n')
	buf.WriteString(signerPasswordForTest)
	buf.WriteByte('\n')
	kb, err := cKeys.New(cosmos.KeyringServiceName(), cKeys.BackendFile, thorcliDir, buf)
	c.Assert(err, IsNil)
	info, _, err := kb.NewMnemonic(signerNameForTest, cKeys.English, cmd.BASEChainHDPath, signerPasswordForTest, hd.Secp256k1)
	c.Assert(err, IsNil)
	c.Logf("name:%s", info.GetName())
	return thorcliDir
}

func (kts *KeyGenTestSuite) TestNewTssKenGen(c *C) {
	oldStdIn := os.Stdin
	defer func() {
		os.Stdin = oldStdIn
	}()
	os.Stdin = nil
	folder := kts.setupKeysForTest(c)
	defer func() {
		err := os.RemoveAll(folder)
		c.Assert(err, IsNil)
	}()
	kb, _, err := mayaclient.GetKeyringKeybase(folder, signerNameForTest, signerPasswordForTest)
	c.Assert(err, IsNil)
	k := mayaclient.NewKeysWithKeybase(kb, signerNameForTest, signerPasswordForTest)
	c.Assert(k, NotNil)
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		c.Logf("requestUri:%s", req.RequestURI)
	}))
	b, err := mayaclient.NewMayachainBridge(config.BifrostClientConfiguration{
		ChainID:      "mayachain",
		ChainHost:    server.Listener.Addr().String(),
		SignerName:   "bob",
		SignerPasswd: "password",
	}, nil, k)
	c.Assert(err, IsNil)
	kg, err := NewTssKeyGen(k, nil, b)
	c.Assert(err, IsNil)
	c.Assert(kg, NotNil)
}
