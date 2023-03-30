package main

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/dsrvlabs/sui-validator-manager/view"
)

func main() {
	rootCmd := &cobra.Command{}
	console := &cobra.Command{
		Use:     "monitor",
		Aliases: []string{"m"},
		Run: func(cmd *cobra.Command, args []string) {
			for {
				cli := rpc.NewClient([]string{"https://fullnode.testnet.sui.io:443"})

				seq, err := cli.LatestCheckpointSequenceNumber()
				cp, err := cli.Checkpoint(seq)

				state, err := cli.LatestSuiSystemState()

				nView := view.NewNetworkView(&view.NetworkViewData{
					EpochNo:           cp.Epoch.String(),
					CheckpointNo:      cp.SequenceNumber.String(),
					TXCount:           cp.NetworkTotalTransactions.String(),
					ComputationCost:   cp.EpochRollingGasCostSummary.ComputationCost.String(),
					StorageCost:       cp.EpochRollingGasCostSummary.StorageCost.String(),
					StorageRebate:     cp.EpochRollingGasCostSummary.StorageRebate.String(),
					ReferenceGasPrice: state.ReferenceGasPrice.String(),
				})

				vView := view.NewValidatorView(state)
				_ = err

				nView.Render()
				vView.Render()
				time.Sleep(5 * time.Second)
			}
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			cli := rpc.NewClient([]string{"https://fullnode.testnet.sui.io:443"})

			seq, err := cli.LatestCheckpointSequenceNumber()
			cp, err := cli.Checkpoint(seq)

			state, err := cli.LatestSuiSystemState()

			nView := view.NewNetworkView(&view.NetworkViewData{
				EpochNo:           cp.Epoch.String(),
				CheckpointNo:      cp.SequenceNumber.String(),
				TXCount:           cp.NetworkTotalTransactions.String(),
				ComputationCost:   cp.EpochRollingGasCostSummary.ComputationCost.String(),
				StorageCost:       cp.EpochRollingGasCostSummary.StorageCost.String(),
				StorageRebate:     cp.EpochRollingGasCostSummary.StorageRebate.String(),
				ReferenceGasPrice: state.ReferenceGasPrice.String(),
			})
			nView.Render()

			vView := view.NewValidatorView(state)
			vView.Render()

			_ = err
		},
	}

	rootCmd.AddCommand(console)
	rootCmd.AddCommand(list)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
