package delegates

import (
	"errors"
	"math/big"
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v2"
	"github.com/stretchr/testify/assert"
)

func Test_processDelegation(t *testing.T) {
	type want struct {
		err                bool
		errContains        string
		delegationEarnings *DelegationEarnings
	}

	cases := []struct {
		name  string
		input *processDelegationInput
		want  want
	}{
		{
			"is successful",
			&processDelegationInput{
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				false,
				"",
				&DelegationEarnings{Fee: big.NewInt(4000000), GrossRewards: big.NewInt(80000000), NetRewards: big.NewInt(76000000), Share: 0.1},
			},
		},
		{
			"handles failure",
			&processDelegationInput{
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				true,
				"failed to get balance",
				nil,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{
				fee: 0.05,
				gt: &gotezosMock{
					err: tt.want.err,
				},
			}
			out, err := baker.processDelegation(tt.input)
			checkErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.delegationEarnings, out)
		})
	}
}

func Test_processDelegations(t *testing.T) {
	type want struct {
		err        bool
		errcount   int
		successful int
	}

	cases := []struct {
		name  string
		input *processDelegationsInput
		want  want
	}{
		{
			"is successful",
			&processDelegationsInput{
				delegations: &[]string{
					"some_delegation",
					"some_delegation1",
					"some_delegation2",
				},
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				false,
				3,
				0,
			},
		},
		{
			"handles failure",
			&processDelegationsInput{
				delegations: &[]string{
					"some_delegation",
					"some_delegation1",
					"some_delegation2",
				},
				stakingBalance: big.NewInt(100000000000),
				frozenBalanceRewards: &gotezos.FrozenBalance{
					Rewards: gotezos.Int{Big: big.NewInt(800000000)},
				},
			},
			want{
				true,
				0,
				3,
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			baker := &Baker{
				fee: 0.05,
				gt: &gotezosMock{
					err: tt.want.err,
				},
			}
			out := baker.proccessDelegations(tt.input)
			successful := 0
			errcount := 0
			returns := 0

			for o := range out {
				if returns != len(*tt.input.delegations) {
					if o.err != nil {
						errcount++
					} else {
						successful++
					}
				} else {
					break
				}
				returns++
			}

			assert.Equal(t, tt.want.successful, successful)
			assert.Equal(t, tt.want.errcount, errcount)
		})
	}
}

func Test_ParseBlackList(t *testing.T) {
	cases := []struct {
		name  string
		input string
		want  []string
	}{
		{
			"is successful",
			"some_address, some_other_address, yet_another_address",
			[]string{
				"some_address",
				"some_other_address",
				"yet_another_address",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			out := ParseBlackList(tt.input)
			assert.Equal(t, tt.want, out)
		})
	}

}

func checkErr(t *testing.T, wantErr bool, errContains string, err error) {
	if wantErr {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errContains)
	} else {
		assert.Nil(t, err)
	}
}

type gotezosMock struct {
	gotezos.IFace
	err bool
}

func (g *gotezosMock) Balance(blockhash, address string) (*big.Int, error) {
	if g.err {
		return big.NewInt(0), errors.New("failed to get balance")
	}
	return big.NewInt(10000000000), nil
}
