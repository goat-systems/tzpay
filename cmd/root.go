package cmd

import (
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "payman",
		Short: "A bulk payout tool for bakers in the Tezos Ecosystem",
	}

	rootCommand.AddCommand(
		newPayoutCommand(),
	)

	return rootCommand
}

func Execute() error {
	return newRootCommand().Execute()
}
