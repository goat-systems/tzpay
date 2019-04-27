package reporting

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"encoding/csv"

	goTezos "github.com/DefinitelyNotAGoat/go-tezos"
	"github.com/olekukonko/tablewriter"
)

// Reporter is a structer that contains a general logger and a csv writer for payout reports
type Reporter struct {
	general *log.Logger
	report  *csv.Writer
}

// Log uses the genral logger and writes the message
func (r *Reporter) Log(msg interface{}) {
	r.general.Println(msg)
}

// NewReporter creates a new reporter for general logging and payout reports (csv)
func NewReporter(general *log.Logger) (Reporter, error) {
	r := Reporter{general: general}
	report, err := r.getCSVWriter()
	if err != nil {
		return r, err
	}
	r.report = report
	return r, nil
}

// PrintPaymentsTable takes in payments and prints them to a table for general logging
func (r *Reporter) PrintPaymentsTable(payments goTezos.DelegateReport) {
	total := []string{}
	data := r.formatData(payments)
	if len(data) > 0 {
		total = data[len(data)-1]
		data = data[:len(data)-1]
	}

	table := tablewriter.NewWriter(r.general.Writer())
	table.SetHeader([]string{"Address", "Share", "Gross", "Fee", "Net"})
	table.SetFooter(total)

	for _, v := range data {
		table.Append(v)
	}
	table.Render()
}

// formatData parses payments into a double array of data for table or csv printing
func (r *Reporter) formatData(payments goTezos.DelegateReport) [][]string {
	var data [][]string
	var totalNet float64
	var totalGross float64
	var totalFee float64
	for _, payment := range payments.Delegations {
		share := payment.Share * 100
		strShare := fmt.Sprintf("%.6f", share)
		fee, _ := strconv.Atoi(payment.Fee)
		floatFee := float64(fee) / float64(goTezos.MUTEZ)
		gross, _ := strconv.Atoi(payment.GrossRewards)
		floatGross := float64(gross) / float64(goTezos.MUTEZ)
		net, _ := strconv.Atoi(payment.NetRewards)
		floatNet := float64(net) / float64(goTezos.MUTEZ)

		totalNet = totalNet + floatNet
		totalGross = totalGross + floatGross
		totalFee = totalFee + floatFee
		data = append(data, []string{payment.DelegationPhk, strShare, fmt.Sprintf("%.6f", floatGross), fmt.Sprintf("%.6f", floatFee), fmt.Sprintf("%.6f", floatNet)})
	}
	data = append(data, []string{"", "Total", fmt.Sprintf("%.6f", totalGross), fmt.Sprintf("%.6f", totalFee), fmt.Sprintf("%.6f", totalNet)})
	return data
}

// WriteCSVReport writes payments to a csv file for reporting
func (r *Reporter) WriteCSVReport(payments goTezos.DelegateReport) {
	data := r.formatData(payments)
	if r.report != nil {
		for _, value := range data {
			r.report.Write(value)
		}
	}
	r.report.Flush()
}

// getCSVWriter opens a file with the current date for its name and
// returns a csv.Writer for that file
func (r *Reporter) getCSVWriter() (*csv.Writer, error) {
	fileName := r.buildFileName()
	f, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	report := csv.NewWriter(f)
	return report, nil
}

// buildFileName returns a string that represents a filename based
// of the current date
func (r *Reporter) buildFileName() string {
	return time.Now().Format(time.RFC3339) + ".csv"
}
