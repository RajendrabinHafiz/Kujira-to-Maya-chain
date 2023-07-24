# Generated by the protocol buffer compiler.  DO NOT EDIT!
# sources: mayachain/v1/common/common.proto
# plugin: python-betterproto
from dataclasses import dataclass
from typing import List

import betterproto
from betterproto.grpc.grpclib_server import ServiceBase


@dataclass(eq=False, repr=False)
class Asset(betterproto.Message):
    chain: str = betterproto.string_field(1)
    symbol: str = betterproto.string_field(2)
    ticker: str = betterproto.string_field(3)
    synth: bool = betterproto.bool_field(4)


@dataclass(eq=False, repr=False)
class Coin(betterproto.Message):
    asset: "Asset" = betterproto.message_field(1)
    amount: str = betterproto.string_field(2)
    decimals: int = betterproto.int64_field(3)


@dataclass(eq=False, repr=False)
class PubKeySet(betterproto.Message):
    """PubKeySet contains two pub keys , secp256k1 and ed25519"""

    secp256_k1: str = betterproto.string_field(1)
    ed25519: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class Tx(betterproto.Message):
    id: str = betterproto.string_field(1)
    chain: str = betterproto.string_field(2)
    from_address: str = betterproto.string_field(3)
    to_address: str = betterproto.string_field(4)
    coins: List["Coin"] = betterproto.message_field(5)
    gas: List["Coin"] = betterproto.message_field(6)
    memo: str = betterproto.string_field(7)


@dataclass(eq=False, repr=False)
class Fee(betterproto.Message):
    coins: List["Coin"] = betterproto.message_field(1)
    pool_deduct: str = betterproto.string_field(2)


@dataclass(eq=False, repr=False)
class ProtoUint(betterproto.Message):
    value: str = betterproto.string_field(1)
