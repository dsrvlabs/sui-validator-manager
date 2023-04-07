package main

import (
	"time"

	"github.com/spf13/cobra"

	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/dsrvlabs/sui-validator-manager/view"
)

var (
	rpcURL = "https://fullnode.testnet.sui.io:443"
)

func main() {
	rootCmd := &cobra.Command{}
	console := &cobra.Command{
		Use:     "monitor",
		Aliases: []string{"m"},
		Run: func(cmd *cobra.Command, args []string) {
			for {
				cli := rpc.NewClient([]string{rpcURL})

				seq, err := cli.LatestCheckpointSequenceNumber()
				if err != nil {
					time.Sleep(time.Second * 1)
					continue
				}

				cp, err := cli.Checkpoint(seq)
				if err != nil {
					time.Sleep(time.Second * 1)
					continue
				}

				state, err := cli.LatestSuiSystemState()
				if err != nil {
					time.Sleep(time.Second * 1)
					continue
				}

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

				nView.Render()
				vView.Render()

				time.Sleep(time.Second * 5)
			}
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			cli := rpc.NewClient([]string{rpcURL})

			seq, err := cli.LatestCheckpointSequenceNumber()
			if err != nil {
				return
			}

			cp, err := cli.Checkpoint(seq)
			if err != nil {
				return
			}

			state, err := cli.LatestSuiSystemState()
			if err != nil {
				return
			}

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
		},
	}

	rootCmd.AddCommand(console)
	rootCmd.AddCommand(list)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
