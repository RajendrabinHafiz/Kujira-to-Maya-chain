# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: mayachain/v1/x/mayachain/types/misc.proto, mayachain/v1/x/mayachain/types/msg_add_liquidity.proto, mayachain/v1/x/mayachain/types/msg_ban.proto, mayachain/v1/x/mayachain/types/msg_bond.proto, mayachain/v1/x/mayachain/types/msg_consolidate.proto, mayachain/v1/x/mayachain/types/msg_deposit.proto, mayachain/v1/x/mayachain/types/msg_donate.proto, mayachain/v1/x/mayachain/types/msg_errata.proto, mayachain/v1/x/mayachain/types/msg_forgive_slash.proto, mayachain/v1/x/mayachain/types/msg_leave.proto, mayachain/v1/x/mayachain/types/msg_manage_mayaname.proto, mayachain/v1/x/mayachain/types/msg_migrate.proto, mayachain/v1/x/mayachain/types/msg_mimir.proto, mayachain/v1/x/mayachain/types/msg_network_fee.proto, mayachain/v1/x/mayachain/types/msg_node_pause_chain.proto, mayachain/v1/x/mayachain/types/msg_noop.proto, mayachain/v1/x/mayachain/types/msg_observed_txin.proto, mayachain/v1/x/mayachain/types/msg_observed_txout.proto, mayachain/v1/x/mayachain/types/msg_ragnarok.proto, mayachain/v1/x/mayachain/types/msg_refund.proto, mayachain/v1/x/mayachain/types/msg_reserve_contributor.proto, mayachain/v1/x/mayachain/types/msg_send.proto, mayachain/v1/x/mayachain/types/msg_set_aztec_address.proto, mayachain/v1/x/mayachain/types/msg_set_ip_address.proto, mayachain/v1/x/mayachain/types/msg_set_node_keys.proto, mayachain/v1/x/mayachain/types/msg_solvency.proto, mayachain/v1/x/mayachain/types/msg_swap.proto, mayachain/v1/x/mayachain/types/msg_tss_keysign_fail.proto, mayachain/v1/x/mayachain/types/msg_tss_pool.proto, mayachain/v1/x/mayachain/types/msg_tx_outbound.proto, mayachain/v1/x/mayachain/types/msg_unbond.proto, mayachain/v1/x/mayachain/types/msg_version.proto, mayachain/v1/x/mayachain/types/msg_withdraw_liquidity.proto, mayachain/v1/x/mayachain/types/msg_yggdrasil.proto, mayachain/v1/x/mayachain/types/type_ban_voter.proto, mayachain/v1/x/mayachain/types/type_blame.proto, mayachain/v1/x/mayachain/types/type_chain_contract.proto, mayachain/v1/x/mayachain/types/type_errata_tx_voter.proto, mayachain/v1/x/mayachain/types/type_events.proto, mayachain/v1/x/mayachain/types/type_forgive_slash_voter.proto, mayachain/v1/x/mayachain/types/type_jail.proto, mayachain/v1/x/mayachain/types/type_keygen.proto, mayachain/v1/x/mayachain/types/type_liquidity_auction_tier.proto, mayachain/v1/x/mayachain/types/type_liquidity_provider.proto, mayachain/v1/x/mayachain/types/type_mayaname.proto, mayachain/v1/x/mayachain/types/type_mimir.proto, mayachain/v1/x/mayachain/types/type_network.proto, mayachain/v1/x/mayachain/types/type_network_fee.proto, mayachain/v1/x/mayachain/types/type_node_account.proto, mayachain/v1/x/mayachain/types/type_node_pause_chain.proto, mayachain/v1/x/mayachain/types/type_observed_network_fee.proto, mayachain/v1/x/mayachain/types/type_observed_tx.proto, mayachain/v1/x/mayachain/types/type_pol.proto, mayachain/v1/x/mayachain/types/type_pool.proto, mayachain/v1/x/mayachain/types/type_ragnarok.proto, mayachain/v1/x/mayachain/types/type_reserve_contributor.proto, mayachain/v1/x/mayachain/types/type_solvency_voter.proto, mayachain/v1/x/mayachain/types/type_tss.proto, mayachain/v1/x/mayachain/types/type_tss_keysign.proto, mayachain/v1/x/mayachain/types/type_tss_metric.proto, mayachain/v1/x/mayachain/types/type_tx_out.proto, mayachain/v1/x/mayachain/types/type_vault.proto
# plugin: python-betterproto
from dataclasses import dataclass
from typing import List

import betterproto
from betterproto.grpc.grpclib_server import ServiceBase


class Status(betterproto.Enum):
    incomplete = 0
    done = 1
    reverted = 2


