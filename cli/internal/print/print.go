package print

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/cli/internal/db/model"
	"github.com/goat-systems/tzpay/v2/cli/internal/enviroment"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// Table prints a payout in table format
func Table(ctx context.Context, payouts *model.Payout) {
	base := enviroment.GetEnviromentFromContext(ctx)
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cylce", "Baker", "Wallet", "Rewards", "Operation"})
	table.Append([]string{
		strconv.Itoa(payouts.Cycle),
		base.Delegate,
		base.Wallet.Address,
		fmt.Sprintf("%.6f", float64(payouts.FrozenBalance.Int64())/float64(gotezos.MUTEZ)),
		groomOperations(payouts.OperationsLink...),
	})

	table.Render()

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Delegation", "Share", "Gross", "Net", "Fee"})

	var net, fee float64
	sort.Sort(payouts.DelegationEarnings)
	for _, delegation := range payouts.DelegationEarnings {
		table.Append([]string{
			delegation.Address,
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
func JSON(payouts *model.Payout) error {
	prettyJSON, err := json.MarshalIndent(payouts, "", "    ")
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
		op = strings.TrimSuffix(strings.TrimPrefix(op, "\""), "\"")

		if i == 0 {
			operation = op
		} else {
			operation = fmt.Sprintf(", %s", op)
		}
	}

	return operation
}
