package constants

import (
	"fmt"

	"github.com/blang/semver"
)

// ConstantName the name we used to get constant values
type ConstantName int

const (
	BlocksPerDay ConstantName = iota
	BlocksPerYear
	OutboundTransactionFee
	NativeTransactionFee
	KillSwitchStart
	KillSwitchDuration
	PoolCycle
	MinCacaoPoolDepth
	MaxAvailablePools
	StagedPoolCost
	MinimumNodesForYggdrasil
	MinimumNodesForBFT
	DesiredValidatorSet
	AsgardSize
	ChurnInterval
	ChurnRetryInterval
	ValidatorsChangeWindow
	LeaveProcessPerBlockHeight
	BadValidatorRedline
	BadValidatorRate
	OldValidatorRate
	LowBondValidatorRate
	LackOfObservationPenalty
	SigningTransactionPeriod
	DoubleSignMaxAge
	PauseBond
	PauseUnbond
	MinimumBondInCacao
	FundMigrationInterval
	ArtificialRagnarokBlockHeight
	MaximumLiquidityCacao
	StrictBondLiquidityRatio
	DefaultPoolStatus
	MaxOutboundAttempts
	SlashPenalty
	PauseOnSlashThreshold
	FailKeygenSlashPoints
	FailKeysignSlashPoints
	LiquidityLockUpBlocks
	ObserveSlashPoints
	ObservationDelayFlexibility
	ForgiveSlashPeriod
	YggFundLimit
	YggFundRetry
	JailTimeKeygen
	JailTimeKeysign
	NodePauseChainBlocks
	MinSwapsPerBlock
	MaxSwapsPerBlock
	MaxSlashRatio
	MaxSynthPerAssetDepth
	VirtualMultSynths
	VirtualMultSynthsBasisPoints
	MinSlashPointsForBadValidator
	FullImpLossProtectionBlocks
	BondLockupPeriod
	MaxBondProviders
	NumberOfNewNodesPerChurn
	MinTxOutVolumeThreshold
	TxOutDelayRate
	TxOutDelayMax
	MaxTxOutOffset
	TNSRegisterFee
	TNSFeeOnSale
	TNSFeePerBlock
	PermittedSolvencyGap
	NodeOperatorFee
	ValidatorMaxRewardRatio
	PoolDepthForYggFundingMin
	MaxNodeToChurnOutForLowVersion
	MayaFundPerc
	MinCacaoForMayaFundDist
	WithdrawLimitTier1
	WithdrawLimitTier2
	WithdrawLimitTier3
	WithdrawDaysTier1
	WithdrawDaysTier2
	WithdrawDaysTier3
	WithdrawTier1
	WithdrawTier2
	WithdrawTier3
	InflationPercentageThreshold
	InflationPoolPercentage
	InflationFormulaMulValue
	InflationFormulaSumValue
	IBCReceiveEnabled
	IBCSendEnabled
	RagnarokProcessNumOfLPPerIteration
	SwapOutDexAggregationDisabled
	POLMaxNetworkDeposit
	POLMaxPoolMovement
	POLSynthUtilization
	POLBuffer
	SynthYieldBasisPoints
	SynthYieldCycle
	MinimumL1OutboundFeeUSD
	MinimumPoolLiquidityFee
	SubsidizeReserveMultiplier
	LiquidityAuction
	IncentiveCurveControl
	FullImpLossProtectionBlocksTimes4
	ZeroImpLossProtectionBlocks
	AllowWideBlame
)

