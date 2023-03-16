package types

import (
	"encoding/json"
	"math/big"
)

type Epoch struct {
	value *big.Int
}

func (m *Epoch) UnmarshalJSON(data []byte) error {
	v := new(big.Int)
	_, _ = v.SetString(string(data), 10)
	m.value = v
	return nil
}

func (m *Epoch) String() string {
	if m.value == nil {
		return ""
	}

	return m.value.String()
}

// Mist is smallest unit of Sui token.
type Mist struct {
	value *big.Int
}

func (m *Mist) UnmarshalJSON(data []byte) error {
	v := new(big.Int)
	_, _ = v.SetString(string(data), 10)
	m.value = v
	return nil
}

func (m *Mist) String() string {
	if m.value == nil {
		return ""
	}

	return m.value.String()
}

func (m *Mist) Sui() *big.Float {
	f := new(big.Float)
	f.SetString(m.value.String())

	r := new(big.Float)
	r.Quo(f, big.NewFloat(1000000000.0))

	return r
}

type Checkpoint struct {
	Epoch                      Epoch       `json:"epoch,omitempty"`
	SequenceNumber             json.Number `json:"sequenceNumber,omitempty"`
	Digest                     string      `json:"digest,omitempty"`
	NetworkTotalTransactions   json.Number `json:"networkTotalTransactions,omitempty"`
	PreviousDigest             string      `json:"previousDigest,omitempty"`
	EpochRollingGasCostSummary struct {
		ComputationCost Mist `json:"computationCost,omitempty"`
		StorageCost     Mist `json:"storageCost,omitempty"`
		StorageRebate   Mist `json:"storageRebate,omitempty"`
	} `json:"epochRollingGasCostSummary,omitempty"`
	TimestampMs           uint64   `json:"timestampMs,omitempty"`
	Transactions          []string `json:"transactions,omitempty"`
	CheckpointCommitments []string `json:"checkpointCommitments,omitempty"`
	// EndOfEpochData // TODO:
}

type SuiSystemState struct {
	Epoch                          Epoch  `json:"epoch,omitempty"`
	ProtocolVersion                uint64 `json:"protocolVersion,omitempty"`
	SystemStateVersion             uint64 `json:"systemStateVersion,omitempty"`
	StorageFund                    Mist   `json:"storageFund,omitempty"`
	ReferenceGasPrice              Mist   `json:"referenceGasPrice,omitempty"`
	SafeMode                       bool   `json:"safeMode,omitempty"`
	EpochStartTimestampMs          uint64 `json:"epochStartTimestampMs,omitempty"`
	GovernanceStartEpoch           Epoch  `json:"governanceStartEpoch,omitempty"`
	EpochDurationMs                uint64 `json:"epochDurationMs,omitempty"`
	StakeSubsidyEpochCounter       Epoch  `json:"stakeSubsidyEpochCounter,omitempty"`
	StakeSubsidyBalance            Mist   `json:"stakeSubsidyBalance,omitempty"`
	StakeSubsidyCurrentEpochAmount Mist   `json:"stakeSubsidyCurrentEpochAmount,omitempty"`
	TotalStake                     Mist   `json:"totalStake,omitempty"`

	ActiveValidators []struct {
		SuiAddress             string `json:"suiAddress,omitempty"`
		ProtocolPubkeyBytes    string `json:"protocolPubkeyBytes,omitempty"`
		NetworkPubkeyBytes     string `json:"networkPubkeyBytes,omitempty"`
		WorkerPubkeyBytes      string `json:"workerPubkeyBytes,omitempty"`
		ProofOfPossessionBytes string `json:"proofOfPossessionBytes,omitempty"`
		Name                   string `json:"name,omitempty"`
		Description            string `json:"description,omitempty"`
		ImageUrl               string `json:"imageUrl,omitempty"`
		ProjectUrl             string `json:"projectUrl,omitempty"`

		NetAddress     string `json:"netAddress,omitempty"`
		P2PAddress     string `json:"p2pAddress,omitempty"`
		PrimaryAddress string `json:"primaryAddress,omitempty"`
		WorkerAddress  string `json:"workerAddress,omitempty"`

		NextEpochProtocolPubkeyBytes string `json:"nextEpochProtocolPubkeyBytes,omitempty"`
		NextNetworkPubkeyBytes       string `json:"nextNetworkPubkeyBytes,omitempty"`
		NextWorkerPubkeyBytes        string `json:"nextWorkerPubkeyBytes,omitempty"`
		NextProofOfPossessionBytes   string `json:"nextProofOfPossessionBytes,omitempty"`
		NextNetAddress               string `json:"nextNetAddress,omitempty"`
		NextP2PAddress               string `json:"nextP2PAddress,omitempty"`
		NextPrimaryAddress           string `json:"nextPrimaryAddress,omitempty"`
		NextWorkerAddress            string `json:"nextWorkerAddress,omitempty"`

		VotingPower    uint32 `json:"votingPower,omitempty"`
		OperationCapID string `json:"operationCapId,omitempty"`
		GasPrice       Mist   `json:"gasPrice,omitempty"`
		CommissionRate uint32 `json:"commissionRate,omitempty"`

		NextEpochStake          Mist   `json:"nextEpochStake,omitempty"`
		NextEpochGasPrice       Mist   `json:"nextEpochGasPrice,omitempty"`
		NextEpochCommissionRate uint32 `json:"nextEpochCommissionRate,omitempty"`

		StakingPoolID                string `json:"stakingPoolId,omitempty"`
		StakingPoolActivationEpoch   Epoch  `json:"stakingPoolActivationEpoch,omitempty"`
		StakingPoolDeactivationEpoch Epoch  `json:"stakingPoolDeactivationEpoch,omitempty"`
		StakingPoolSuiBalance        Mist   `json:"stakingPoolSuiBalance,omitempty"`

		RewardsPool      Mist `json:"rewardsPool,omitempty"`
		PoolTokenBalance Mist `json:"poolTokenBalance,omitempty"`

		PendingStake             Mist `json:"pendingStake,omitempty"`
		PendingTotalSuiWithdraw  Mist `json:"pendingTotalSuiWithdraw,omitempty"`
		PendingPoolTokenWithdraw Mist `json:"pendingPoolTokenWithdraw,omitempty"`

		ExchangeRatesID   string `json:"exchangeRatesId,omitempty"`
		ExchangeRatesSize uint64 `json:"exchangeRatesSize,omitempty"`
	} `json:"activeValidators,omitempty"`

	PendingActiveValidatorsID   string `json:"pendingActiveValidatorsId,omitempty"`
	PendingActiveValidatorsSize uint64 `json:"pendingActiveValidatorsSize,omitempty"`
	// PendingRemovals: [],
	StakingPoolMappingsId   string `json:"stakingPoolMappingsId,omitempty"`
	StakingPoolMappingsSize uint64 `json:"stakingPoolMappingsSize,omitempty"`
	InactivePoolsId         string `json:"inactivePoolsId,omitempty"`
	InactivePoolsSize       uint64 `json:"inactivePoolsSize,omitempty"`
	ValidatorCandidatesId   string `json:"validatorCandidatesId,omitempty"`
	ValidatorCandidatesSize uint64 `json:"validatorCandidatesSize,omitempty"`
	// AtRiskValidators: [],
	// ValidatorReportRecords: []
}
