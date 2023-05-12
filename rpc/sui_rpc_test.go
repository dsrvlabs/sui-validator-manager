package rpc

import (
	"io"
	"math/big"
	"net/http"
	"os"
	"sort"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"

	"github.com/dsrvlabs/sui-validator-manager/types"
)

const (
	url = "https://fullnode.devnet.sui.io:443"
)

func TestGetLatestCheckpointSequeceNumber(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	fixture := "./fixtures/getLatestCheckpointSequenceNumber.json"

	f, err := os.Open(fixture)
	defer f.Close()

	assert.Nil(t, err)

	data, err := io.ReadAll(f)
	assert.Nil(t, err)

	httpmock.RegisterResponder(
		http.MethodPost,
		url,
		httpmock.NewStringResponder(http.StatusOK, string(data)))

	cli := NewClient([]string{url})
	num, err := cli.LatestCheckpointSequenceNumber()

	assert.Nil(t, err)
	assert.Zero(t, big.NewInt(903369).Cmp(num))
	assert.NotNil(t, num)
}

func TestCheckpoint(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	fixture := "./fixtures/getCheckpoint.json"

	f, err := os.Open(fixture)
	defer f.Close()

	assert.Nil(t, err)

	data, err := io.ReadAll(f)
	assert.Nil(t, err)

	cli := NewClient([]string{url})

	httpmock.RegisterResponder(
		http.MethodPost,
		url,
		httpmock.NewStringResponder(http.StatusOK, string(data)))

	cp, err := cli.Checkpoint(big.NewInt(0)) // Simple use dummy argument.

	assert.Nil(t, err)
	assert.Equal(t, "0", cp.Epoch.String())
	assert.Equal(t, "2000", cp.SequenceNumber.String())
	assert.Equal(t, "4tvXxeUpTU1XYw8DqvZBYfZkTek5CDPR3myuyBA6vwCH", cp.Digest)
	assert.Equal(t, "2001", cp.NetworkTotalTransactions.String())
	assert.Equal(t, "GF9J2Xv8HUUcx8yH5CsVZYRxvn88jWXqP71HoxjMmxcv", cp.PreviousDigest)
	assert.Equal(t, "1", cp.EpochRollingGasCostSummary.ComputationCost.String())
	assert.Equal(t, "2", cp.EpochRollingGasCostSummary.StorageCost.String())
	assert.Equal(t, "3", cp.EpochRollingGasCostSummary.StorageRebate.String())
	assert.Equal(t, "4", cp.EpochRollingGasCostSummary.NonRefundableStorageFee.String())

	assert.Equal(t, int64(1681394996136)/1000, cp.TimestampMs.Unix())
	assert.Equal(t, "DZDJjt92zFt8xjnm81ooN84R9m6fF5tNMdL71hwMLATB", cp.Transactions[0])
	assert.Equal(t, []string{}, cp.CheckpointCommitments)
	assert.Equal(t, "jN7LF69ln0JdZ5BDD+o/Rl1qtDcSx798m1t/ASZibdSSqBLjgzZlsqDU/Zh/j8JR",
		cp.ValidatorSignature)
}

