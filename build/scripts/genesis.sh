#!/bin/sh

set -o pipefail

. "$(dirname "$0")/core.sh"

if [ "$NET" = "mocknet" ] || [ "$NET" = "testnet" ] || [ "$NET" = "stagenet" ]; then
  echo "Loading unsafe init for mocknet and testnet..."
  . "$(dirname "$0")/core-unsafe.sh"
  . "$(dirname "$0")/testnet/state.sh"
fi

NODES="${NODES:=1}"
SEED="${SEED:=mayanode}" # the hostname of the master node
ETH_HOST="${ETH_HOST:=http://ethereum:8545}"
THOR_BLOCK_TIME="${THOR_BLOCK_TIME:=5s}"
CHAIN_ID=${CHAIN_ID:=mayachain}

# this is required as it need to run mayanode init, otherwise tendermint related command doesn't work
if [ "$SEED" = "$(hostname)" ]; then
  if [ ! -f ~/.mayanode/config/priv_validator_key.json ]; then
    init_chain
    # remove the original generate genesis file, as below will init chain again
    rm -rf ~/.mayanode/config/genesis.json
  fi
fi

create_thor_user "$SIGNER_NAME" "$SIGNER_PASSWD" "$SIGNER_SEED_PHRASE"

VALIDATOR=$(mayanode tendermint show-validator | mayanode pubkey --bech cons)
NODE_ADDRESS=$(echo "$SIGNER_PASSWD" | mayanode keys show "$SIGNER_NAME" -a --keyring-backend file)
NODE_PUB_KEY=$(echo "$SIGNER_PASSWD" | mayanode keys show "$SIGNER_NAME" -p --keyring-backend file | mayanode pubkey)
VERSION=$(fetch_version)

