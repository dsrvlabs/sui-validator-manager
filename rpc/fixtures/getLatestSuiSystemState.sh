#!/bin/bash

RPC=https://fullnode.testnet.sui.io:443/

DATA='{"jsonrpc":"2.0","id":"1","method":"suix_getLatestSuiSystemState","params":[]}'

curl -X POST -H "Content-Type: application/json" \
    -d $DATA \
    $RPC | jq . > getLatestSuiSystemState.json
