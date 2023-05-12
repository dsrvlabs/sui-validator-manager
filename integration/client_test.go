package integration

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/stretchr/testify/assert"
)

func TestLatestCheckpointSequenceNumber(t *testing.T) {
	var rpcURL = "https://sui-rpc-mainnet.testnet-pride.com:443"
	cli := rpc.NewClient([]string{rpcURL})

	_, err := cli.LatestCheckpointSequenceNumber()

	assert.Nil(t, err)
}

func TestCheckpoint(t *testing.T) {
	var rpcURL = "https://sui-rpc-mainnet.testnet-pride.com:443"
	cli := rpc.NewClient([]string{rpcURL})

	cp, err := cli.Checkpoint(new(big.Int))

	fmt.Println(cp, err)
	assert.Nil(t, err)
}

func TestLatestSuiSystemState(t *testing.T) {
	var rpcURL = "https://sui-rpc-mainnet.testnet-pride.com:443"
	cli := rpc.NewClient([]string{rpcURL})

	state, err := cli.LatestSuiSystemState()

	assert.Nil(t, err)

	_ = state
}
