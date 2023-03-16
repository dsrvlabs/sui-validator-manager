package main

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/dsrvlabs/sui-validator-manager/types"
)

type ValidatorCollector struct {
	systemState *types.SuiSystemState
}

func (v *ValidatorCollector) Refresh() {
	// TODO:
	cli := rpc.NewClient([]string{"https://wave3-rpc.testnet.sui.io:443"})

	var err error
	v.systemState, err = cli.LatestSuiSystemState()

	_ = err
}

func (v *ValidatorCollector) Render() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"#", "Name", "Vote(%)",
		"Stake(SUI)", "Next Stake(SUI)",
		"Reward Pool(SUI)", "RGP(MIST)",
	})

	p := message.NewPrinter(language.English)

	validators := v.systemState.ActiveValidators
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
}
