package rpc

import (
	"io"
	"math/big"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, "9QSYond1gVerejrBgvnZxC5CD3jxKcsubDZdcZni5Ua2", cp.Digest)
	assert.Equal(t, "2001", cp.NetworkTotalTransactions.String())
	assert.Equal(t, "5LXvxVpbvCwLFArQ74gxmQg4LLa9kQmv9FtReR1FmPhh", cp.PreviousDigest)
	assert.Equal(t, "1", cp.EpochRollingGasCostSummary.ComputationCost.String())
	assert.Equal(t, "2", cp.EpochRollingGasCostSummary.StorageCost.String())
	assert.Equal(t, "3", cp.EpochRollingGasCostSummary.StorageRebate.String())
	assert.Equal(t, "4", cp.EpochRollingGasCostSummary.NonRefundableStorageFee.String())

	assert.Equal(t, uint64(1680014365060), cp.TimestampMs)
	assert.Equal(t, "7dNSFkcVAm4GDQ1wPsMC1W2JpWtnHYCaT6Wq7UvH3naV", cp.Transactions[0])
	assert.Equal(t, []string{}, cp.CheckpointCommitments)
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

	assert.Equal(t, "2", state.Epoch.String())
	assert.Equal(t, uint64(1), state.ProtocolVersion)
	assert.Equal(t, uint64(1), state.SystemStateVersion)
	assert.Equal(t, "21772", state.StorageFundTotalObjectStorageRebate.String())
	assert.Equal(t, "104", state.StorageFundNonRefundableBalance.String())
	assert.Equal(t, "1000", state.ReferenceGasPrice.String())
	assert.False(t, state.SafeMode)
	assert.Equal(t, "0", state.SafeModeStorageRewards.String())
	assert.Equal(t, "0", state.SafeModeComputationRewards.String())
	assert.Equal(t, "0", state.SafeModeNonRefundableStorageFee.String())
	assert.Equal(t, uint64(1680109210601), state.EpochStartTimestampMs)
	assert.Equal(t, uint64(86400000), state.EpochDurationMs)
	assert.Equal(t, "2", state.StakeSubsidyStartEpoch.String())
	assert.Equal(t, uint64(150), state.MaxValidatorCount)
	assert.Equal(t, "3", state.MinValidatorJoiningStake.String()) // TODO
	assert.Equal(t, "2", state.ValidatorLowStakeThreshold.String()) // TODO
	assert.Equal(t, "15000000000000000", state.ValidatorVeryLowStakeThreshold.String())
	assert.Equal(t, "7", state.ValidatorLowStakeGracePeriod.String())
	assert.Equal(t, "1", state.StakeSubsidyBalance.String()) // TODO
	assert.Equal(t, "0", state.StakeSubsidyDistributionCounter.String())
	assert.Equal(t, "10", state.StakeSubsidyPeriodLength.String())
	assert.Equal(t, uint64(1000), state.StakeSubsidyDecreaseRate)
	assert.Equal(t, "4635000185940212000", state.TotalStake.String())

	assert.True(t, len(state.ActiveValidators) > 0)

	v := state.ActiveValidators[56]

	assert.Equal(t, "0x0a392298244ca2694098d015b00cf49ae1168118b28d13cb0baafd5884e5559a", v.SuiAddress)
	assert.Equal(
		t,
		"qF1a3fDbUaqZg0YczpP6ip/lhRJMNWJhouhbIP6h2v5l0wAqj74LCB7C/07PKq95BVzj7xqcm9XPwcyn0WPZksw/1BgjwtC7j3Y0gwY9qlNIvbnpNB1Zjbcaak0p5zxO",
		v.ProtocolPubkeyBytes,
	)
	assert.Equal(t, "A+8uE0bqPFltKslTOCrnTH5VBdXsWFL82aQ5KCrd46Y=", v.NetworkPubkeyBytes)
	assert.Equal(t, "IRH0i9Jit+Bjdc8SaAakqMMM2tHEkQlZtAl4SmBofjE=", v.WorkerPubkeyBytes)
	assert.Equal(t, "hA3IGtzySIZyKIZTAaEzs1ymZ8+L0UKNH2bwq0MjAk7rLPK9JFJmvgLH96X3zPMb", v.ProofOfPossessionBytes)
	assert.Equal(t, "DSRV", v.Name)
	assert.Equal(t, "Everything distributed, served complete.", v.Description)
	assert.Equal(
		t,
		"https://raw.githubusercontent.com/dsrvlabs/sui/main/logo/favicon_black_square.svg",
		v.ImageUrl,
	)
	assert.Equal(t, "https://dsrvlabs.com", v.ProjectUrl)

	assert.Equal(t, "/dns/sui-01.validator.dsrvlabs.dev/tcp/8080/http", v.NetAddress)
	assert.Equal(t, "/dns/sui-01.validator.dsrvlabs.dev/udp/8084", v.P2PAddress)
	assert.Equal(t, "/dns/sui-01.validator.dsrvlabs.dev/udp/8081", v.PrimaryAddress)
	assert.Equal(t, "/dns/sui-01.validator.dsrvlabs.dev/udp/8082", v.WorkerAddress)

	assert.Equal(t, uint32(43), v.VotingPower)
	assert.Equal(t, "0x894d7bba53117c88e9711b47fab8bec4b0def3990109601b9f904470ea294c5b", v.OperationCapID)
	assert.Equal(t, "900", v.GasPrice.String())
	assert.Equal(t, uint32(1000), v.CommissionRate)

	assert.Equal(t, "20000186400768504", v.NextEpochStake.String())
	assert.Equal(t, "900", v.NextEpochGasPrice.String())
	assert.Equal(t, uint32(1000), v.NextEpochCommissionRate)
	assert.Equal(t, "0x093136a86b72b6aa1c84e84e72a00ca2260246441976f1ce070b136dbfc6b90f", v.StakingPoolID)
	assert.Equal(t, "0", v.StakingPoolActivationEpoch.String())
	assert.Equal(t, "0", v.StakingPoolDeactivationEpoch.String())
	assert.Equal(t, "20000000008628504", v.StakingPoolSuiBalance.String())
	assert.Equal(t, "7765656", v.RewardsPool.String())

	assert.Equal(t, "20000000000862850", v.PoolTokenBalance.String())
	assert.Equal(t, "186392140000", v.PendingStake.String())
	assert.Equal(t, "0", v.PendingTotalSuiWithdraw.String())
	assert.Equal(t, "0", v.PendingPoolTokenWithdraw.String())
	assert.Equal(t, "0xbbf6cb712a8cab65919d29cac95fadbfec80e2118ade0580a7e2653c49706ee2", v.ExchangeRatesID)
	assert.Equal(t, uint64(3), v.ExchangeRatesSize)

	assert.Equal(t, "0x57efe3d331c728ddac6baf7a0474fbfe14d1e65a314fb56c6024a068c24777cb", state.PendingActiveValidatorsID)
	assert.Equal(t, uint64(1), state.PendingActiveValidatorsSize)

	assert.Equal(t, "0xbc85c779db3a0fe9afc7747f070000829857a2dabb987e0ea48aeed55db86d37", state.StakingPoolMappingsId)
	assert.Equal(t, uint64(95), state.StakingPoolMappingsSize)
	assert.Equal(t, "0x893c9de72b8c59e666e172014e6cd817ba930c584171f100c5124a3741db237d", state.InactivePoolsId)
	assert.Equal(t, uint64(2), state.InactivePoolsSize)
	assert.Equal(t, "0x289405ce92b7479287ef74dde807caaf3341e659c39ef5bf43f39236d12c73d2", state.ValidatorCandidatesId)
	assert.Equal(t, uint64(0), state.ValidatorCandidatesSize)
}