class KeygenType(betterproto.Enum):
    UnknownKeygen = 0
    AsgardKeygen = 1
    YggdrasilKeygen = 2


class NodeStatus(betterproto.Enum):
    Unknown = 0
    Whitelisted = 1
    Standby = 2
    Ready = 3
    Active = 4
    Disabled = 5


class NodeType(betterproto.Enum):
    TypeValidator = 0
    TypeVault = 1
    TypeUnknown = 2


class PoolStatus(betterproto.Enum):
    """
    |    State    | Swap | Add   | Withdraw  | Refunding | | ----------- | ----
    | ----- | --------- | --------- | | `staged`    | no   | yes   | yes
    | Refund Invalid Add/Remove Liquidity && all Swaps | | `available` | yes  |
    yes   | yes       | Refund Invalid Tx | | `suspended` | no   | no    | no
    | Refund all |
    """

    UnknownPoolStatus = 0
    Available = 1
    Staged = 2
    Suspended = 4


class PendingLiquidityType(betterproto.Enum):
    add = 0
    withdraw = 1


class BondType(betterproto.Enum):
    bond_paid = 0
    bond_returned = 1
    bond_reward = 2
    bond_cost = 3


class VaultType(betterproto.Enum):
    UnknownVault = 0
    AsgardVault = 1
    YggdrasilVault = 2


class VaultStatus(betterproto.Enum):
    InactiveVault = 0
    ActiveVault = 1
    RetiringVault = 2
    InitVault = 3


class OrderType(betterproto.Enum):
    market = 0
    limit = 1


@dataclass(eq=False, repr=False)
class LiquidityAuctionTier(betterproto.Message):
    address: str = betterproto.string_field(1)
    tier: int = betterproto.int64_field(2)


@dataclass(eq=False, repr=False)
class ReserveContributor(betterproto.Message):
    address: str = betterproto.string_field(1)
    amount: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class MsgDeposit(betterproto.Message):
    coins: List["_common__.Coin"] = betterproto.message_field(1)
    memo: str = betterproto.string_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class TxOutItem(betterproto.Message):
    chain: str = betterproto.string_field(1)
    to_address: str = betterproto.string_field(2)
    vault_pub_key: str = betterproto.string_field(3)
    coin: "_common__.Coin" = betterproto.message_field(4)
    memo: str = betterproto.string_field(5)
    max_gas: List["_common__.Coin"] = betterproto.message_field(6)
    gas_rate: int = betterproto.int64_field(7)
    in_hash: str = betterproto.string_field(8)
    out_hash: str = betterproto.string_field(9)
    module_name: str = betterproto.string_field(10)
    aggregator: str = betterproto.string_field(11)
    aggregator_target_asset: str = betterproto.string_field(12)
    aggregator_target_limit: str = betterproto.string_field(13)


@dataclass(eq=False, repr=False)
class TxOut(betterproto.Message):
    height: int = betterproto.int64_field(1)
    tx_array: List["TxOutItem"] = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class ObservedTx(betterproto.Message):
    tx: "_common__.Tx" = betterproto.message_field(1)
    status: "Status" = betterproto.enum_field(2)
    out_hashes: List[str] = betterproto.string_field(3)
    block_height: int = betterproto.int64_field(4)
    signers: List[str] = betterproto.string_field(5)
    observed_pub_key: str = betterproto.string_field(6)
    keysign_ms: int = betterproto.int64_field(7)
    finalise_height: int = betterproto.int64_field(8)
    aggregator: str = betterproto.string_field(9)
    aggregator_target: str = betterproto.string_field(10)
    aggregator_target_limit: str = betterproto.string_field(11)


@dataclass(eq=False, repr=False)
class ObservedTxVoter(betterproto.Message):
    tx_id: str = betterproto.string_field(1)
    tx: "ObservedTx" = betterproto.message_field(2)
    height: int = betterproto.int64_field(3)
    txs: List["ObservedTx"] = betterproto.message_field(4)
    actions: List["TxOutItem"] = betterproto.message_field(5)
    out_txs: List["_common__.Tx"] = betterproto.message_field(6)
    finalised_height: int = betterproto.int64_field(7)
    updated_vault: bool = betterproto.bool_field(8)
    reverted: bool = betterproto.bool_field(9)
    outbound_height: int = betterproto.int64_field(10)


@dataclass(eq=False, repr=False)
class MsgMigrate(betterproto.Message):
    tx: "ObservedTx" = betterproto.message_field(1)
    block_height: int = betterproto.int64_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class NodeMimir(betterproto.Message):
    key: str = betterproto.string_field(1)
    value: int = betterproto.int64_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class NodeMimirs(betterproto.Message):
    mimirs: List["NodeMimir"] = betterproto.message_field(1)


