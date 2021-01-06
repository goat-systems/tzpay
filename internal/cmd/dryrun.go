package cmd

import (
	"strconv"

	"github.com/goat-systems/tzpay/v3/internal/config"
	"github.com/goat-systems/tzpay/v3/internal/payout"
	"github.com/goat-systems/tzpay/v3/internal/print"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// DryRun -
type DryRun struct {
	payout payout.IFace
	config config.Config
	cycle  int
	table  bool
}

// NewDryRun returns a new dryrun
func NewDryRun(cycle string, table bool) DryRun {
	config, err := config.New()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to load config.")
	}

	// Clear sensitive data if loaded
	config.Key.Password = ""
	config.Key.Esk = ""

	c, err := strconv.Atoi(cycle)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to parse cycle argument into integer.")
	}

	payout, err := payout.New(config, c, false, false)
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to intialize payout.")
	}

	return DryRun{
		payout: payout,
		config: config,
		cycle:  c,
		table:  table,
	}
}

// DryRunCommand returns the cobra command for dryrun
func DryRunCommand() *cobra.Command {
	var table bool

	var dryrun = &cobra.Command{
		Use:     "dryrun",
		Short:   "dryrun simulates a payout",
		Long:    "dryrun simulates a payout and prints the result in json or a table",
		Example: `tzpay dryrun <cycle>`,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				log.Fatal("Missing cycle as argument.")
			}

			dryrun := NewDryRun(args[0], table)
			dryrun.execute()
		},
	}
	dryrun.PersistentFlags().BoolVarP(&table, "table", "t", false, "formats result into a table (Default: json)")

	return dryrun
}

func (d *DryRun) execute() {
	rewardsSplit, err := d.payout.Execute()
	if err != nil {
		log.WithField("error", err.Error()).Fatal("Failed to execute payout.")
	}

	if d.table {
		print.Table(d.cycle, d.config.Baker.Address, rewardsSplit)
	} else {
		err := print.JSON(rewardsSplit)
		if err != nil {
			log.WithField("error", err.Error()).Fatal("Failed to print JSON report.")
		}
	}
}