var nameToString = map[ConstantName]string{
	BlocksPerDay:                       "BlocksPerDay",
	BlocksPerYear:                      "BlocksPerYear",
	OutboundTransactionFee:             "OutboundTransactionFee",
	NativeTransactionFee:               "NativeTransactionFee",
	PoolCycle:                          "PoolCycle",
	MinCacaoPoolDepth:                  "MinRunePoolDepth", // Can't change the string value, because we would have to account for the version change when mimir is used
	MaxAvailablePools:                  "MaxAvailablePools",
	StagedPoolCost:                     "StagedPoolCost",
	KillSwitchStart:                    "KillSwitchStart",
	KillSwitchDuration:                 "KillSwitchDuration",
	MinimumNodesForYggdrasil:           "MinimumNodesForYggdrasil",
	MinimumNodesForBFT:                 "MinimumNodesForBFT",
	DesiredValidatorSet:                "DesiredValidatorSet",
	AsgardSize:                         "AsgardSize",
	ChurnInterval:                      "ChurnInterval",
	ChurnRetryInterval:                 "ChurnRetryInterval",
	ValidatorsChangeWindow:             "ValidatorsChangeWindow",
	LeaveProcessPerBlockHeight:         "LeaveProcessPerBlockHeight",
	BadValidatorRedline:                "BadValidatorRedline",
	BadValidatorRate:                   "BadValidatorRate",
	OldValidatorRate:                   "OldValidatorRate",
	LowBondValidatorRate:               "LowBondValidatorRate",
	LackOfObservationPenalty:           "LackOfObservationPenalty",
	SigningTransactionPeriod:           "SigningTransactionPeriod",
	DoubleSignMaxAge:                   "DoubleSignMaxAge",
	PauseBond:                          "PauseBond",
	PauseUnbond:                        "PauseUnbond",
	MinimumBondInCacao:                 "MinimumBondInRune", // Can't change the string value, because we would have to account for the version change when mimir is used
	MaxBondProviders:                   "MaxBondProviders",
	FundMigrationInterval:              "FundMigrationInterval",
	ArtificialRagnarokBlockHeight:      "ArtificialRagnarokBlockHeight",
	MaximumLiquidityCacao:              "MaximumLiquidityRune", // Can't change the string value, because we would have to account for the version change when mimir is used
	StrictBondLiquidityRatio:           "StrictBondLiquidityRatio",
	DefaultPoolStatus:                  "DefaultPoolStatus",
	MaxOutboundAttempts:                "MaxOutboundAttempts",
	SlashPenalty:                       "SlashPenalty",
	PauseOnSlashThreshold:              "PauseOnSlashThreshold",
	FailKeygenSlashPoints:              "FailKeygenSlashPoints",
	FailKeysignSlashPoints:             "FailKeysignSlashPoints",
	LiquidityLockUpBlocks:              "LiquidityLockUpBlocks",
	ObserveSlashPoints:                 "ObserveSlashPoints",
	ObservationDelayFlexibility:        "ObservationDelayFlexibility",
	ForgiveSlashPeriod:                 "ForgiveSlashPeriod",
	YggFundLimit:                       "YggFundLimit",
	YggFundRetry:                       "YggFundRetry",
	JailTimeKeygen:                     "JailTimeKeygen",
	JailTimeKeysign:                    "JailTimeKeysign",
	NodePauseChainBlocks:               "NodePauseChainBlocks",
	MinSwapsPerBlock:                   "MinSwapsPerBlock",
	MaxSwapsPerBlock:                   "MaxSwapsPerBlock",
	VirtualMultSynths:                  "VirtualMultSynths",
	VirtualMultSynthsBasisPoints:       "VirtualMultSynthsBasisPoints",
	MaxSynthPerAssetDepth:              "MaxSynthPerAssetDepth",
	MinSlashPointsForBadValidator:      "MinSlashPointsForBadValidator",
	MaxSlashRatio:                      "MaxSlashRatio",
	FullImpLossProtectionBlocks:        "FullImpLossProtectionBlocks",
	BondLockupPeriod:                   "BondLockupPeriod",
	NumberOfNewNodesPerChurn:           "NumberOfNewNodesPerChurn",
	MinTxOutVolumeThreshold:            "MinTxOutVolumeThreshold",
	TxOutDelayRate:                     "TxOutDelayRate",
	TxOutDelayMax:                      "TxOutDelayMax",
	MaxTxOutOffset:                     "MaxTxOutOffset",
	TNSRegisterFee:                     "TNSRegisterFee",
	TNSFeeOnSale:                       "TNSFeeOnSale",
	TNSFeePerBlock:                     "TNSFeePerBlock",
	PermittedSolvencyGap:               "PermittedSolvencyGap",
	ValidatorMaxRewardRatio:            "ValidatorMaxRewardRatio",
	NodeOperatorFee:                    "NodeOperatorFee",
	PoolDepthForYggFundingMin:          "PoolDepthForYggFundingMin",
	MaxNodeToChurnOutForLowVersion:     "MaxNodeToChurnOutForLowVersion",
	MayaFundPerc:                       "MayaFundPerc",
	MinCacaoForMayaFundDist:            "MinRuneForMayaFundDist", // Can't change the string value, because we would have to account for the version change when mimir is used
	WithdrawLimitTier1:                 "WithdrawLimitTier1",
	WithdrawLimitTier2:                 "WithdrawLimitTier2",
	WithdrawLimitTier3:                 "WithdrawLimitTier3",
	WithdrawDaysTier1:                  "WithdrawDaysTier1",
	WithdrawDaysTier2:                  "WithdrawDaysTier2",
	WithdrawDaysTier3:                  "WithdrawDaysTier3",
	WithdrawTier1:                      "WithdrawTier1",
	WithdrawTier2:                      "WithdrawTier2",
	WithdrawTier3:                      "WithdrawTier3",
	InflationPercentageThreshold:       "InflationPercentageThreshold",
	InflationPoolPercentage:            "InflationPoolPercentage",
	InflationFormulaMulValue:           "InflationFormulaMulValue",
	InflationFormulaSumValue:           "InflationFormulaSumValue",
	IBCReceiveEnabled:                  "IBCReceiveEnabled",
	IBCSendEnabled:                     "IBCSendEnabled",
	SwapOutDexAggregationDisabled:      "SwapOutDexAggregationDisabled",
	POLMaxNetworkDeposit:               "POLMaxNetworkDeposit",
	POLMaxPoolMovement:                 "POLMaxPoolMovement",
	POLSynthUtilization:                "POLSynthUtilization",
	POLBuffer:                          "POLBuffer",
	RagnarokProcessNumOfLPPerIteration: "RagnarokProcessNumOfLPPerIteration",
	SynthYieldBasisPoints:              "SynthYieldBasisPoints",
	SynthYieldCycle:                    "SynthYieldCycle",
	MinimumL1OutboundFeeUSD:            "MinimumL1OutboundFeeUSD",
	MinimumPoolLiquidityFee:            "MinimumPoolLiquidityFee",
	SubsidizeReserveMultiplier:         "SubsidizeReserveMultiplier",
	LiquidityAuction:                   "LiquidityAuction",
	IncentiveCurveControl:              "IncentiveCurveControl",
	FullImpLossProtectionBlocksTimes4:  "FullImpLossProtectionBlocksTimes4",
	ZeroImpLossProtectionBlocks:        "ZeroImpLossProtectionBlocks",
	AllowWideBlame:                     "AllowWideBlame",
}

// String implement fmt.stringer
func (cn ConstantName) String() string {
	val, ok := nameToString[cn]
	if !ok {
		return "NA"
	}
	return val
}

// ConstantValues define methods used to get constant values
type ConstantValues interface {
	fmt.Stringer
	GetInt64Value(name ConstantName) int64
	GetBoolValue(name ConstantName) bool
	GetStringValue(name ConstantName) string
}

// GetConstantValues will return an  implementation of ConstantValues which provide ways to get constant values
func GetConstantValues(ver semver.Version) ConstantValues {
	switch {
	case ver.GTE(semver.MustParse("1.102.0")):
		return NewConstantValue102()
	case ver.GTE(semver.MustParse("0.1.0")):
		return NewConstantValue010()
	default:
		return nil
	}
}
