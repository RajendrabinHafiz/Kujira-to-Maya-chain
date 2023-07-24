package keeperv1

import (
	"gitlab.com/mayachain/mayanode/common"
	"gitlab.com/mayachain/mayanode/x/mayachain/types"
)

const (
	ModuleName  = types.ModuleName
	ReserveName = types.ReserveName
	AsgardName  = types.AsgardName
	BondName    = types.BondName
	StoreKey    = types.StoreKey
	MayaFund    = types.MayaFund

	// Vaults
	AsgardVault    = types.VaultType_AsgardVault
	YggdrasilVault = types.VaultType_YggdrasilVault
	ActiveVault    = types.VaultStatus_ActiveVault
	InitVault      = types.VaultStatus_InitVault
	RetiringVault  = types.VaultStatus_RetiringVault
	InactiveVault  = types.VaultStatus_InactiveVault

	// Node status
	NodeActive  = types.NodeStatus_Active
	NodeStandby = types.NodeStatus_Standby
	NodeUnknown = types.NodeStatus_Unknown

	// Node type
	NodeTypeUnknown   = types.NodeType_TypeUnknown
	NodeTypeValidator = types.NodeType_TypeValidator
	NodeTypeVault     = types.NodeType_TypeVault

	// Bond type
	AsgardKeygen = types.KeygenType_AsgardKeygen

	// Pool Status
	PoolAvailable = types.PoolStatus_Available
	PoolStaged    = types.PoolStatus_Staged
	PoolSuspended = types.PoolStatus_Suspended
	AllPoolStatus = PoolAvailable & PoolStaged & PoolSuspended
)

var (
	NewPool                    = types.NewPool
	NewJail                    = types.NewJail
	NewNetwork                 = types.NewNetwork
	NewProtocolOwnedLiquidity  = types.NewProtocolOwnedLiquidity
	NewObservedTx              = types.NewObservedTx
	NewTssVoter                = types.NewTssVoter
	NewBanVoter                = types.NewBanVoter
	NewForgiveSlashVoter       = types.NewForgiveSlashVoter
	NewErrataTxVoter           = types.NewErrataTxVoter
	NewObservedTxVoter         = types.NewObservedTxVoter
	NewKeygen                  = types.NewKeygen
	NewKeygenBlock             = types.NewKeygenBlock
	NewTxOut                   = types.NewTxOut
	HasSuperMajority           = types.HasSuperMajority
	RegisterCodec              = types.RegisterCodec
	NewNodeAccount             = types.NewNodeAccount
	NewBondProviders           = types.NewBondProviders
	NewBondProvider            = types.NewBondProvider
	NewVault                   = types.NewVault
	NewReserveContributor      = types.NewReserveContributor
	NewMAYAName                = types.NewMAYAName
	GetRandomTx                = types.GetRandomTx
	GetRandomValidatorNode     = types.GetRandomValidatorNode
	GetRandomVaultNode         = types.GetRandomVaultNode
	GetRandomBNBAddress        = types.GetRandomBNBAddress
	GetRandomBTCAddress        = types.GetRandomBTCAddress
	GetRandomBCHAddress        = types.GetRandomBCHAddress
	GetRandomBaseAddress       = types.GetRandomBaseAddress
	GetRandomTxHash            = types.GetRandomTxHash
	GetRandomBech32Addr        = types.GetRandomBech32Addr
	GetRandomPubKey            = types.GetRandomPubKey
	GetRandomPubKeySet         = types.GetRandomPubKeySet
	GetCurrentVersion          = types.GetCurrentVersion
	NewObservedNetworkFeeVoter = types.NewObservedNetworkFeeVoter
	NewNetworkFee              = types.NewNetworkFee
	NewTssKeysignFailVoter     = types.NewTssKeysignFailVoter
	SetupConfigForTest         = types.SetupConfigForTest
	NewChainContract           = types.NewChainContract
	GetLiquidityPools          = types.GetLiquidityPools
)

type (
	MsgSwap                  = types.MsgSwap
	Pool                     = types.Pool
	Pools                    = types.Pools
	LiquidityProvider        = types.LiquidityProvider
	LiquidityProviders       = types.LiquidityProviders
	ObservedTxs              = types.ObservedTxs
	ObservedTxVoter          = types.ObservedTxVoter
	BanVoter                 = types.BanVoter
	ForgiveSlashVoter        = types.ForgiveSlashVoter
	ErrataTxVoter            = types.ErrataTxVoter
	TssVoter                 = types.TssVoter
	TssKeysignFailVoter      = types.TssKeysignFailVoter
	TxOutItem                = types.TxOutItem
	TxOut                    = types.TxOut
	KeygenBlock              = types.KeygenBlock
	ReserveContributors      = types.ReserveContributors
	Vault                    = types.Vault
	Vaults                   = types.Vaults
	Jail                     = types.Jail
	BondProvider             = types.BondProvider
	BondProviders            = types.BondProviders
	NodeAccount              = types.NodeAccount
	NodeAccounts             = types.NodeAccounts
	NodeStatus               = types.NodeStatus
	NodeType                 = types.NodeType
	Network                  = types.Network
	VaultStatus              = types.VaultStatus
	NetworkFee               = types.NetworkFee
	ObservedNetworkFeeVoter  = types.ObservedNetworkFeeVoter
	RagnarokWithdrawPosition = types.RagnarokWithdrawPosition
	TssKeygenMetric          = types.TssKeygenMetric
	TssKeysignMetric         = types.TssKeysignMetric
	ChainContract            = types.ChainContract
	MAYAName                 = types.MAYAName
	MAYANameAlias            = types.MAYANameAlias
	SolvencyVoter            = types.SolvencyVoter
	NodeMimir                = types.NodeMimir
	NodeMimirs               = types.NodeMimirs
	LiquidityAuctionTier     = types.LiquidityAuctionTier
	ProtocolOwnedLiquidity   = types.ProtocolOwnedLiquidity

	ProtoInt64        = types.ProtoInt64
	ProtoUint64       = types.ProtoUint64
	ProtoAccAddresses = types.ProtoAccAddresses
	ProtoStrings      = types.ProtoStrings
	ProtoUint         = common.ProtoUint
	ProtoBools        = types.ProtoBools
)
