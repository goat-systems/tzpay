package notifier

import (
	"testing"

	gotezos "github.com/goat-systems/go-tezos/v3"
	"github.com/goat-systems/tzpay/v2/internal/test"
	"github.com/stretchr/testify/assert"
)

func Test_isEndorsementSuccessful(t *testing.T) {
	m := MissedOpportunityNotifier{
		baker: "some_baker",
	}
	ok := m.isEndorsementSuccessful(&gotezos.Block{
		Operations: [][]gotezos.Operations{
			{
				{
					Contents: gotezos.Contents{
						Endorsements: []gotezos.Endorsement{
							{
								Metadata: &gotezos.EndorsementMetadata{
									Delegate: "some_baker",
								},
							},
						},
					},
				},
			},
		},
	})
	assert.Equal(t, true, ok)

	ok = m.isEndorsementSuccessful(&gotezos.Block{
		Operations: [][]gotezos.Operations{
			{
				{
					Contents: gotezos.Contents{
						Endorsements: []gotezos.Endorsement{
							{
								Metadata: &gotezos.EndorsementMetadata{
									Delegate: "some_other_baker",
								},
							},
						},
					},
				},
			},
		},
	})
	assert.Equal(t, false, ok)
}

func Test_isBakeSuccessful(t *testing.T) {
	m := MissedOpportunityNotifier{
		baker: "some_baker",
	}
	ok := m.isBakeSuccessful(&gotezos.Block{
		Metadata: gotezos.Metadata{
			Baker: "some_baker",
		},
	})
	assert.Equal(t, true, ok)

	ok = m.isBakeSuccessful(&gotezos.Block{
		Metadata: gotezos.Metadata{
			Baker: "some_other_baker",
		},
	})
	assert.Equal(t, false, ok)
}

func Test_getRights(t *testing.T) {
	type input struct {
		m MissedOpportunityNotifier
	}

	type want struct {
		err         bool
		errContains string
		erights     *gotezos.EndorsingRights
		brights     *gotezos.BakingRights
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"is successful",
			input{
				MissedOpportunityNotifier{
					gt: &test.GoTezosMock{},
				},
			},
			want{
				false,
				"",
				&gotezos.EndorsingRights{
					{
						Level:    100,
						Delegate: "some_delegate",
					},
				},
				&gotezos.BakingRights{
					{
						Level:    100,
						Delegate: "some_delegate",
					},
				},
			},
		},
		{
			"handles failure to get endorsing rights",
			input{
				MissedOpportunityNotifier{
					gt: &test.GoTezosMock{
						EndorsingRightsErr: true,
					},
				},
			},
			want{
				true,
				"failed to get endorsing rights",
				&gotezos.EndorsingRights{},
				&gotezos.BakingRights{},
			},
		},
		{
			"handles failure to get baking rights",
			input{
				MissedOpportunityNotifier{
					gt: &test.GoTezosMock{
						BakingRightsErr: true,
					},
				},
			},
			want{
				true,
				"failed to get baking rights",
				&gotezos.EndorsingRights{},
				&gotezos.BakingRights{},
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			brights, erights, err := tt.input.m.getRights(0)
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
			assert.Equal(t, tt.want.erights, erights)
			assert.Equal(t, tt.want.brights, brights)
		})
	}

}

func Test_notify(t *testing.T) {
	type input struct {
		notifier ClientIFace
	}

	type want struct {
		err         bool
		errContains string
	}

	cases := []struct {
		name  string
		input input
		want  want
	}{
		{
			"handles notifier failure",
			input{
				notifier: &MockClient{
					WantSendErr: true,
				},
			},
			want{
				true,
				"failed to send message",
			},
		},
		{
			"is successful",
			input{
				notifier: &MockClient{},
			},
			want{
				false,
				"",
			},
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			notifier := MissedOpportunityNotifier{
				notifiers: []ClientIFace{
					tt.input.notifier,
				},
			}

			err := notifier.notify("bogus message")
			test.CheckErr(t, tt.want.err, tt.want.errContains, err)
		})
	}
}
