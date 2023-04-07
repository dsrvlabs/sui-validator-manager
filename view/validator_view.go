package view

import (
	"os"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/dsrvlabs/sui-validator-manager/types"
)

type ByNextRGP []types.Validator

func (v ByNextRGP) Len() int      { return len(v) }
func (v ByNextRGP) Swap(i, j int) { v[i], v[j] = v[j], v[i] }
func (v ByNextRGP) Less(i, j int) bool {
	return v[i].NextEpochGasPrice.Sui().Cmp(v[j].NextEpochGasPrice.Sui()) == -1
}

type ValidatorView struct {
	systemState *types.SuiSystemState
}

func (v *ValidatorView) Render() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{
		"#", "Name", "Vote(%)", "Acc Votes(%)",
		"Stake(SUI)", "Next Stake(SUI)",
		"Reward Pool(SUI)", "RGP(MIST)",
		"Next RGP(MIST)",
	})

	p := message.NewPrinter(language.English)

	var (
		validators ByNextRGP
		nextRGP    types.Mist
	)

	validators = v.systemState.ActiveValidators
	sort.Sort(validators)

	accVotes := uint32(0)
	isCut := false

	for i, v := range validators {
		stake, _ := v.StakingPoolSuiBalance.Sui().Float64()
		nextStake, _ := v.NextEpochStake.Sui().Float64()
		rewards, _ := v.RewardsPool.Sui().Float64()

		accVotes += v.VotingPower

		t.AppendRow(table.Row{
			i + 1,
			v.Name,
			float32(v.VotingPower) / 100.0,
			float32(accVotes) / 100.0,
			p.Sprintf("%f", stake),
			p.Sprintf("%f", nextStake),
			p.Sprintf("%f", rewards),
			v.GasPrice.String(),
			v.NextEpochGasPrice.String(),
		})

		if float32(accVotes) > float32(10000.0*2.0/3.0) && !isCut {
			t.AppendSeparator()
			nextRGP = v.NextEpochGasPrice
			isCut = true
		}
	}

	t.AppendFooter(table.Row{"", "", "", "", "", "", "", "Next RGP", nextRGP.String()})

	t.Render()
}

func NewValidatorView(state *types.SuiSystemState) Renderer {
	return &ValidatorView{systemState: state}
}
