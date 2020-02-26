package print

import (
	"encoding/json"
	"fmt"

	"github.com/goat-systems/payman/v2/cmd/internal/delegates"
	"github.com/pkg/errors"
)

// Printer is an interface to tzpay print functions.
type Printer interface {
	PrintJSON()
	PrintTable()
}

// ReportPrinter is a printer for the report command.
type ReportPrinter struct{}

// NewReportPrinter returns a pointer to a new ReportPrinter
func NewReportPrinter() *ReportPrinter {
	return &ReportPrinter{}
}

// PrintJSON prints an array of DelegationEarnings in json format.
func (r *ReportPrinter) PrintJSON(delegationEarnings *[]delegates.DelegationEarnings) error {
	out, err := json.MarshalIndent(delegationEarnings, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to print json")
	}
	fmt.Println(out)

	return nil
}

// PrintTable prints an array of DelegationEarnings in table format.
// func (r *ReportPrinter) PrintTable(cycle int, rewards string, delegationEarnings *[]delegates.DelegationEarnings) error {
// 	table := tablewriter.NewWriter(r.general.Writer())
// 	table.SetHeader([]string{"Delegate", "Cycle", "Rewards", "Fees"})
// }

// func getTotalFees(delegationEarnings *[]delegates.DelegationEarnings) string {
// 	var fees string
// 	bigIntFees := big.NewInt(0)
// 	for _, d := range *delegationEarnings {
// 		bigIntFees = d.Fee.Add(d.Fee., bigIntFees)
// 	}
// }
