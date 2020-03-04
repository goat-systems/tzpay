package print

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	gotezos "github.com/goat-systems/go-tezos/v2"
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
func Table(ctx context.Context, payouts *baker.Payout, operations ...string) {
	base := enviroment.GetEnviromentFromContext(ctx)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Baker", "Wallet", "Rewards", "Operation"})
	table.Append([]string{
		base.Delegate,
		base.Wallet.Address,
		fmt.Sprintf("%.6f", float64(payouts.FrozenBalance.Int64())/float64(gotezos.MUTEZ)),
		groomOperations(operations...),
	})

	table.Render()

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Delegation", "Share", "Gross", "Net", "Fee"})

	var net, fee float64
	sort.Sort(payouts.DelegationEarnings)
	for _, delegation := range payouts.DelegationEarnings {
		table.Append([]string{
			delegation.Delegation,
			fmt.Sprintf("%.6f", delegation.Share),
			fmt.Sprintf("%.6f", float64(delegation.GrossRewards.Int64())/float64(gotezos.MUTEZ)),
			fmt.Sprintf("%.6f", float64(delegation.NetRewards.Int64())/float64(gotezos.MUTEZ)),
			fmt.Sprintf("%.6f", float64(delegation.Fee.Int64())/float64(gotezos.MUTEZ)),
		})
		net += float64(delegation.NetRewards.Int64()) / float64(gotezos.MUTEZ)
		fee += float64(delegation.Fee.Int64()) / float64(gotezos.MUTEZ)
	}

	table.SetFooter([]string{"", "", "TOTAL", fmt.Sprintf("%.6f", net), fmt.Sprintf("%.6f", fee)}) // Add Footer

	table.Render()
}

// JSON prints a payout to json
func JSON(payouts *baker.Payout, operations ...string) error {
	prettyJSON, err := json.MarshalIndent(Payout{groomOperations(operations...), payouts}, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to print json")
	}

	fmt.Println(string(prettyJSON))

	return nil
}

func groomOperations(operations ...string) string {
	var operation string
	if operations == nil {
		operation = "N/A"
	}
	for i, op := range operations {
		if i == 0 {
			operation = op
		} else {
			operation = fmt.Sprintf(", %s", op)
		}
	}

	return operation
}
