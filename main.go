package main

import (
	"github.com/goat-systems/tzpay/v3/internal/cmd"
	"github.com/spf13/cobra"
)

func main() {
	rootCommand := &cobra.Command{
		Use:   "tzpay",
		Short: "A bulk payout tool for bakers in the Tezos Ecosystem",
	}
	rootCommand.AddCommand(
		cmd.DryRunCommand(),
		cmd.ServCommand(),
		cmd.RunCommand(),
		cmd.NewVersionCommand(),
		cmd.NewSetupCommand(),
	)

	rootCommand.Execute()
}
