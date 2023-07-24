package chainclients

import (
	"github.com/rs/zerolog/log"
	"gitlab.com/thorchain/tss/go-tss/tss"

	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/avalanche"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/dash"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/dogecoin"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/gaia"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/thorchain"

	"gitlab.com/mayachain/mayanode/bifrost/mayaclient"
	"gitlab.com/mayachain/mayanode/bifrost/metrics"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/binance"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/bitcoin"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/bitcoincash"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/ethereum"
	"gitlab.com/mayachain/mayanode/bifrost/pkg/chainclients/litecoin"
	"gitlab.com/mayachain/mayanode/bifrost/pubkeymanager"
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/config"
)

// LoadChains returns chain clients from chain configuration
func LoadChains(thorKeys *mayaclient.Keys,
	cfg map[common.Chain]config.BifrostChainConfiguration,
	server *tss.TssServer,
	mayachainBridge mayaclient.MayachainBridge,
	m *metrics.Metrics,
	pubKeyValidator pubkeymanager.PubKeyValidator,
	poolMgr mayaclient.PoolManager,
) map[common.Chain]ChainClient {
	logger := log.Logger.With().Str("module", "bifrost").Logger()
	chains := make(map[common.Chain]ChainClient)

	for _, chain := range cfg {
		if chain.Disabled {
			logger.Info().Msgf("%s chain is disabled by configure", chain.ChainID)
			continue
		}
		switch chain.ChainID {
		case common.BNBChain:
			bnb, err := binance.NewBinance(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			chains[common.BNBChain] = bnb
		case common.ETHChain:
			eth, err := ethereum.NewClient(thorKeys, chain, server, mayachainBridge, m, pubKeyValidator, poolMgr)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			chains[common.ETHChain] = eth
		case common.BTCChain:
			btc, err := bitcoin.NewClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			pubKeyValidator.RegisterCallback(btc.RegisterPublicKey)
			chains[common.BTCChain] = btc
		case common.BCHChain:
			bch, err := bitcoincash.NewClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			pubKeyValidator.RegisterCallback(bch.RegisterPublicKey)
			chains[common.BCHChain] = bch
		case common.LTCChain:
			ltc, err := litecoin.NewClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			pubKeyValidator.RegisterCallback(ltc.RegisterPublicKey)
			chains[common.LTCChain] = ltc
		case common.DOGEChain:
			doge, err := dogecoin.NewClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			pubKeyValidator.RegisterCallback(doge.RegisterPublicKey)
			chains[common.DOGEChain] = doge
		case common.THORChain:
			logger.Debug().Msg("Loading THORCHAIN")
			thor, err := thorchain.NewCosmosClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			chains[common.THORChain] = thor
		case common.DASHChain:
			dashClient, err := dash.NewClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			pubKeyValidator.RegisterCallback(dashClient.RegisterPublicKey)
			chains[common.DASHChain] = dashClient
		case common.AVAXChain:
			avax, err := avalanche.NewAvalancheClient(thorKeys, chain, server, mayachainBridge, m, pubKeyValidator, poolMgr)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			chains[common.AVAXChain] = avax
		case common.GAIAChain:
			gaia, err := gaia.NewCosmosClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			chains[common.GAIAChain] = gaia
		case common.KUJIChain:
			kuji, err := gaia.NewCosmosClient(thorKeys, chain, server, mayachainBridge, m)
			if err != nil {
				logger.Fatal().Err(err).Str("chain_id", chain.ChainID.String()).Msg("fail to load chain")
				continue
			}
			chains[common.KUJIChain] = kuji
		}
	}

	return chains
}
