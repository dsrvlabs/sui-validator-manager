package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/dsrvlabs/sui-validator-manager/view"

	"github.com/dsrvlabs/sui-validator-manager/config"
)

var rpcURL = "https://sui-rpc-mainnet.testnet-pride.com:443"

func main() {
	rootCmd := &cobra.Command{}
	console := &cobra.Command{
		Use:     "monitor",
		Aliases: []string{"m"},
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := viper.GetString("console_config")
			config, err := config.Load(configFile)
			if err != nil {
				return err
			}

			for {
				cli := rpc.NewClient([]string{config.RPC[0].Endpoint})

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
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := viper.GetString("list_config")
			config, err := config.Load(configFile)
			if err != nil {
				return err
			}

			cli := rpc.NewClient([]string{config.RPC[0].Endpoint})

			seq, err := cli.LatestCheckpointSequenceNumber()
			if err != nil {
				return err
			}

			cp, err := cli.Checkpoint(seq)
			if err != nil {
				return err
			}

			state, err := cli.LatestSuiSystemState()
			if err != nil {
				return err
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

			return nil
		},
	}

	_ = console.PersistentFlags().StringP("config", "c", "", "Config file path")
	viper.BindPFlag("console_config", console.PersistentFlags().Lookup("config"))
	rootCmd.AddCommand(console)

	_ = list.PersistentFlags().StringP("config", "c", "", "Config file path")
	viper.BindPFlag("list_config", list.PersistentFlags().Lookup("config"))
	rootCmd.AddCommand(list)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
