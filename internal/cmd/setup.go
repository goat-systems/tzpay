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
			sb.WriteString("TZPAY_BAKER=<TODO (e.g. tz1SUgyRB8T5jXgXAwS33pgRHAKrafyg87Yc)>\n")
			sb.WriteString("TZPAY_BAKER_FEE=<TODO (e.g. 0.05 for 5%)>\n")
			sb.WriteString("TZPAY_WALLET_ESK=<TODO (e.g. edesk1fddn27MaLcQVEdZpAYiyGQNm6UjtWiBfNP2ZenTy3CFsoSVJgeHM9pP9cvLJ2r5Xp2quQ5mYexW1LRKee2)>\n")
			sb.WriteString("TZPAY_WALLET_PASSWORD=<TODO (e.g. password12345##)>\n")
			sb.WriteString("###### OPTIONAL ENVIROMENT VARIABLES ######\n")
			sb.WriteString("TZPAY_BAKER_MINIMUM_PAYMENT=<TODO (e.g. MUTEZ 10000)>\n")
			sb.WriteString("TZPAY_BAKER_EARNINGS_ONLY=<TODO (e.g. True)>\n")
			sb.WriteString("TZPAY_BAKER_BLACK_LIST=<TODO (e.g. KT19Aro5JcjKH7J7RA6sCRihPiBQzQED3oQC, KT1CQiyDJ3mMVDoEqLY8Fz1onFXo5ycp5BDN)>\n")
			sb.WriteString("TZPAY_BAKER_LIQUIDITY_CONTRACTS=<TODO (e.g. KT19Aro5JcjKH7J7RA6sCRihPiBQzQED3oQC, KT1CQiyDJ3mMVDoEqLY8Fz1onFXo5ycp5BDN)>\n")
			sb.WriteString("TZPAY_API_TZKT=<TODO (e.g. https://api.tzkt.io )>\n")
			sb.WriteString("TZPAY_API_TEZOS=<TODO (e.g. https://tezos.giganode.io/)>\n")
			sb.WriteString("TZPAY_OPERATIONS_NETWORK_FEE=<TODO (e.g. 2941)>\n")
			sb.WriteString("TZPAY_OPERATIONS_GAS_LIMIT=<TODO (e.g. 26283)>\n")
			sb.WriteString("TZPAY_OPERATIONS_BATCH_SIZE=<TODO (e.g. 125)>\n")
			fmt.Println(sb.String())
		},
	}
	return setup
}