@dataclass(eq=False, repr=False)
class MsgSetNodeKeys(betterproto.Message):
    pub_key_set_set: "_common__.PubKeySet" = betterproto.message_field(1)
    validator_cons_pub_key: str = betterproto.string_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class MsgObservedTxIn(betterproto.Message):
    txs: List["ObservedTx"] = betterproto.message_field(1)
    signer: bytes = betterproto.bytes_field(2)


@dataclass(eq=False, repr=False)
class MsgYggdrasil(betterproto.Message):
    tx: "_common__.Tx" = betterproto.message_field(1)
    pub_key: str = betterproto.string_field(2)
    add_funds: bool = betterproto.bool_field(3)
    coins: List["_common__.Coin"] = betterproto.message_field(4)
    block_height: int = betterproto.int64_field(5)
    signer: bytes = betterproto.bytes_field(6)


@dataclass(eq=False, repr=False)
class MsgObservedTxOut(betterproto.Message):
    txs: List["ObservedTx"] = betterproto.message_field(1)
    signer: bytes = betterproto.bytes_field(2)


@dataclass(eq=False, repr=False)
class MsgDonate(betterproto.Message):
    asset: "_common__.Asset" = betterproto.message_field(1)
    asset_amount: str = betterproto.string_field(2)
    cacao_amount: str = betterproto.string_field(3)
    tx: "_common__.Tx" = betterproto.message_field(4)
    signer: bytes = betterproto.bytes_field(5)


@dataclass(eq=False, repr=False)
class MsgSetIpAddress(betterproto.Message):
    ip_address: str = betterproto.string_field(1)
    signer: bytes = betterproto.bytes_field(2)


@dataclass(eq=False, repr=False)
class MsgManageMayaName(betterproto.Message):
    name: str = betterproto.string_field(1)
    chain: str = betterproto.string_field(2)
    address: str = betterproto.string_field(3)
    coin: "_common__.Coin" = betterproto.message_field(4)
    expire_block_height: int = betterproto.int64_field(5)
    preferred_asset: "_common__.Asset" = betterproto.message_field(6)
    owner: bytes = betterproto.bytes_field(7)
    signer: bytes = betterproto.bytes_field(8)


@dataclass(eq=False, repr=False)
class MsgSend(betterproto.Message):
    from_address: bytes = betterproto.bytes_field(1)
    to_address: bytes = betterproto.bytes_field(2)
    amount: List["_cosmos_base_v1_beta1__.Coin"] = betterproto.message_field(3)


@dataclass(eq=False, repr=False)
class MsgNoOp(betterproto.Message):
    observed_tx: "ObservedTx" = betterproto.message_field(1)
    signer: bytes = betterproto.bytes_field(2)
    action: str = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class MsgConsolidate(betterproto.Message):
    observed_tx: "ObservedTx" = betterproto.message_field(1)
    signer: bytes = betterproto.bytes_field(2)


@dataclass(eq=False, repr=False)
class ObservedNetworkFeeVoter(betterproto.Message):
    block_height: int = betterproto.int64_field(1)
    report_block_height: int = betterproto.int64_field(2)
    chain: str = betterproto.string_field(3)
    signers: List[str] = betterproto.string_field(4)
    fee_rate: int = betterproto.int64_field(5)


@dataclass(eq=False, repr=False)
class ChainContract(betterproto.Message):
    chain: str = betterproto.string_field(1)
    router: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class MsgSetAztecAddress(betterproto.Message):
    aztec_address: str = betterproto.string_field(1)
    signer: bytes = betterproto.bytes_field(2)


@dataclass(eq=False, repr=False)
class Node(betterproto.Message):
    pubkey: str = betterproto.string_field(1)
    blame_data: bytes = betterproto.bytes_field(2)
    blame_signature: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class Blame(betterproto.Message):
    fail_reason: str = betterproto.string_field(1)
    is_unicast: bool = betterproto.bool_field(2)
    blame_nodes: List["Node"] = betterproto.message_field(3)


@dataclass(eq=False, repr=False)
class Keygen(betterproto.Message):
    id: str = betterproto.string_field(1)
    type: "KeygenType" = betterproto.enum_field(2)
    members: List[str] = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class KeygenBlock(betterproto.Message):
    height: int = betterproto.int64_field(1)
    keygens: List["Keygen"] = betterproto.message_field(4)


