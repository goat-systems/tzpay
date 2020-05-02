package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// NewSetupCommand returns a new setup cobra command
func NewSetupCommand() *cobra.Command {
	var setup = &cobra.Command{
		Use:     "setup",
		Short:   "setup prints a list of enviroment variables needed to get started.",
		Example: `tzpay setup`,
		Run: func(cmd *cobra.Command, args []string) {
			var sb strings.Builder
			sb.WriteString("###### REQUIRED ENVIROMENT VARIABLES ######\n")
			sb.WriteString("TZPAY_HOST_NODE=<TODO (e.g. http://127.0.0.1:8732)>\n")
			sb.WriteString("TZPAY_BAKERS_FEE=<TODO (e.g. 0.05 for 5%)>\n")
			sb.WriteString("TZPAY_DELEGATE=<TODO (e.g. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc)>\n")
			sb.WriteString("TZPAY_WALLET_SECRET=<TODO (e.g. edesk...)>\n")
			sb.WriteString("TZPAY_WALLET_PASSWORD=<TODO (e.g. password)>\n\n")
			sb.WriteString("###### OPTIONAL ENVIROMENT VARIABLES ######\n")
			sb.WriteString("TZPAY_BLACKLIST=<TODO (e.g. some_blacklist_file.json)>\n")
			sb.WriteString("TZPAY_NETWORK_GAS_LIMIT=<TODO (e.g. 30000)>\n")
			sb.WriteString("TZPAY_NETWORK_FEE=<TODO (e.g. 3000)>\n")
			sb.WriteString("TZPAY_MINIMUM_PAYMENT=<TODO (e.g. 3000)>\n")

			fmt.Println(sb.String())
		},
	}
	return setup
}