if [ "$SEED" = "$(hostname)" ]; then
  echo "Setting MAYANode as genesis"
  if [ ! -f ~/.mayanode/config/genesis.json ]; then
    NODE_IP_ADDRESS=${EXTERNAL_IP:=$(curl -s http://whatismyip.akamai.com)}

    init_chain "$NODE_ADDRESS"
    add_node_account "$NODE_ADDRESS" "$VALIDATOR" "$NODE_PUB_KEY" "$VERSION" "$NODE_ADDRESS" "$NODE_PUB_KEY_ED25519" "$NODE_IP_ADDRESS"

    # disable default bank transfer, and opt to use our own custom one
    disable_bank_send
    disable_mint

    # for mocknet, add initial balances
    echo "Using NET $NET"
    if [ "$NET" = "mocknet" ]; then
      echo "Adding pool and lp bonder for first node"
      add_pool "BTC.BTC" 8
      # 100K cacao of bond
      add_account "$NODE_ADDRESS" cacao 100000000000
      add_account "$NODE_ADDRESS" maya 8440000000
      add_lp "BTC.BTC" "$NODE_ADDRESS" "" 10000000000 "$NODE_ADDRESS"

      # founders maya tokens
      add_account tmaya1m3t5wwrpfylss8e9g3jvq5chsv2xl3ucn9rdrf cacao 10000000000
      add_account tmaya1m3t5wwrpfylss8e9g3jvq5chsv2xl3ucn9rdrf maya 1560000000

      # local cluster accounts for genesis nodes
      CAT=tmaya1uuds8pd92qnnq0udw0rpg0szpgcslc9p8gps0z
      add_account "$CAT" cacao 30000000000 # cat
      add_lp "BTC.BTC" "$CAT" "" 10000000000 "$NODE_ADDRESS"
      DOG=tmaya1zf3gsk7edzwl9syyefvfhle37cjtql35hdgtzt
      add_account "$DOG" cacao 10000000000000000 # dog
      add_lp "BTC.BTC" "$DOG" "" 10000000000 "$NODE_ADDRESS"
      FOX=tmaya1qk8c8sfrmfm0tkncs0zxeutc8v5mx3pjjcq6rv
      add_account "$FOX" cacao 300000000000 # fox
      add_lp "BTC.BTC" "$FOX" "" 10000000000 "$NODE_ADDRESS"
      PIG=tmaya13wrmhnh2qe98rjse30pl7u6jxszjjwl4fd6gwn
      add_account "$PIG" cacao 300000000000 # pig
      add_lp "BTC.BTC" "$PIG" "" 10000000000 "$NODE_ADDRESS"

      add_mayaname itzamna 5256000 tmaya18z343fsdlav47chtkyp0aawqt6sgxsh3g94648 MAYA
      add_mayaname aaluxx 5256000 tmaya1thtn0cdcu2sujtklknlc9wl2dg09a07v0lgnkr tmaya1thtn0cdcu2sujtklknlc9wl2dg09a07v0lgnkr MAYA
      add_mayaname wr 5256000 tmaya14lkndecaw0zkzu0yq4a0qq869hrs8hh7cqajxy 0x259fe33B59ae7C24B34C0cb957271Ae0f7786878 ETH
      add_mayaname wr 5256000 tmaya14lkndecaw0zkzu0yq4a0qq869hrs8hh7cqajxy tmaya1a427q3v96psuj4fnughdw8glt5r7j38lk7v2wj MAYA
      add_mayaname wr 5256000 tmaya14lkndecaw0zkzu0yq4a0qq869hrs8hh7cqajxy tthor1a427q3v96psuj4fnughdw8glt5r7j38lkfjxcz THOR

      echo "Setting up accounts"
      # heimdall accounts
      add_account tmaya1z63f3mzwv3g75az80xwmhrawdqcjpaekkcgpz9 cacao 500000000000000
      add_account tmaya1wz78qmrkplrdhy37tw0tnvn0tkm5pqd6z6lxzw cacao 2500000000000100
      add_account tmaya18f55frcvknxvcpx2vvpfedvw4l8eutuhkt3sy2 cacao 2500000000000100
      add_account tmaya1xwusttz86hqfuk5z7amcgqsg7vp6g8zhsk2n26 cacao 509000000000000

      reserve 22000000000000000
      asgard 50000000000000

      # deploy eth contract
      deploy_eth_contract "$ETH_HOST"

      # override block time for faster smoke tests
      block_time "$THOR_BLOCK_TIME"
    else
      echo "ETH Contract Address: $CONTRACT"
      set_eth_contract "$CONTRACT"
    fi

    if [ "$NET" = "mainnet" ]; then
      add_pool "BTC.BTC" 8
      add_lp "BTC.BTC" "$NODE_ADDRESS" "" 1000000000000 "$NODE_ADDRESS"

      echo "Setting up accounts"
      add_account maya18z343fsdlav47chtkyp0aawqt6sgxsh3vjy2vz cacao 9000000000000000
      add_account maya18z343fsdlav47chtkyp0aawqt6sgxsh3vjy2vz maya 8440000000
      add_account maya1m3t5wwrpfylss8e9g3jvq5chsv2xl3uchjja6v maya 1560000000
      # genesis nodes
      add_account maya1q3jj8n8pkvl2kjv3pajdyju4hp92cmxnadknd2 cacao 1000000000
      add_account maya10sy79jhw9hw9sqwdgu0k4mw4qawzl7czewzs47 cacao 1000000000
      add_account maya1gv85v0jvc0rsjunku3qxempax6kmrg5jqh8vmg cacao 1000000000
      add_account maya1vm43yk3jq0evzn2u6a97mh2k9x4xf5mzp62g23 cacao 1000000000
      add_account maya1g8dzs4ywxhf8hynaddw4mhwzlwzjfccakkfch7 cacao 1000000000

      add_lp "BTC.BTC" "maya1q3jj8n8pkvl2kjv3pajdyju4hp92cmxnadknd2" "" 1000000000000 "maya1q3jj8n8pkvl2kjv3pajdyju4hp92cmxnadknd2"
      add_lp "BTC.BTC" "maya10sy79jhw9hw9sqwdgu0k4mw4qawzl7czewzs47" "" 1000000000000 "maya10sy79jhw9hw9sqwdgu0k4mw4qawzl7czewzs47"
      add_lp "BTC.BTC" "maya1gv85v0jvc0rsjunku3qxempax6kmrg5jqh8vmg" "" 1000000000000 "maya1gv85v0jvc0rsjunku3qxempax6kmrg5jqh8vmg"
      add_lp "BTC.BTC" "maya1vm43yk3jq0evzn2u6a97mh2k9x4xf5mzp62g23" "" 1000000000000 "maya1vm43yk3jq0evzn2u6a97mh2k9x4xf5mzp62g23"
      add_lp "BTC.BTC" "maya1g8dzs4ywxhf8hynaddw4mhwzlwzjfccakkfch7" "" 1000000000000 "maya1g8dzs4ywxhf8hynaddw4mhwzlwzjfccakkfch7"

      # Add mayanames
      add_mayaname itzamna 5256000 maya18z343fsdlav47chtkyp0aawqt6sgxsh3vjy2vz maya18z343fsdlav47chtkyp0aawqt6sgxsh3vjy2vz MAYA
      add_mayaname aaluxx 5256000 maya1thtn0cdcu2sujtklknlc9wl2dg09a07vtger0x maya1thtn0cdcu2sujtklknlc9wl2dg09a07vtger0x MAYA
      add_mayaname wr 5256000 maya14lkndecaw0zkzu0yq4a0qq869hrs8hh7uhvzlp 0x259fe33B59ae7C24B34C0cb957271Ae0f7786878 ETH
      add_mayaname wr 5256000 maya14lkndecaw0zkzu0yq4a0qq869hrs8hh7uhvzlp maya1a427q3v96psuj4fnughdw8glt5r7j38ljfa6hh MAYA
      add_mayaname wr 5256000 maya14lkndecaw0zkzu0yq4a0qq869hrs8hh7uhvzlp thor1a427q3v96psuj4fnughdw8glt5r7j38lj7rkp8 THOR

      reserve 993994900000000
      # liquidity of genesis nodes
      asgard 6000000000000
    fi

    if [ "$NET" = "stagenet" ]; then
      echo "Adding pool and lp bonder for first node"
      add_pool "BTC.BTC" 8
      add_lp "BTC.BTC" "$NODE_ADDRESS" "" 1000000000000 "$NODE_ADDRESS"

      echo "Setting up accounts"
      add_account smaya18z343fsdlav47chtkyp0aawqt6sgxsh3ctcu6u cacao 9000000000000000
      add_account smaya18z343fsdlav47chtkyp0aawqt6sgxsh3ctcu6u maya 10000000000
      # genesis nodes
      add_account smaya10sy79jhw9hw9sqwdgu0k4mw4qawzl7czdh7xrq cacao 1000000000
      add_account smaya1gv85v0jvc0rsjunku3qxempax6kmrg5j5wm6dk cacao 1000000000
      add_account smaya1vm43yk3jq0evzn2u6a97mh2k9x4xf5mz4rk7u0 cacao 1000000000
      add_account smaya1q3jj8n8pkvl2kjv3pajdyju4hp92cmxnf529m5 cacao 1000000000

      add_lp "BTC.BTC" "smaya10sy79jhw9hw9sqwdgu0k4mw4qawzl7czdh7xrq" "" 10000000000000 "smaya10sy79jhw9hw9sqwdgu0k4mw4qawzl7czdh7xrq"
      add_lp "BTC.BTC" "smaya1gv85v0jvc0rsjunku3qxempax6kmrg5j5wm6dk" "" 10000000000000 "smaya1gv85v0jvc0rsjunku3qxempax6kmrg5j5wm6dk"
      add_lp "BTC.BTC" "smaya1vm43yk3jq0evzn2u6a97mh2k9x4xf5mz4rk7u0" "" 10000000000000 "smaya1vm43yk3jq0evzn2u6a97mh2k9x4xf5mz4rk7u0"
      add_lp "BTC.BTC" "smaya1q3jj8n8pkvl2kjv3pajdyju4hp92cmxnf529m5" "" 10000000000000 "smaya1q3jj8n8pkvl2kjv3pajdyju4hp92cmxnf529m5"

      reserve 994996000000000
    fi

    if [ "$NET" = "testnet" ]; then
      # mint 1m RUNE to reserve for testnet
      reserve 100000000000000

      # add testnet account and balances
      testnet_add_accounts
    fi

    echo "Genesis content"
    cat ~/.mayanode/config/genesis.json
    mayanode validate-genesis --trace
  fi
fi

# setup peer connection, typically only used for some mocknet configurations
if [ "$SEED" != "$(hostname)" ]; then
  if [ ! -f ~/.mayanode/config/genesis.json ]; then
    echo "Setting MAYANode as peer not genesis"

    init_chain "$NODE_ADDRESS"
    fetch_genesis "$SEED"
    NODE_ID=$(fetch_node_id "$SEED")
    echo "NODE ID: $NODE_ID"
    export THOR_TENDERMINT_P2P_PERSISTENT_PEERS="$NODE_ID@$SEED:$PORT_P2P"

    cat ~/.mayanode/config/genesis.json
  fi
fi

# render tendermint and cosmos configuration files
mayanode render-config

export SIGNER_NAME
export SIGNER_PASSWD
exec mayanode start
