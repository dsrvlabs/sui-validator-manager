package view

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type NetworkViewData struct {
	EpochNo           string
	CheckpointNo      string
	TXCount           string
	ComputationCost   string
	StorageCost       string
	StorageRebate     string
	ReferenceGasPrice string
}

type NetworkView struct {
	data *NetworkViewData
}

func (c *NetworkView) Render() {
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
		"RGP(MIST)",
	})

	t.AppendRow(table.Row{
		c.data.EpochNo,
		c.data.CheckpointNo,
		c.data.TXCount,
		c.data.ComputationCost,
		c.data.StorageCost,
		c.data.StorageRebate,
		c.data.ReferenceGasPrice,
	})

	t.Render()
}

func NewNetworkView(data *NetworkViewData) Renderer {
	return &NetworkView{data: data}
}
