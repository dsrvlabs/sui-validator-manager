package rgp

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dsrvlabs/sui-validator-manager/types"
)

func TestRGP_RewardCalculation(t *testing.T) {
	systemState := types.SuiSystemState{
		TotalStake:                            new(types.Mist).SetString("4720871917040883702"),
		StakeSubsidyCurrentDistributionAmount: new(types.Mist).SetString("1111111111111111"),
		ActiveValidators:                      make([]types.Validator, 100),
	}

	testSuiAddress := "0x6f4e73ee97bfae95e054d31dff1361a839aaadf2cfdb873ad2b07d479507905a"
	validator := types.Validator{
		SuiAddress: testSuiAddress,
		VotingPower:           json.Number("56"),
		CommissionRate:        json.Number("1000"),
		StakingPoolSuiBalance: new(types.Mist).SetString("26020681646333831"),
	}

	systemState.ActiveValidators[0] = validator

	// TODO: StoragefuncNonRefuncdableBalance + StorageFuncTotalObjectStorageRebate
	storageFund, _ := new(big.Int).SetString("2011080000000", 10)
	computationCost, _ := new(big.Int).SetString("2848461056220", 10)

	selfStake, _ := new(big.Int).SetString("124472926677", 10)

	reward, err := CalculateReward(
		systemState,
		storageFund,
		computationCost,
		selfStake,
		testSuiAddress,
	)

	assert.Nil(t, err)
	assert.NotNil(t, reward)

	// TODO: Result value is not accurate.
	assert.Equal(t, "623848459744", reward.String())
}
