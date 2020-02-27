package cli

import (
	"github.com/goat-systems/tzpay/v2/cli/internal/cmd"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:   "tzpay",
		Short: "A bulk payout tool for bakers in the Tezos Ecosystem",
	}

	rootCommand.AddCommand(
		cmd.NewDryRunCommand(),
		cmd.NewVersionCommand(),
	)

	return rootCommand
}

// Execute executes the user command.
func Execute() error {
	return newRootCommand().Execute()
}
