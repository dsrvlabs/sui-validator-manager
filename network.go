package main

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/dsrvlabs/sui-validator-manager/types"
)

type NetworkCollector struct {
	checkpoint *types.Checkpoint
}

func (c *NetworkCollector) Refresh() {
	cli := rpc.NewClient([]string{"https://wave3-rpc.testnet.sui.io:443"})

	seq, err := cli.LatestCheckpointSequenceNumber()
	cp, err := cli.Checkpoint(seq)

	c.checkpoint = cp

	_ = err
}

func (c *NetworkCollector) Render() {
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
		"# TXs", "C. Cost", "S. Cost", "S. Rebate",
	})

	cp := c.checkpoint
	t.AppendRow(table.Row{
		cp.Epoch.String(), cp.SequenceNumber.String(),
		cp.NetworkTotalTransactions.String(),
		cp.EpochRollingGasCostSummary.ComputationCost.String(),
		cp.EpochRollingGasCostSummary.StorageCost.String(),
		cp.EpochRollingGasCostSummary.StorageRebate.String(),
	})

	t.Render()
}
