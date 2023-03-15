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
	assert.Equal(t, "4uYJt5NcXCxoZFRehsNeZn8dtrUujmn7qhrFaV4Zt1as", cp.Digest)
	assert.Equal(t, "2001", cp.NetworkTotalTransactions.String())
	assert.Equal(t, "FrtoQ9tACumuxb83iiSdwPW2Q122KRHubcHKgjKDRxEb", cp.PreviousDigest)
	assert.Equal(t, "1", cp.EpochRollingGasCostSummary.ComputationCost.String())
	assert.Equal(t, "2", cp.EpochRollingGasCostSummary.StorageCost.String())
	assert.Equal(t, "3", cp.EpochRollingGasCostSummary.StorageRebate.String())

	assert.Equal(t, uint64(1678893005933), cp.TimestampMs)
	assert.Equal(t, "2vXkHf76Y3CCAgUQnYiNXwnHzBSTdQrBUKXpYBg31m1Q", cp.Transactions[0])
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

	assert.Equal(t, "1", state.Epoch.String())
	assert.Equal(t, uint64(1), state.ProtocolVersion)
	assert.Equal(t, uint64(1), state.SystemStateVersion)
	assert.Equal(t, "13164", state.StorageFund.String())
	assert.Equal(t, "1000", state.ReferenceGasPrice.String())
	assert.False(t, state.SafeMode)
	assert.Equal(t, uint64(1678922854471), state.EpochStartTimestampMs)
	assert.Equal(t, "0", state.GovernanceStartEpoch.String())
	assert.Equal(t, uint64(86400000), state.EpochDurationMs)
	assert.Equal(t, "1", state.StakeSubsidyEpochCounter.String())
	assert.Equal(t, "99000000000000000", state.StakeSubsidyBalance.String())
	assert.Equal(t, "1000000000000000", state.StakeSubsidyCurrentEpochAmount.String())
	assert.Equal(t, "1901000000041827000", state.TotalStake.String())

	assert.True(t, len(state.ActiveValidators) > 0)

	v := state.ActiveValidators[0]

	assert.Equal(t, "0x8925c11a13cf4b30a64a30ee9f3ca401e58b541b34517d99122e779aa81e3bc9", v.SuiAddress)
	assert.Equal(
		t,
		"mBGzDVHKDA2mooaQt7hdjXyCupQQEt5HHTWUkTF4w92nvKTsF448a2CAU1edYPyTAyn4Ph6vD5H6zBNQQy5JRLBY9kL44N1LG2NvF1GjgDU6N6pidnt2eXKZdhVPvTSxHTc",
		v.ProtocolPubkeyBytes,
	)
	assert.Equal(t, "ncCAUxd9tDLAz7pjWAFkApjgr93o3QqGffdofeDgnzc", v.NetworkPubkeyBytes)
	assert.Equal(t, "DLKEdhMQHNDu6eCTkvgRMpMcwBCwDbRiFFy4u78WC58S", v.WorkerPubkeyBytes)
	assert.Equal(t, "6S4fC8eRLSdiX1Rj6gsGJBSwhEZ9PoCHFw5Tg6dUH8Jw2LY1KRnFPawCQbiq1AEzg6", v.ProofOfPossessionBytes)
	assert.Equal(t, "Overclock", v.Name)
	assert.Equal(t, "High uptime, high performance bare metal validator", v.Description)
	assert.Equal(
		t,
		"https://avatars.githubusercontent.com/u/116830051?s=400&u=ee4d79d2bde09a2aa216980a198b5dabccc05545&v=4",
		v.ImageUrl,
	)
	assert.Equal(t, "https://www.overclock.one", v.ProjectUrl)

	assert.Equal(t, "/ip4/186.233.186.11/tcp/8080/http", v.NetAddress)
	assert.Equal(t, "/ip4/186.233.186.11/udp/8084", v.P2PAddress)
	assert.Equal(t, "/ip4/186.233.186.11/udp/8081", v.PrimaryAddress)
	assert.Equal(t, "/ip4/186.233.186.11/udp/8082", v.WorkerAddress)

	assert.Equal(t, uint32(157), v.VotingPower)
	assert.Equal(t, "0x16bb51edad1714f14cae3a3b1265f619900cee6126bafed1e743e68f2928a536", v.OperationCapID)
	assert.Equal(t, "1000", v.GasPrice.String())
	assert.Equal(t, uint32(1500), v.CommissionRate)

	///
	assert.Equal(t, "30016666667363784", v.NextEpochStake.String())
	assert.Equal(t, "1000", v.NextEpochGasPrice.String())
	assert.Equal(t, uint32(1500), v.NextEpochCommissionRate)
	assert.Equal(t, "0xb0cf62387a5526d51debd3ee28576c60417ba2e8187b94f7d311a312ba0c4c08", v.StakingPoolID)

	assert.Equal(t, "0", v.StakingPoolActivationEpoch.String())
	assert.Equal(t, "0", v.StakingPoolDeactivationEpoch.String())
	assert.Equal(t, "30016666667363784", v.StakingPoolSuiBalance.String())
	assert.Equal(t, "14166667259216", v.RewardsPool.String())

	assert.Equal(t, "30002498820106136", v.PoolTokenBalance.String())
	assert.Equal(t, "0", v.PendingStake.String())
	assert.Equal(t, "0", v.PendingTotalSuiWithdraw.String())
	assert.Equal(t, "0", v.PendingPoolTokenWithdraw.String())
	assert.Equal(t, "0x4d3f360db730ab473ac315c7d5c96840698e1bf49b8e5bb3d49047c8d131519b", v.ExchangeRatesID)
	assert.Equal(t, uint64(2), v.ExchangeRatesSize)
}
