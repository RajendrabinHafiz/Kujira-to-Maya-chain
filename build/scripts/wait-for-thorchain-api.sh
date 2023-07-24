#!/bin/sh

# https://docs.docker.com/compose/startup-order/

set -e

echo "Waiting for $2..."

until curl -s "$1/thorchain/ping" >/dev/null; do
  # echo "Rest server is unavailable - sleeping"
  sleep 1
done

echo "$2 ready"
