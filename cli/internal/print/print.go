package print

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goat-systems/tzpay/v2/cli/internal/baker"
	"github.com/pkg/errors"
)

// Payout is the struct form of the JSON to print
type Payout struct {
	Operation string `json:"operation"`
	*baker.Payout
}

// JSON prints a payout to json
func JSON(ctx context.Context, operation string, payouts *baker.Payout) error {
	prettyJSON, err := json.MarshalIndent(Payout{operation, payouts}, "", "    ")
	if err != nil {
		return errors.Wrap(err, "failed to print json")
	}

	fmt.Println(prettyJSON)

	return nil
}
