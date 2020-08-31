package print

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/internal/payout"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Table prints a payout in table format
func Table(delegate, walletAddress string, report payout.Report) {
	if walletAddress == "" {
		walletAddress = "N/A"
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cylce", "Baker", "Wallet", "Rewards", "Operation"})
	table.Append([]string{
		strconv.Itoa(report.Cycle),
		delegate,
		walletAddress,
		fmt.Sprintf("%.6f", float64(report.FrozenBalance.Int64())/float64(gotezos.MUTEZ)),
		groomOperations(report.OperationsLink...),
	})

	table.Render()

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Delegation", "Share", "Gross", "Net", "Fee"})

	var net, fee float64
	sort.Sort(report.DelegationEarnings)
	for _, delegation := range report.DelegationEarnings {
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
func JSON(report payout.Report) error {
	prettyJSON, err := json.Marshal(report)
	if err != nil {
		return errors.Wrap(err, "failed to parse report into json")
	}

	log.WithField("payout", string(prettyJSON)).Info("Payout for cycle complete.")
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
