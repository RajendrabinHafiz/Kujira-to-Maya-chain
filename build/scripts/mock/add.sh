#!/bin/sh

# ./mock-add.bash <mock binance IP address> <thornode peer> <bnb address> <amt> <asset>
# ./mock-add.bash 127.0.0.1 bnbZZZZZ 3000 RUNE-A1F

if [ -z "$1" ]; then
  echo "Missing mock binance address (address:port)"
  exit 1
fi

if [ -z "$2" ]; then
  echo "Missing peer thorchain"
fi

if [ -z "$3" ]; then
  echo "Missing bnb address argument"
  exit 1
fi

if [ -z "$4" ]; then
  echo "Missing thorchain address argument"
  exit 1
fi

if [ -z "$5" ]; then
  echo "Missing amount argument"
  exit 1
fi

if [ -z "$6" ]; then
  echo "Missing asset"
  exit 1
fi

INBOUND_ADDRESS=$(curl -s "$2:1317/thorchain/inbound_addresses" | jq -r '.[]|select(.chain=="BNB") | .address')

echo "$1"
echo "$2"
echo "$3"
echo "$INBOUND_ADDRESS"
echo "$4"
echo "$5"
echo "$6"

curl -vvv -s -X POST -d "[{
  \"from\": \"$3\",
  \"to\": \"$INBOUND_ADDRESS\",
  \"coins\":[
      {\"denom\": \"$6\", \"amount\": $5}
  ],
  \"memo\": \"ADD:BNB.$6:$4\"
}]" "$1:26660/broadcast/easy"
