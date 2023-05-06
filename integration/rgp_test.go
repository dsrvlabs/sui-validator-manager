package integration

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsrvlabs/sui-validator-manager/rgp"
	"github.com/dsrvlabs/sui-validator-manager/rpc"
)

func TestRGP_CalculateReward(t *testing.T) {
	validatorAddr := "0x6f4e73ee97bfae95e054d31dff1361a839aaadf2cfdb873ad2b07d479507905a"
	url := "https://fullnode.mainnet.sui.io:443"

	cli := rpc.NewClient([]string{url})

	state, err := cli.LatestSuiSystemState()

	assert.Nil(t, err)

	storageFund, _ := new(big.Int).SetString("8386140000000", 10)
	computationCost, _ := new(big.Int).SetString("3852472157232", 10)

	selfStake, _ := new(big.Int).SetString("8491585775550", 10)

	reward, err :=rgp.CalculateReward(*state, storageFund, computationCost, selfStake, validatorAddr)

	fmt.Println(reward, err)

	assert.Nil(t, err)

	_ = reward
}
