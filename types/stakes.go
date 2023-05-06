package types

import (
	"math/big"
)

type StakeInfo struct {
	ValidatorAddress string `json:"validatorAddress"`
	StakingPool      string `json:"stakingPool"`
	Stakes           Stakes `json:"stakes"`
}

type StakeInfoList []StakeInfo

func (s StakeInfoList) StakeSum() *Mist {
	total := new(big.Int)

	for _, info := range s {
		total = new(big.Int).Add(total, info.Stakes.Sum().BigInt())
	}

	return &Mist{value: total}
}

type Stake struct {
	StakedSuiId       string `json:"stakedSuiId"`
	StakeRequestEpoch Epoch  `json:"stakeRequestEpoch"`
	StakeActiveEpoch  Epoch  `json:"stakeActiveEpoch"`
	Principal         Mist   `json:"principal"`
	Status            string `json:"status"`
	EstimatedReward   Mist   `json:"estimatedReward"`
}

type Stakes []Stake

func (s Stakes) Sum() *Mist {
	sum := new(big.Int).SetInt64(0)
	for _, stake := range s {
		v := stake.Principal.BigInt()
		sum = new(big.Int).Add(sum, v)
	}

	return &Mist{value: sum}
}

type StakeByEpoch []Stake

func (s StakeByEpoch) Len() int      { return len(s) }
func (s StakeByEpoch) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s StakeByEpoch) Less(i, j int) bool {
	return s[i].StakeActiveEpoch.BigInt().Cmp(s[j].StakeActiveEpoch.BigInt()) < 0
}
