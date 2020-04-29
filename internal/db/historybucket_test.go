package db

import (
	"math/big"
	"os"
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/goat-systems/tzpay/v2/internal/db/model"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_SavePayout(t *testing.T) {
	type want struct {
		err bool
	}

	cases := []struct {
		name  string
		input model.Payout
		want  want
	}{
		{
			"is successful",
			model.Payout{
				StakingBalance: big.NewInt(0),
				CycleHash:      "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
			},
			want{
				false,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(nil, "./tzpay.db")
			assert.Nil(t, err)
			assert.NotNil(t, db)

			err = db.SavePayout(tt.input)
			checkErr(t, tt.want.err, "", err)

			err = os.Remove("./tzpay.db")
			assert.Nil(t, err)
		})
	}
}

func Test_GetPayout(t *testing.T) {
	type input struct {
		cycle  int
		payout model.Payout
	}

	type want struct {
		err         bool
		errContains string
		payout      *model.Payout
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				10,
				model.Payout{
					StakingBalance: big.NewInt(0),
					CycleHash:      "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
				},
			},
			want{
				false,
				"",
				&model.Payout{
					StakingBalance: big.NewInt(0),
					CycleHash:      "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
				},
			},
		},
		{
			"handles cycle error",
			input{
				10,
				model.Payout{
					StakingBalance: big.NewInt(0),
					CycleHash:      "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
				},
			},
			want{
				true,
				"failed to get payout from history bucket: failed to get cycle",
				nil,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			db, err := New(nil, "./tzpay.db")
			db.gt = &gotezosMock{cycleErr: tt.want.err}
			assert.Nil(t, err)
			assert.NotNil(t, db)

			err = db.SavePayout(tt.input.payout)
			checkErr(t, false, "", err)

			payout, err := db.GetPayout(tt.input.cycle)
			checkErr(t, tt.want.err, tt.want.errContains, err)

			assert.Equal(t, tt.want.payout, payout)

			err = os.Remove("./tzpay.db")
			assert.Nil(t, err)
		})
	}
}

type gotezosMock struct {
	gotezos.IFace
	cycleErr bool
}

func (g *gotezosMock) Cycle(cycle int) (*gotezos.Cycle, error) {
	if g.cycleErr {
		return &gotezos.Cycle{}, errors.New("failed to get cycle")
	}
	return &gotezos.Cycle{
		RandomSeed:   "some_seed",
		RollSnapshot: 10,
		BlockHash:    "BLfEWKVudXH15N8nwHZehyLNjRuNLoJavJDjSZ7nq8ggfzbZ18p",
	}, nil
}
