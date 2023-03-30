#!/bin/bash

RPC=https://fullnode.testnet.sui.io:443/

DATA='{"jsonrpc":"2.0","id":"1","method":"sui_getCheckpoint","params":["2000"]}'

curl -X POST -H "Content-Type: application/json" \
    -d $DATA \
    $RPC | jq . > getCheckpoint.json