func TestLatestSuiSystemState(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	fixture := "./fixtures/getLatestSuiSystemState.json"

	f, err := os.Open(fixture)
	defer f.Close()

	assert.Nil(t, err)

	data, err := io.ReadAll(f)
	assert.Nil(t, err)

	cli := NewClient([]string{url})

	httpmock.RegisterResponder(
		http.MethodPost,
		url,
		httpmock.NewStringResponder(http.StatusOK, string(data)))

	state, err := cli.LatestSuiSystemState() // Simple use dummy argument.

	assert.Nil(t, err)

	assert.Equal(t, "1", state.Epoch.String())
	assert.Equal(t, "4", state.ProtocolVersion)
	assert.Equal(t, "2", state.SystemStateVersion)
	assert.Equal(t, "0", state.StorageFundTotalObjectStorageRebate.String())
	assert.Equal(t, "0", state.StorageFundNonRefundableBalance.String())
	assert.Equal(t, "1000", state.ReferenceGasPrice.String())
	assert.False(t, state.SafeMode)
	assert.Equal(t, "0", state.SafeModeStorageRewards.String())
	assert.Equal(t, "0", state.SafeModeComputationRewards.String())
	assert.Equal(t, "0", state.SafeModeNonRefundableStorageFee.String())
	assert.Equal(t, int64(1681405202097)/1000, state.EpochStartTimestampMs.Unix())
	assert.Equal(t, "86400000", state.EpochDurationMs.String())
	assert.Equal(t, "13", state.StakeSubsidyStartEpoch.String())
	assert.Equal(t, "150", state.MaxValidatorCount.String())
	assert.Equal(t, "30000000000000000", state.MinValidatorJoiningStake.String())
	assert.Equal(t, "20000000000000000", state.ValidatorLowStakeThreshold.String())
	assert.Equal(t, "15000000000000000", state.ValidatorVeryLowStakeThreshold.String())
	assert.Equal(t, "7", state.ValidatorLowStakeGracePeriod.String())
	assert.Equal(t, "1000000000000000000", state.StakeSubsidyBalance.String())
	assert.Equal(t, "0", state.StakeSubsidyDistributionCounter.String())
	assert.Equal(t, "90", state.StakeSubsidyPeriodLength.String())
	assert.Equal(t, uint64(1000), state.StakeSubsidyDecreaseRate)
	assert.Equal(t, "4655000000000000000", state.TotalStake.String())

	assert.True(t, len(state.ActiveValidators) > 0)

	var v types.Validator
	for _, curValidator := range state.ActiveValidators {
		if curValidator.Name == "DSRV" {
			v = curValidator
			break
		}
	}

	assert.Equal(t, "0x6f4e73ee97bfae95e054d31dff1361a839aaadf2cfdb873ad2b07d479507905a", v.SuiAddress)
	assert.Equal(
		t,
		"qVVlhyH6YXw07KX09Rv0aPOFV6PpxaSncgna7RO8CQ860wLZte1MU/m6OGEPR1T7AXepfKPvD1CdXLZXYqJbsnQ3vYl/yZ2yunW9Ydd1huUv0vJ5pvOE/gh3kkB5tVx9",
		v.ProtocolPubkeyBytes,
	)
	assert.Equal(t, "ysDnrm+FsRxOJh7gcacb6E6YpX35SZqCVPOw83NBUww=", v.NetworkPubkeyBytes)
	assert.Equal(t, "KOCVurcFP1f2OaxJT2uLyT7c2Estm6FJoixdcqxn2cg=", v.WorkerPubkeyBytes)
	assert.Equal(t, "pb2SfjvcFtUDNeLKDWoD5o2YI0V21oyMvPmmnt7pf6Bn8sDGrczDLpv5gLV95WHa", v.ProofOfPossessionBytes)
	assert.Equal(t, "DSRV", v.Name)
	assert.Equal(t, "Everything distributed, served complete.", v.Description)
	assert.Equal(
		t,
		"https://raw.githubusercontent.com/dsrvlabs/sui/main/logo/favicon_black_square.svg",
		v.ImageUrl,
	)
	assert.Equal(t, "https://dsrvlabs.com", v.ProjectUrl)

	assert.Equal(t, "/dns/validator-01.sui.dsrvlabs.net/tcp/8080/http", v.NetAddress)
	assert.Equal(t, "/dns/validator-01.sui.dsrvlabs.net/udp/8084", v.P2PAddress)
	assert.Equal(t, "/dns/validator-01.sui.dsrvlabs.net/udp/8081", v.PrimaryAddress)
	assert.Equal(t, "/dns/validator-01.sui.dsrvlabs.net/udp/8082", v.WorkerAddress)

	assert.Equal(t, "54", v.VotingPower.String())
	assert.Equal(t, "0x018a1c3e90d7db8632260f7ed4b751ade6d75354bb05ed57e8834cafd1936427", v.OperationCapID)
	assert.Equal(t, "1000", v.GasPrice.String())
	assert.Equal(t, "200", v.CommissionRate.String())

	assert.Equal(t, "25000000000000000", v.NextEpochStake.String())
	assert.Equal(t, "1000", v.NextEpochGasPrice.String())
	assert.Equal(t, "200", v.NextEpochCommissionRate.String())
	assert.Equal(t, "0xca5d9d37e55d5e44bd1db57c4764a764a33859bb2fd8074d67070f4c3714907d", v.StakingPoolID)
	assert.Equal(t, "0", v.StakingPoolActivationEpoch.String())
	assert.Equal(t, "0", v.StakingPoolDeactivationEpoch.String())
	assert.Equal(t, "25000000000000000", v.StakingPoolSuiBalance.String())
	assert.Equal(t, "0", v.RewardsPool.String())

	assert.Equal(t, "25000000000000000", v.PoolTokenBalance.String())
	assert.Equal(t, "0", v.PendingStake.String())
	assert.Equal(t, "0", v.PendingTotalSuiWithdraw.String())
	assert.Equal(t, "0", v.PendingPoolTokenWithdraw.String())
	assert.Equal(t, "0xe3af0012a2802c37e862d07bad78f50348211c311e8c18cfdf5fe1ae997d349a", v.ExchangeRatesID)
	assert.Equal(t, "2", v.ExchangeRatesSize.String())

	assert.Equal(t, "0x719fdd5d050b2a1be364ab385ac3d163b7ac407e234721392d3c716a6332caf3", state.PendingActiveValidatorsID)
	assert.Equal(t, "0", state.PendingActiveValidatorsSize.String())

	assert.Equal(t, "0x3a4ec1afc6f550aa838aa4e823380a2c7c9567cf12e8e4dcc81ea7d411e544c8", state.StakingPoolMappingsId)
	assert.Equal(t, "100", state.StakingPoolMappingsSize.String())
	assert.Equal(t, "0xf2dfc014f09869d512c7965d743e1f513853f492d9c7c0d755597154cb3ff8cd", state.InactivePoolsId)
	assert.Equal(t, "0", state.InactivePoolsSize.String())
	assert.Equal(t, "0x94f89425db4d5bfa8d982d17f5746a1ac35ccb863662ee9486587e1e2d922763", state.ValidatorCandidatesId)
	assert.Equal(t, "0", state.ValidatorCandidatesSize.String())
}

func TestGetStakes(t *testing.T) {
	httpmock.Activate()
	defer httpmock.Deactivate()

	fixture := "./fixtures/getStakes.json"
	f, err := os.Open(fixture)
	defer f.Close()

	data, err := io.ReadAll(f)

	httpmock.RegisterResponder(
		http.MethodPost,
		url,
		httpmock.NewStringResponder(http.StatusOK, string(data)),
	)

	cli := NewClient([]string{url})

	stakeInfo, err := cli.GetStakes("0x6f4e73ee97bfae95e054d31dff1361a839aaadf2cfdb873ad2b07d479507905a")
	assert.Nil(t, err)

	var stakes []types.Stake = stakeInfo[0].Stakes
	sort.Sort(types.StakeByEpoch(stakes))

	assert.Equal(t, "0x6f4e73ee97bfae95e054d31dff1361a839aaadf2cfdb873ad2b07d479507905a", stakeInfo[0].ValidatorAddress)
	assert.Equal(t, "1349674360863", stakeInfo[0].Stakes.Sum().String())
	assert.Equal(t, "1349674360863", stakeInfo.StakeSum().String())
}
