package rgp

import (
	"fmt"
	"math/big"

	"github.com/dsrvlabs/sui-validator-manager/types"
)

/**
TODOs
- Add tally score.
- gamma value from onchain.
- Calculate adjusted voting power.
*/
func CalculateReward(
	systemState types.SuiSystemState,
	storageFund,
	computationCost,
	selfStake *big.Int,
	validatorAddress string,
) (*big.Int, error) {
	var validator *types.Validator
	for _, val := range systemState.ActiveValidators {
		if val.SuiAddress == validatorAddress {
			validator = &val
			break
		}
	}

	if validator == nil {
		return nil, fmt.Errorf("%s is not an active validator", validatorAddress)
	}

	gamma := 0.95
	tallyScore := 1.0 // mu := 1.0

	_totalStake := new(big.Float).SetInt(systemState.TotalStake.BigInt())
	_storageFund := new(big.Float).SetInt(storageFund)
	totalFund := new(big.Float).Add(_totalStake, _storageFund)

	_alpha := new(big.Float).Quo(_totalStake, totalFund)
	alpha, accuracy := _alpha.Float64()

	_selfStake := new(big.Float).SetInt(selfStake)
	_validatorStake := new(big.Float).SetInt(validator.StakingPoolSuiBalance.BigInt())

	_selfStakeRatio := new(big.Float).Quo(_selfStake, _validatorStake)
	selfStakeRatio, accuracy := _selfStakeRatio.Float64()

	stakeReward := new(big.Int).Add(
		systemState.StakeSubsidyCurrentDistributionAmount.BigInt(),
		computationCost)

	commissionRate, err := validator.CommissionRate.Float64()
	if err != nil {
		return nil, err
	}
	commissionRate /= 10000

	validatorShare, err := validator.VotingPower.Float64()
	if err != nil {
		return nil, err
	}
	validatorShare /= 10000

	k := alpha * (selfStakeRatio + commissionRate*(1.0-selfStakeRatio)) * tallyScore * validatorShare
	k += (1 - alpha) * (gamma / float64(len(systemState.ActiveValidators)))

	validatorReward := new(big.Float).Mul(big.NewFloat(k), new(big.Float).SetInt(stakeReward))

	resultReward := new(big.Int)
	_, accuracy = validatorReward.Int(resultReward)

	_ = accuracy

	return resultReward, nil
}
