package types

import (
	"encoding/json"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"time"
)

type Epoch struct {
	value *big.Int
}

func (m *Epoch) UnmarshalJSON(data []byte) error {
	v := new(big.Int)
	_, _ = v.SetString(strings.Trim(string(data), "\""), 10)
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
	_, _ = v.SetString(strings.Trim(string(data), "\""), 10)
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

type Timestamp struct {
	time time.Time
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	val, err := strconv.ParseInt(strings.Trim(string(data), "\""), 10, 64)
	if err != nil {
		return err
	}

	t.time = time.UnixMilli(int64(val))

	return nil
}

func (t *Timestamp) Unix() int64 {
	return t.time.Unix()
}

type BigInt struct {
	value *big.Int
}

func (b *BigInt) UnmarshalJSON(data []byte) error {
	n := new(big.Int)
	bn, ok := n.SetString(strings.Trim(string(data), "\""), 10)
	if !ok {
		return errors.New("cannot parse value")
	}

	b.value = bn

	return nil
}

func (b *BigInt) String() string {
	return b.value.String()
}

type Checkpoint struct {
	Epoch                      Epoch       `json:"epoch,omitempty"`
	SequenceNumber             json.Number `json:"sequenceNumber,omitempty"`
	Digest                     string      `json:"digest,omitempty"`
	NetworkTotalTransactions   json.Number `json:"networkTotalTransactions,omitempty"`
	PreviousDigest             string      `json:"previousDigest,omitempty"`
	EpochRollingGasCostSummary struct {
		ComputationCost         Mist `json:"computationCost,omitempty"`
		StorageCost             Mist `json:"storageCost,omitempty"`
		StorageRebate           Mist `json:"storageRebate,omitempty"`
		NonRefundableStorageFee Mist `json:"nonRefundableStorageFee"`
	} `json:"epochRollingGasCostSummary,omitempty"`
	TimestampMs           Timestamp `json:"timestampMs,omitempty"`
	Transactions          []string  `json:"transactions,omitempty"`
	CheckpointCommitments []string  `json:"checkpointCommitments,omitempty"`
	ValidatorSignature    string    `json:"validatorSignature,omitempty"`
	// EndOfEpochData // TODO:
}

type SuiSystemState struct {
	Epoch                               Epoch  `json:"epoch,omitempty"`
	ProtocolVersion                     string `json:"protocolVersion,omitempty"`
	SystemStateVersion                  string `json:"systemStateVersion,omitempty"`
	StorageFundTotalObjectStorageRebate Mist   `json:"storageFundTotalObjectStorageRebates,omitempty"`
	StorageFundNonRefundableBalance     Mist   `json:"storageFundNonRefundableBalance,omitempty"`
	ReferenceGasPrice                   Mist   `json:"referenceGasPrice,omitempty"`

	SafeMode                        bool `json:"safeMode,omitempty"`
	SafeModeStorageRewards          Mist `json:"safeModeStorageRewards,omitempty"`
	SafeModeComputationRewards      Mist `json:"safeModeComputationRewards,omitempty"`
	SafeModeNonRefundableStorageFee Mist `json:"safeModeNonRefundableStorageFee,omitempty"`

	EpochStartTimestampMs  Timestamp   `json:"epochStartTimestampMs,omitempty"`
	EpochDurationMs        json.Number `json:"epochDurationMs,omitempty"`
	StakeSubsidyStartEpoch Epoch       `json:"stakeSubsidyStartEpoch,omitempty"`
	MaxValidatorCount      json.Number `json:"maxValidatorCount,omitempty"`

	MinValidatorJoiningStake   Mist `json:"minValidatorJoiningStake"`
	ValidatorLowStakeThreshold Mist `json:"validatorLowStakeThreshold"`

	ValidatorVeryLowStakeThreshold Mist  `json:"validatorVeryLowStakeThreshold"`
	ValidatorLowStakeGracePeriod   Epoch `json:"validatorLowStakeGracePeriod,omitempty"`

	StakeSubsidyBalance                   Mist   `json:"stakeSubsidyBalance,omitempty"`
	StakeSubsidyDistributionCounter       Epoch  `json:"stakeSubsidyDistributionCounter"`
	StakeSubsidyCurrentDistributionAmount Mist   `json:"stakeSubsidyCurrentDistributionAmount"`
	StakeSubsidyPeriodLength              Epoch  `json:"stakeSubsidyPeriodLength"`
	StakeSubsidyDecreaseRate              uint64 `json:"stakeSubsidyDecreaseRate"`

	TotalStake Mist `json:"totalStake,omitempty"`

	ActiveValidators []Validator `json:"activeValidators,omitempty"`

	PendingActiveValidatorsID   string      `json:"pendingActiveValidatorsId,omitempty"`
	PendingActiveValidatorsSize json.Number `json:"pendingActiveValidatorsSize,omitempty"`
	// PendingRemovals: [],
	StakingPoolMappingsId   string      `json:"stakingPoolMappingsId,omitempty"`
	StakingPoolMappingsSize json.Number `json:"stakingPoolMappingsSize,omitempty"`
	InactivePoolsId         string      `json:"inactivePoolsId,omitempty"`
	InactivePoolsSize       json.Number `json:"inactivePoolsSize,omitempty"`
	ValidatorCandidatesId   string      `json:"validatorCandidatesId,omitempty"`
	ValidatorCandidatesSize json.Number `json:"validatorCandidatesSize,omitempty"`
	// AtRiskValidators: [],
	// ValidatorReportRecords: []
}

type Validator struct {
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

	VotingPower    json.Number `json:"votingPower,omitempty"`
	OperationCapID string      `json:"operationCapId,omitempty"`
	GasPrice       Mist        `json:"gasPrice,omitempty"`
	CommissionRate json.Number `json:"commissionRate,omitempty"`

	NextEpochStake          Mist        `json:"nextEpochStake,omitempty"`
	NextEpochGasPrice       Mist        `json:"nextEpochGasPrice,omitempty"`
	NextEpochCommissionRate json.Number `json:"nextEpochCommissionRate,omitempty"`

	StakingPoolID                string `json:"stakingPoolId,omitempty"`
	StakingPoolActivationEpoch   Epoch  `json:"stakingPoolActivationEpoch,omitempty"`
	StakingPoolDeactivationEpoch Epoch  `json:"stakingPoolDeactivationEpoch,omitempty"`
	StakingPoolSuiBalance        Mist   `json:"stakingPoolSuiBalance,omitempty"`

	RewardsPool      Mist `json:"rewardsPool,omitempty"`
	PoolTokenBalance Mist `json:"poolTokenBalance,omitempty"`

	PendingStake             Mist `json:"pendingStake,omitempty"`
	PendingTotalSuiWithdraw  Mist `json:"pendingTotalSuiWithdraw,omitempty"`
	PendingPoolTokenWithdraw Mist `json:"pendingPoolTokenWithdraw,omitempty"`

	ExchangeRatesID   string      `json:"exchangeRatesId,omitempty"`
	ExchangeRatesSize json.Number `json:"exchangeRatesSize,omitempty"`
}
