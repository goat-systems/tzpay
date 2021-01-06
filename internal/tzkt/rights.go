package tzkt

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
)

type Rights []struct {
	Type      string    `json:"type"`
	Cycle     int       `json:"cycle"`
	Level     int       `json:"level"`
	Timestamp time.Time `json:"timestamp"`
	Priority  int       `json:"priority"`
	Slots     int       `json:"slots"`
	Baker     struct {
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"baker"`
	Status string `json:"status"`
}

/*
GetRights -
See: https://api.tzkt.io/#operation/Rights_Get
*/
func (t *Tzkt) GetRights(options ...URLParameters) (Rights, error) {
	resp, err := t.get("/v1/rights", options...)
	if err != nil {
		return Rights{}, errors.Wrapf(err, "failed to get reward split")
	}

	var rights Rights
	if err := json.Unmarshal(resp, &rights); err != nil {
		return Rights{}, errors.Wrap(err, "failed to get reward split")
	}

	return rights, nil
}
