package print

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v3/internal/tzkt"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Table prints a payout in table format
func Table(cycle int, delegate string, rewards tzkt.RewardsSplit) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Cylce", "Baker", "Share", "Rewards", "Fees", "Total", "Operations"})
	table.Append([]string{
		strconv.Itoa(cycle),
		delegate,
		fmt.Sprintf("%.6f", rewards.BakerShare),
		fmt.Sprintf("%.6f", float64(rewards.BakerRewards)/float64(gotezos.MUTEZ)),
		fmt.Sprintf("%.6f", float64(rewards.BakerCollectedFees)/float64(gotezos.MUTEZ)),
		fmt.Sprintf("%.6f", float64(rewards.BakerRewards+rewards.BakerCollectedFees)/float64(gotezos.MUTEZ)),
		groomOperations(rewards.OperationLink...),
	})

	table.Render()

	table = tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Delegation", "Share", "Gross", "Net", "Fee"})

	liquidityProviderTable := tablewriter.NewWriter(os.Stdout)
	liquidityProviderTable.SetHeader([]string{"Liquidiy Provider", "Contract", "Share", "Gross", "Net", "Fee"})

	var net, fee, liquidityNet, liquidityFee, gross, share float64

	for _, delegation := range rewards.Delegators {
		for _, lp := range delegation.LiquidityProviders {
			liquidityProviderTable.Append([]string{
				lp.Address,
				delegation.Address,
				fmt.Sprintf("%.6f", lp.Share),
				fmt.Sprintf("%.6f", float64(lp.GrossRewards)/float64(gotezos.MUTEZ)),
				fmt.Sprintf("%.6f", float64(lp.NetRewards)/float64(gotezos.MUTEZ)),
				fmt.Sprintf("%.6f", float64(lp.Fee)/float64(gotezos.MUTEZ)),
			})

			liquidityNet += float64(lp.NetRewards) / float64(gotezos.MUTEZ)
			liquidityFee += float64(lp.Fee) / float64(gotezos.MUTEZ)
			gross += float64(lp.GrossRewards) / float64(gotezos.MUTEZ)
			share += lp.Share
		}

		table.Append([]string{
			delegation.Address,
			fmt.Sprintf("%.6f", delegation.Share),
			fmt.Sprintf("%.6f", float64(delegation.GrossRewards)/float64(gotezos.MUTEZ)),
			fmt.Sprintf("%.6f", float64(delegation.NetRewards)/float64(gotezos.MUTEZ)),
			fmt.Sprintf("%.6f", float64(delegation.Fee)/float64(gotezos.MUTEZ)),
		})
		net += float64(delegation.NetRewards) / float64(gotezos.MUTEZ)
		fee += float64(delegation.Fee) / float64(gotezos.MUTEZ)
	}

	table.SetFooter([]string{"", "", "TOTAL", fmt.Sprintf("%.6f", net), fmt.Sprintf("%.6f", fee)}) // Add Footer
	liquidityProviderTable.SetFooter([]string{"", "TOTAL", fmt.Sprintf("%.6f", share), fmt.Sprintf("%.6f", gross), fmt.Sprintf("%.6f", liquidityNet), fmt.Sprintf("%.6f", liquidityFee)})

	table.Render()

	if liquidityProviderTable.NumLines() > 0 {
		liquidityProviderTable.Render()
	}
}

// JSON prints a payout to json
func JSON(rewards tzkt.RewardsSplit) error {
	prettyJSON, err := json.Marshal(rewards)
	if err != nil {
		return errors.Wrap(err, "failed to parse reward split into json")
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
