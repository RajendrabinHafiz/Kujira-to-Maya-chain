package thorchain

import (
	"crypto/x509"
	"fmt"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	ctypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/common/cosmos"
	tcTypes "gitlab.com/mayachain/mayanode/x/mayachain/types"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultMaxMsgSize = 1024 * 1024 * 1024 // 1GB

// buildUnsigned takes a MsgSend and other parameters and returns a txBuilder
// It can be used to simulateTx or as the input to signMsg before BraodcastTx
func buildUnsigned(
	txConfig client.TxConfig,
	msg *tcTypes.MsgSend,
	pubkey common.PubKey,
	memo string,
	fee ctypes.Coins,
	account uint64,
	sequence uint64,
) (client.TxBuilder, error) {
	cpk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, pubkey.String())
	if err != nil {
		return nil, fmt.Errorf("unable to GetPubKeyFromBech32 from cosmos: %w", err)
	}
	txBuilder := txConfig.NewTxBuilder()

	err = txBuilder.SetMsgs(msg)
	if err != nil {
		return nil, fmt.Errorf("unable to SetMsgs on txBuilder: %w", err)
	}

	txBuilder.SetMemo(memo)
	// Currently no fees set for THORChain
	// txBuilder.SetFeeAmount(fee)
	txBuilder.SetGasLimit(GasLimit)

	sigData := &signingtypes.SingleSignatureData{
		SignMode: signingtypes.SignMode_SIGN_MODE_DIRECT,
	}
	sig := signingtypes.SignatureV2{
		PubKey:   cpk,
		Data:     sigData,
		Sequence: sequence,
	}

	err = txBuilder.SetSignatures(sig)
	if err != nil {
		return nil, fmt.Errorf("unable to initial SetSignatures on txBuilder: %w", err)
	}

	return txBuilder, nil
}

func fromCosmosToThorchain(c cosmos.Coin) (common.Coin, error) {
	cosmosAsset, exists := GetAssetByCosmosDenom(c.Denom)
	if !exists {
		return common.NoCoin, fmt.Errorf("asset does not exist / not whitelisted by client")
	}

	thorAsset, err := common.NewAsset(fmt.Sprintf("%s.%s", common.THORChain.String(), cosmosAsset.THORChainSymbol))
	if err != nil {
		return common.NoCoin, fmt.Errorf("invalid thorchain asset: %w", err)
	}

	return common.Coin{
		Asset:    thorAsset,
		Amount:   cosmos.NewUint(c.Amount.Uint64()),
		Decimals: int64(cosmosAsset.CosmosDecimals),
	}, nil
}

func fromThorchainToCosmos(coin common.Coin) (cosmos.Coin, error) {
	asset, exists := GetAssetByThorchainSymbol(coin.Asset.Symbol.String())
	if !exists {
		return cosmos.Coin{}, fmt.Errorf("asset does not exist / not whitelisted by client")
	}

	amount := coin.Amount.BigInt()
	return cosmos.NewCoin(asset.CosmosDenom, ctypes.NewIntFromBigInt(amount)), nil
}

func parseEnvMaxMsgSize(envVar string) (int, error) {
	maxMsgSizeEnv := os.Getenv(envVar)
	if maxMsgSizeEnv == "" {
		return defaultMaxMsgSize, nil
	}

	parsed, err := strconv.Atoi(maxMsgSizeEnv)
	if err != nil {
		return defaultMaxMsgSize, err
	}

	return parsed, nil
}

// Bifrost only supports an "RPCHost" in its configuration.
// We also need to access GRPC for Cosmos chains
func getGRPCConn(host string, tls bool) (*grpc.ClientConn, error) {
	// load system certificates or proceed with insecure if tls disabled
	var creds credentials.TransportCredentials
	if tls {
		certs, err := x509.SystemCertPool()
		if err != nil {
			return &grpc.ClientConn{}, fmt.Errorf("unable to load system certs: %w", err)
		}
		creds = credentials.NewClientTLSFromCert(certs, "")
	} else {
		creds = insecure.NewCredentials()
	}

	maxMsgSize, err := parseEnvMaxMsgSize("THOR_GRPC_MAX_MSG_SIZE")
	if err != nil {
		return &grpc.ClientConn{}, fmt.Errorf("unable to parse THOR_GRPC_MAX_MSG_SIZE: %w", err)
	}

	return grpc.Dial(host, grpc.WithTransportCredentials(creds), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
}

func unmarshalJSONToPb(filePath string, msg proto.Message) error {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	u := new(jsonpb.Unmarshaler)
	u.AllowUnknownFields = true
	return u.Unmarshal(jsonFile, msg)
}

// TODO Consider implement optimized solution
// github.com/cosmos/cosmos-sdk/pull/8717
// github.com/cosmos/cosmos-sdk/pull/8694
func accAddressToString(acc ctypes.AccAddress, prefix string) (string, error) {
	if acc.Empty() {
		return "", fmt.Errorf("empty account address")
	}

	bech32Addr, err := bech32.ConvertAndEncode(prefix, acc)
	if err != nil {
		return "", fmt.Errorf("unable to encode AccAddress to bech32: %w", err)
	}

	return bech32Addr, nil
}

func accAddressFromBech32(address common.Address) (ctypes.AccAddress, error) {
	if address.IsEmpty() {
		return nil, fmt.Errorf("empty address %s", address.String())
	}

	hrp, bz, err := bech32.DecodeAndConvert(address.String())
	if err != nil {
		return nil, err
	}

	if hrp != "thor" && hrp != "tthor" {
		return nil, fmt.Errorf("invalid address prefix %s", hrp)
	}

	return ctypes.AccAddress(bz), nil
}