@dataclass(eq=False, repr=False)
class MsgTssPool(betterproto.Message):
    id: str = betterproto.string_field(1)
    pool_pub_key: str = betterproto.string_field(2)
    keygen_type: "KeygenType" = betterproto.enum_field(3)
    pub_keys: List[str] = betterproto.string_field(4)
    height: int = betterproto.int64_field(5)
    blame: "Blame" = betterproto.message_field(6)
    chains: List[str] = betterproto.string_field(7)
    signer: bytes = betterproto.bytes_field(8)
    keygen_time: int = betterproto.int64_field(9)


@dataclass(eq=False, repr=False)
class MsgRefundTx(betterproto.Message):
    tx: "ObservedTx" = betterproto.message_field(1)
    in_tx_id: str = betterproto.string_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class MsgReserveContributor(betterproto.Message):
    tx: "_common__.Tx" = betterproto.message_field(1)
    contributor: "ReserveContributor" = betterproto.message_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class Jail(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(1)
    release_height: int = betterproto.int64_field(2)
    reason: str = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class MsgNodePauseChain(betterproto.Message):
    value: int = betterproto.int64_field(1)
    signer: bytes = betterproto.bytes_field(2)


@dataclass(eq=False, repr=False)
class NodeAccount(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(1)
    status: "NodeStatus" = betterproto.enum_field(2)
    pub_key_set: "_common__.PubKeySet" = betterproto.message_field(3)
    aztec_address: str = betterproto.string_field(17)
    validator_cons_pub_key: str = betterproto.string_field(4)
    bond: str = betterproto.string_field(5)
    active_block_height: int = betterproto.int64_field(6)
    bond_address: str = betterproto.string_field(7)
    status_since: int = betterproto.int64_field(8)
    signer_membership: List[str] = betterproto.string_field(9)
    requested_to_leave: bool = betterproto.bool_field(10)
    forced_to_leave: bool = betterproto.bool_field(11)
    leave_score: int = betterproto.uint64_field(12)
    ip_address: str = betterproto.string_field(13)
    version: str = betterproto.string_field(14)
    type: "NodeType" = betterproto.enum_field(15)
    reward: str = betterproto.string_field(16)


@dataclass(eq=False, repr=False)
class BondProvider(betterproto.Message):
    bond_address: bytes = betterproto.bytes_field(1)
    bonded: bool = betterproto.bool_field(3)


@dataclass(eq=False, repr=False)
class BondProviders(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(1)
    node_operator_fee: str = betterproto.string_field(2)
    providers: List["BondProvider"] = betterproto.message_field(3)


@dataclass(eq=False, repr=False)
class TssKeysignFailVoter(betterproto.Message):
    id: str = betterproto.string_field(1)
    height: int = betterproto.int64_field(4)
    signers: List[str] = betterproto.string_field(6)


@dataclass(eq=False, repr=False)
class Network(betterproto.Message):
    bond_reward_rune: str = betterproto.string_field(1)
    total_bond_units: str = betterproto.string_field(2)
    lp_income_split: int = betterproto.int64_field(3)
    node_income_split: int = betterproto.int64_field(4)


@dataclass(eq=False, repr=False)
class Pool(betterproto.Message):
    balance_cacao: str = betterproto.string_field(1)
    balance_asset: str = betterproto.string_field(2)
    asset: "_common__.Asset" = betterproto.message_field(3)
    lp_units: str = betterproto.string_field(4)
    status: "PoolStatus" = betterproto.enum_field(5)
    status_since: int = betterproto.int64_field(10)
    decimals: int = betterproto.int64_field(6)
    synth_units: str = betterproto.string_field(7)
    pending_inbound_cacao: str = betterproto.string_field(8)
    pending_inbound_asset: str = betterproto.string_field(9)


@dataclass(eq=False, repr=False)
class PoolMod(betterproto.Message):
    asset: "_common__.Asset" = betterproto.message_field(1)
    cacao_amt: str = betterproto.string_field(2)
    cacao_add: bool = betterproto.bool_field(3)
    asset_amt: str = betterproto.string_field(4)
    asset_add: bool = betterproto.bool_field(5)


@dataclass(eq=False, repr=False)
class EventSwap(betterproto.Message):
    pool: "_common__.Asset" = betterproto.message_field(1)
    swap_target: str = betterproto.string_field(2)
    swap_slip: str = betterproto.string_field(3)
    liquidity_fee: str = betterproto.string_field(4)
    liquidity_fee_in_cacao: str = betterproto.string_field(5)
    in_tx: "_common__.Tx" = betterproto.message_field(6)
    out_txs: "_common__.Tx" = betterproto.message_field(7)
    emit_asset: "_common__.Coin" = betterproto.message_field(8)
    synth_units: str = betterproto.string_field(9)


@dataclass(eq=False, repr=False)
class EventAddLiquidity(betterproto.Message):
    pool: "_common__.Asset" = betterproto.message_field(1)
    provider_units: str = betterproto.string_field(2)
    cacao_address: str = betterproto.string_field(3)
    cacao_amount: str = betterproto.string_field(4)
    asset_amount: str = betterproto.string_field(5)
    cacao_tx_id: str = betterproto.string_field(6)
    asset_tx_id: str = betterproto.string_field(7)
    asset_address: str = betterproto.string_field(8)


@dataclass(eq=False, repr=False)
class EventWithdraw(betterproto.Message):
    pool: "_common__.Asset" = betterproto.message_field(1)
    provider_units: str = betterproto.string_field(2)
    basis_points: int = betterproto.int64_field(3)
    asymmetry: bytes = betterproto.bytes_field(4)
    in_tx: "_common__.Tx" = betterproto.message_field(5)
    emit_asset: str = betterproto.string_field(6)
    emit_cacao: str = betterproto.string_field(7)
    imp_loss_protection: str = betterproto.string_field(8)


@dataclass(eq=False, repr=False)
class EventPendingLiquidity(betterproto.Message):
    pool: "_common__.Asset" = betterproto.message_field(1)
    pending_type: "PendingLiquidityType" = betterproto.enum_field(2)
    cacao_address: str = betterproto.string_field(3)
    cacao_amount: str = betterproto.string_field(4)
    asset_address: str = betterproto.string_field(5)
    asset_amount: str = betterproto.string_field(6)
    cacao_tx_id: str = betterproto.string_field(7)
    asset_tx_id: str = betterproto.string_field(8)


@dataclass(eq=False, repr=False)
class EventDonate(betterproto.Message):
    pool: "_common__.Asset" = betterproto.message_field(1)
    in_tx: "_common__.Tx" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class EventPool(betterproto.Message):
    pool: "_common__.Asset" = betterproto.message_field(1)
    status: "PoolStatus" = betterproto.enum_field(2)


@dataclass(eq=False, repr=False)
class PoolAmt(betterproto.Message):
    asset: "_common__.Asset" = betterproto.message_field(1)
    amount: int = betterproto.int64_field(2)


@dataclass(eq=False, repr=False)
class EventRewards(betterproto.Message):
    bond_reward: str = betterproto.string_field(1)
    pool_rewards: List["PoolAmt"] = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class EventRefund(betterproto.Message):
    code: int = betterproto.uint32_field(1)
    reason: str = betterproto.string_field(2)
    in_tx: "_common__.Tx" = betterproto.message_field(3)
    fee: "_common__.Fee" = betterproto.message_field(4)


@dataclass(eq=False, repr=False)
class EventBond(betterproto.Message):
    amount: str = betterproto.string_field(1)
    bond_type: "BondType" = betterproto.enum_field(2)
    tx_in: "_common__.Tx" = betterproto.message_field(3)


@dataclass(eq=False, repr=False)
class GasPool(betterproto.Message):
    asset: "_common__.Asset" = betterproto.message_field(1)
    cacao_amt: str = betterproto.string_field(2)
    asset_amt: str = betterproto.string_field(3)
    count: int = betterproto.int64_field(4)


@dataclass(eq=False, repr=False)
class EventGas(betterproto.Message):
    pools: List["GasPool"] = betterproto.message_field(1)


@dataclass(eq=False, repr=False)
class EventReserve(betterproto.Message):
    reserve_contributor: "ReserveContributor" = betterproto.message_field(1)
    in_tx: "_common__.Tx" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class EventScheduledOutbound(betterproto.Message):
    out_tx: "TxOutItem" = betterproto.message_field(1)


@dataclass(eq=False, repr=False)
class EventSecurity(betterproto.Message):
    msg: str = betterproto.string_field(1)
    tx: "_common__.Tx" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class EventSlash(betterproto.Message):
    pool: "_common__.Asset" = betterproto.message_field(1)
    slash_amount: List["PoolAmt"] = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class EventSlashLiquidity(betterproto.Message):
    node_bond_address: bytes = betterproto.bytes_field(1)
    asset: "_common__.Asset" = betterproto.message_field(2)
    address: str = betterproto.string_field(3)
    lp_units: str = betterproto.string_field(4)


@dataclass(eq=False, repr=False)
class EventErrata(betterproto.Message):
    tx_id: str = betterproto.string_field(1)
    pools: List["PoolMod"] = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class EventFee(betterproto.Message):
    tx_id: str = betterproto.string_field(1)
    fee: "_common__.Fee" = betterproto.message_field(2)
    synth_units: str = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class EventOutbound(betterproto.Message):
    in_tx_id: str = betterproto.string_field(1)
    tx: "_common__.Tx" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class EventTssKeygenMetric(betterproto.Message):
    pub_key: str = betterproto.string_field(1)
    median_duration_ms: int = betterproto.int64_field(2)


@dataclass(eq=False, repr=False)
class EventTssKeysignMetric(betterproto.Message):
    tx_id: str = betterproto.string_field(1)
    median_duration_ms: int = betterproto.int64_field(2)


@dataclass(eq=False, repr=False)
class EventSlashPoint(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(1)
    slash_points: int = betterproto.int64_field(2)
    reason: str = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class EventPoolBalanceChanged(betterproto.Message):
    pool_change: "PoolMod" = betterproto.message_field(1)
    reason: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class EventSwitch(betterproto.Message):
    to_address: bytes = betterproto.bytes_field(1)
    from_address: str = betterproto.string_field(2)
    burn: "_common__.Coin" = betterproto.message_field(3)
    tx_id: str = betterproto.string_field(4)


@dataclass(eq=False, repr=False)
class EventSwitchV87(betterproto.Message):
    to_address: bytes = betterproto.bytes_field(1)
    from_address: str = betterproto.string_field(2)
    burn: "_common__.Coin" = betterproto.message_field(3)
    tx_id: str = betterproto.string_field(4)
    mint: str = betterproto.string_field(5)


@dataclass(eq=False, repr=False)
class EventMayaName(betterproto.Message):
    name: str = betterproto.string_field(1)
    chain: str = betterproto.string_field(2)
    address: str = betterproto.string_field(3)
    registration_fee: str = betterproto.string_field(4)
    fund_amt: str = betterproto.string_field(5)
    expire: int = betterproto.int64_field(6)
    owner: bytes = betterproto.bytes_field(7)


@dataclass(eq=False, repr=False)
class EventSetMimir(betterproto.Message):
    key: str = betterproto.string_field(1)
    value: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class EventSetNodeMimir(betterproto.Message):
    key: str = betterproto.string_field(1)
    value: str = betterproto.string_field(2)
    address: str = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class SolvencyVoter(betterproto.Message):
    id: str = betterproto.string_field(1)
    chain: str = betterproto.string_field(2)
    pub_key: str = betterproto.string_field(3)
    coins: List["_common__.Coin"] = betterproto.message_field(4)
    height: int = betterproto.int64_field(5)
    consensus_block_height: int = betterproto.int64_field(6)
    signers: List[str] = betterproto.string_field(7)


@dataclass(eq=False, repr=False)
class MsgLeave(betterproto.Message):
    tx: "_common__.Tx" = betterproto.message_field(1)
    node_address: bytes = betterproto.bytes_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class MsgBond(betterproto.Message):
    tx_in: "_common__.Tx" = betterproto.message_field(1)
    node_address: bytes = betterproto.bytes_field(2)
    amount: str = betterproto.string_field(3)
    bond_address: str = betterproto.string_field(4)
    signer: bytes = betterproto.bytes_field(5)
    bond_provider_address: bytes = betterproto.bytes_field(6)
    operator_fee: int = betterproto.int64_field(7)


@dataclass(eq=False, repr=False)
class NetworkFee(betterproto.Message):
    """
    NetworkFee represent the fee rate and typical transaction size outbound
    from MAYANode This is to keep the information reported by bifrost For BTC
    chain, TransactionFeeRate should be sats/vbyte For Binance chain , given
    fee is fixed , thus for single coin , transaction size will be 1, and the
    rate should be 37500, for multiple coin , Transaction size should the
    number of coins
    """

    chain: str = betterproto.string_field(1)
    transaction_size: int = betterproto.uint64_field(2)
    transaction_fee_rate: int = betterproto.uint64_field(3)


@dataclass(eq=False, repr=False)
class Vault(betterproto.Message):
    block_height: int = betterproto.int64_field(1)
    pub_key: str = betterproto.string_field(2)
    coins: List["_common__.Coin"] = betterproto.message_field(3)
    type: "VaultType" = betterproto.enum_field(4)
    status: "VaultStatus" = betterproto.enum_field(5)
    status_since: int = betterproto.int64_field(6)
    membership: List[str] = betterproto.string_field(7)
    chains: List[str] = betterproto.string_field(8)
    inbound_tx_count: int = betterproto.int64_field(9)
    outbound_tx_count: int = betterproto.int64_field(10)
    pending_tx_block_heights: List[int] = betterproto.int64_field(11)
    routers: List["ChainContract"] = betterproto.message_field(22)


@dataclass(eq=False, repr=False)
class BanVoter(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(1)
    block_height: int = betterproto.int64_field(2)
    signers: List[str] = betterproto.string_field(3)


@dataclass(eq=False, repr=False)
class MsgAddLiquidity(betterproto.Message):
    tx: "_common__.Tx" = betterproto.message_field(1)
    asset: "_common__.Asset" = betterproto.message_field(2)
    asset_amount: str = betterproto.string_field(3)
    cacao_amount: str = betterproto.string_field(4)
    cacao_address: str = betterproto.string_field(5)
    asset_address: str = betterproto.string_field(6)
    affiliate_address: str = betterproto.string_field(7)
    affiliate_basis_points: str = betterproto.string_field(8)
    signer: bytes = betterproto.bytes_field(9)
    liquidity_auction_tier: int = betterproto.int64_field(10)


@dataclass(eq=False, repr=False)
class MsgSwap(betterproto.Message):
    tx: "_common__.Tx" = betterproto.message_field(1)
    target_asset: "_common__.Asset" = betterproto.message_field(2)
    destination: str = betterproto.string_field(3)
    trade_target: str = betterproto.string_field(4)
    affiliate_address: str = betterproto.string_field(5)
    affiliate_basis_points: str = betterproto.string_field(6)
    signer: bytes = betterproto.bytes_field(7)
    aggregator: str = betterproto.string_field(8)
    aggregator_target_address: str = betterproto.string_field(9)
    aggregator_target_limit: str = betterproto.string_field(10)
    order_type: "OrderType" = betterproto.enum_field(11)


@dataclass(eq=False, repr=False)
class MsgUnBond(betterproto.Message):
    tx_in: "_common__.Tx" = betterproto.message_field(1)
    node_address: bytes = betterproto.bytes_field(2)
    bond_address: str = betterproto.string_field(5)
    signer: bytes = betterproto.bytes_field(7)
    bond_provider_address: bytes = betterproto.bytes_field(8)


@dataclass(eq=False, repr=False)
class MsgMimir(betterproto.Message):
    key: str = betterproto.string_field(1)
    value: int = betterproto.int64_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class ProtocolOwnedLiquidity(betterproto.Message):
    cacao_deposited: str = betterproto.string_field(1)
    cacao_withdrawn: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class MsgSolvency(betterproto.Message):
    id: str = betterproto.string_field(1)
    chain: str = betterproto.string_field(2)
    pub_key: str = betterproto.string_field(3)
    coins: List["_common__.Coin"] = betterproto.message_field(4)
    height: int = betterproto.int64_field(5)
    signer: bytes = betterproto.bytes_field(6)


@dataclass(eq=False, repr=False)
class ErrataTxVoter(betterproto.Message):
    tx_id: str = betterproto.string_field(1)
    chain: str = betterproto.string_field(2)
    block_height: int = betterproto.int64_field(3)
    signers: List[str] = betterproto.string_field(4)


@dataclass(eq=False, repr=False)
class LiquidityProvider(betterproto.Message):
    asset: "_common__.Asset" = betterproto.message_field(1)
    cacao_address: str = betterproto.string_field(2)
    asset_address: str = betterproto.string_field(3)
    last_add_height: int = betterproto.int64_field(4)
    last_withdraw_height: int = betterproto.int64_field(5)
    units: str = betterproto.string_field(6)
    pending_cacao: str = betterproto.string_field(7)
    pending_asset: str = betterproto.string_field(8)
    pending_tx_id: str = betterproto.string_field(9)
    cacao_deposit_value: str = betterproto.string_field(10)
    asset_deposit_value: str = betterproto.string_field(11)
    node_bond_address: bytes = betterproto.bytes_field(12)
    withdraw_counter: str = betterproto.string_field(13)
    last_withdraw_counter_height: int = betterproto.int64_field(14)


@dataclass(eq=False, repr=False)
class TssVoter(betterproto.Message):
    id: str = betterproto.string_field(1)
    pool_pub_key: str = betterproto.string_field(2)
    pub_keys: List[str] = betterproto.string_field(3)
    block_height: int = betterproto.int64_field(4)
    chains: List[str] = betterproto.string_field(5)
    signers: List[str] = betterproto.string_field(6)
    majority_consensus_block_height: int = betterproto.int64_field(7)


@dataclass(eq=False, repr=False)
class MsgForgiveSlash(betterproto.Message):
    blocks: str = betterproto.string_field(2)
    node_address: bytes = betterproto.bytes_field(3)
    signer: bytes = betterproto.bytes_field(4)


@dataclass(eq=False, repr=False)
class ProtoInt64(betterproto.Message):
    value: int = betterproto.int64_field(1)


@dataclass(eq=False, repr=False)
class ProtoUint64(betterproto.Message):
    value: int = betterproto.uint64_field(1)


@dataclass(eq=False, repr=False)
class ProtoFloat64(betterproto.Message):
    value: float = betterproto.double_field(1)


@dataclass(eq=False, repr=False)
class ProtoAccAddresses(betterproto.Message):
    value: List[bytes] = betterproto.bytes_field(1)


@dataclass(eq=False, repr=False)
class ProtoStrings(betterproto.Message):
    value: List[str] = betterproto.string_field(1)


@dataclass(eq=False, repr=False)
class ProtoBools(betterproto.Message):
    value: List[bool] = betterproto.bool_field(1)


@dataclass(eq=False, repr=False)
class RagnarokWithdrawPosition(betterproto.Message):
    number: int = betterproto.int64_field(1)
    pool: "_common__.Asset" = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class MsgRagnarok(betterproto.Message):
    tx: "ObservedTx" = betterproto.message_field(1)
    block_height: int = betterproto.int64_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class MayaNameAlias(betterproto.Message):
    chain: str = betterproto.string_field(1)
    address: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class MayaName(betterproto.Message):
    name: str = betterproto.string_field(1)
    expire_block_height: int = betterproto.int64_field(2)
    owner: bytes = betterproto.bytes_field(3)
    preferred_asset: "_common__.Asset" = betterproto.message_field(4)
    aliases: List["MayaNameAlias"] = betterproto.message_field(5)


@dataclass(eq=False, repr=False)
class MsgTssKeysignFail(betterproto.Message):
    id: str = betterproto.string_field(1)
    height: int = betterproto.int64_field(2)
    blame: "Blame" = betterproto.message_field(3)
    memo: str = betterproto.string_field(4)
    coins: List["_common__.Coin"] = betterproto.message_field(5)
    pub_key: str = betterproto.string_field(6)
    signer: bytes = betterproto.bytes_field(7)


@dataclass(eq=False, repr=False)
class MsgSetVersion(betterproto.Message):
    version: str = betterproto.string_field(1)
    signer: bytes = betterproto.bytes_field(2)


@dataclass(eq=False, repr=False)
class ForgiveSlashVoter(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(1)
    block_height: int = betterproto.int64_field(2)
    proposed_block_height: int = betterproto.int64_field(3)
    signers: List[str] = betterproto.string_field(4)


@dataclass(eq=False, repr=False)
class NodeTssTime(betterproto.Message):
    address: bytes = betterproto.bytes_field(1)
    tss_time: int = betterproto.int64_field(2)


@dataclass(eq=False, repr=False)
class TssKeygenMetric(betterproto.Message):
    pub_key: str = betterproto.string_field(1)
    node_tss_times: List["NodeTssTime"] = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class TssKeysignMetric(betterproto.Message):
    tx_id: str = betterproto.string_field(1)
    node_tss_times: List["NodeTssTime"] = betterproto.message_field(2)


@dataclass(eq=False, repr=False)
class MsgWithdrawLiquidity(betterproto.Message):
    tx: "_common__.Tx" = betterproto.message_field(1)
    withdraw_address: str = betterproto.string_field(2)
    basis_points: str = betterproto.string_field(3)
    asset: "_common__.Asset" = betterproto.message_field(4)
    withdrawal_asset: "_common__.Asset" = betterproto.message_field(5)
    signer: bytes = betterproto.bytes_field(6)


@dataclass(eq=False, repr=False)
class MsgNetworkFee(betterproto.Message):
    block_height: int = betterproto.int64_field(1)
    chain: str = betterproto.string_field(2)
    transaction_size: int = betterproto.uint64_field(3)
    transaction_fee_rate: int = betterproto.uint64_field(4)
    signer: bytes = betterproto.bytes_field(5)


@dataclass(eq=False, repr=False)
class MsgErrataTx(betterproto.Message):
    tx_id: str = betterproto.string_field(1)
    chain: str = betterproto.string_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class MsgBan(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(2)
    signer: bytes = betterproto.bytes_field(3)


@dataclass(eq=False, repr=False)
class NodePauseChain(betterproto.Message):
    node_address: bytes = betterproto.bytes_field(1)
    block_height: int = betterproto.int64_field(2)


@dataclass(eq=False, repr=False)
class MsgOutboundTx(betterproto.Message):
    tx: "ObservedTx" = betterproto.message_field(1)
    in_tx_id: str = betterproto.string_field(2)
    signer: bytes = betterproto.bytes_field(3)


from .. import common as _common__
from ..cosmos.base import v1beta1 as _cosmos_base_v1_beta1__
