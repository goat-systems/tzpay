package print

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goat-systems/tzpay/v2/cli/internal/baker"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// Payout is the struct form of the JSON to print
type Payout struct {
	Operation string `json:"operation"`
	*baker.Payout
}

// Table prints a payout in table format
func Table(ctx context.Context, operation string, payouts *baker.Payout) {
	base := enviroment.GetEnviromentFromContext(ctx)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Baker", "Wallet", "Rewards", "Operation"})
	table.Append([]string{
		base.Delegate,
		base.Wallet.Address,
		payouts.FrozenBalance.String(),
		operation,
	})

	table.Render()

	table.SetHeader([]string{"Delegation", "Share", "Gross", "Net", "Fee"})
	for _, delegation := range payouts.DelegationEarnings {
		table.Append([]string{
			delegation.Delegation,
			fmt.Sprintf("%.6f", delegation.Share),
			fmt.Sprintf("%d", delegation.GrossRewards),
			fmt.Sprintf("%d", delegation.NetRewards),
			fmt.Sprintf("%d", delegation.Fee),
		})
	}

	table.Render()
}

// JSON prints a payout to json
func JSON(operation string, payouts *baker.Payout) error {
	prettyJSON, err := json.MarshalIndent(Payout{operation, payouts}, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to print json")
	}

	fmt.Println(string(prettyJSON))

	return nil
}
