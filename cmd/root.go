package cmd

import (
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "tzpay",
		Short: "A bulk payout tool for bakers in the Tezos Ecosystem",
	}

	rootCommand.AddCommand(
		newPayoutCommand(),
		newReportCommand(),
	)

	return rootCommand
}

// Execute executes the user command.
func Execute() error {
	return newRootCommand().Execute()
}
