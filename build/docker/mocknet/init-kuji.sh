#!/bin/sh

# /gaiad init --chain-id harpoon-4 local

mkdir -p /root/.kuji/data
cat >/root/.kuji/data/priv_validator_state.json <<EOF
{
  "height": "0",
  "round": 0,
  "step": 0
}
EOF

exec /entrypoint.sh
