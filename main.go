package main

import (
	"time"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{}
	console := &cobra.Command{
		Use:     "monitor",
		Aliases: []string{"m"},
		Run: func(cmd *cobra.Command, args []string) {
			n := NetworkCollector{}
			c := ValidatorCollector{}
			for {
				n.Refresh()
				n.Render()

				c.Refresh()
				c.Render()

				time.Sleep(5 * time.Second)
			}
		},
	}

	list := &cobra.Command{
		Use:     "list",
		Aliases: []string{"l"},
		Run: func(cmd *cobra.Command, args []string) {
			n := NetworkCollector{}
			c := ValidatorCollector{}

			n.Refresh()
			n.Render()

			c.Refresh()
			c.Render()
		},
	}

	rootCmd.AddCommand(console)
	rootCmd.AddCommand(list)

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
