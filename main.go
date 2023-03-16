package main

import (
	"fmt"
	"math/big"
	"os"
	"sync"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/dsrvlabs/sui-validator-manager/types"
)

func main() {
	// TODO:
	cli := rpc.NewClient([]string{"https://wave3-rpc.testnet.sui.io:443"})

	var (
		seq   *big.Int
		cp    *types.Checkpoint
		state *types.SuiSystemState
		err   error
	)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		seq, err = cli.LatestCheckpointSequenceNumber()
		cp, err = cli.Checkpoint(seq)
	}()

	go func() {
		defer wg.Done()

		state, err = cli.LatestSuiSystemState()
	}()

	wg.Wait()

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	rowCfgAutoMerge := table.RowConfig{AutoMerge: true}
	t.AppendHeader(
		table.Row{
			"Network Information",
			"Network Information",
			"Network Information",
			"Network Information",
			"Network Information",
			"Network Information",
		},
		rowCfgAutoMerge,
	)
	t.AppendHeader(table.Row{
		"Epoch", "Checkpoint",
		"# TXs", "C. Cost", "S. Cost", "S. Rebate"})
	t.AppendRow(table.Row{
		cp.Epoch.String(), cp.SequenceNumber.String(),
		cp.NetworkTotalTransactions.String(),
		cp.EpochRollingGasCostSummary.ComputationCost.String(),
		cp.EpochRollingGasCostSummary.StorageCost.String(),
		cp.EpochRollingGasCostSummary.StorageRebate.String(),
	})

	t.Render()

	t = table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"#", "Name", "Vote(%)",
		"Stake(SUI)", "Next Stake(SUI)",
		"Reward Pool(SUI)", "RGP(MIST)",
	})

	p := message.NewPrinter(language.English)

	validators := state.ActiveValidators
	for i, v := range validators {
		stake, _ := v.StakingPoolSuiBalance.Sui().Float64()
		nextStake, _ := v.NextEpochStake.Sui().Float64()
		rewards, _ := v.RewardsPool.Sui().Float64()

		t.AppendRow(table.Row{
			i + 1,
			v.Name,
			float32(v.VotingPower) / 100.0,
			p.Sprintf("%f", stake),
			p.Sprintf("%f", nextStake),
			p.Sprintf("%f", rewards),
			v.GasPrice.String(),
		})
	}

	t.Render()

	_ = err
}
