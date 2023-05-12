package main

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/dsrvlabs/sui-validator-manager/config"
	"github.com/dsrvlabs/sui-validator-manager/pricefeeder"
	"github.com/dsrvlabs/sui-validator-manager/rgp"
	"github.com/dsrvlabs/sui-validator-manager/rpc"
	"github.com/dsrvlabs/sui-validator-manager/types"
	"github.com/dsrvlabs/sui-validator-manager/view"
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

	_ = console.PersistentFlags().StringP("config", "c", "", "Config file path")
	viper.BindPFlag("console_config", console.PersistentFlags().Lookup("config"))
	rootCmd.AddCommand(console)

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

	_ = list.PersistentFlags().StringP("config", "c", "", "Config file path")
	viper.BindPFlag("list_config", list.PersistentFlags().Lookup("config"))
	rootCmd.AddCommand(list)

	predict := &cobra.Command{
		Use:     "predict",
		Aliases: []string{"p"},
		RunE: func(cmd *cobra.Command, args []string) error {
			configFile := viper.GetString("predict_config")
			config, err := config.Load(configFile)
			if err != nil {
				return err
			}

			valAddress := viper.GetString("predict_val_addr")

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

			// Filter by name and address
			var validator *types.Validator
			for _, val := range state.ActiveValidators {
				if val.SuiAddress == valAddress {
					validator = &val
					break
				}
			}

			if validator == nil {
				return errors.New("cannot find validator")
			}

			// Get self stakes
			stakes, err := cli.GetStakes(valAddress)
			if err != nil {
				return err
			}

			totalSelfStakes := stakes.StakeSum().BigInt()
			storageFund := new(big.Int).Sub(
				state.StorageFundNonRefundableBalance.BigInt(),
				state.StorageFundTotalObjectStorageRebate.BigInt())

			// TODO: Predict Computational cost, not use current.
			reward, err := rgp.CalculateReward(
				*state,
				storageFund,
				cp.EpochRollingGasCostSummary.ComputationCost.BigInt(),
				totalSelfStakes,
				valAddress,
			)
			if err != nil {
				return err
			}

			suiReward := types.Mist{}
			suiReward.SetString(reward.String())

			// TODO: Get token price.
			feeder := pricefeeder.NewClient()
			suiUSD, err := feeder.QueryPrice()
			if err != nil {
				return err
			}

			bfSuiInUSD := new(big.Float).SetFloat64(suiUSD)
			rewardUSD := new(big.Float).Mul(suiReward.Sui(), bfSuiInUSD)

			fmt.Printf("Expected reward of %s\n", valAddress)
			fmt.Printf("* %.2f MIST\n", suiReward.Sui())
			fmt.Printf("* $%s\n", rewardUSD.String())

			// TODO: Get daily operational cost.

			// TODO: Calculate profit.
			// TODO: Extract Reference Gas Price.

			return nil
		},
	}

	_ = predict.PersistentFlags().StringP("config", "c", "", "Config file path")
	viper.BindPFlag("predict_config", predict.PersistentFlags().Lookup("config"))

	_ = predict.PersistentFlags().StringP("address", "a", "", "Validator address")
	viper.BindPFlag("predict_val_addr", predict.PersistentFlags().Lookup("address"))
	rootCmd.AddCommand(predict)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
