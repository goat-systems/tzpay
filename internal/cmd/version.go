package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const (
	version = "v2.7.0-alpha"
	changed = "Support for dexter liquidity providers."
)

// NewVersionCommand returns a version cobra command
func NewVersionCommand() *cobra.Command {
	var version = &cobra.Command{
		Use:     "version",
		Short:   "version prints tzpay's version",
		Long:    "version prints tzpay's version to stdout",
		Example: `tzpay version`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s - %s\n", version, changed)
		},
	}
	return version
}
