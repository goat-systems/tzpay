package notifier

import (
	"testing"

	"github.com/goat-systems/go-tezos/v3/rpc"
	"github.com/goat-systems/tzpay/v3/internal/test"
	"github.com/stretchr/testify/assert"
)

func Test_isEndorsementSuccessful(t *testing.T) {
	m := MissedOpportunityNotifier{
		baker: "some_baker",
	}
	ok := m.isEndorsementSuccessful(&rpc.Block{
		Operations: [][]rpc.Operations{
			{
				{
					Contents: rpc.Contents{
						{
							Metadata: &rpc.ContentsHelperMetadata{
								Delegate: "some_baker",
							},
						},
					},
				},
			},
		},
	})
	assert.Equal(t, true, ok)

	ok = m.isEndorsementSuccessful(&rpc.Block{
		Operations: [][]rpc.Operations{
			{
				{
					Contents: rpc.Contents{
						{
							Metadata: &rpc.ContentsHelperMetadata{
								Delegate: "some_other_baker",
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
	ok := m.isBakeSuccessful(&rpc.Block{
		Metadata: rpc.Metadata{
			Baker: "some_baker",
		},
	})
	assert.Equal(t, true, ok)

	ok = m.isBakeSuccessful(&rpc.Block{
		Metadata: rpc.Metadata{
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
		erights     *rpc.EndorsingRights
		brights     *rpc.BakingRights
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
					rpcClient: &test.RPCMock{},
				},
			},
			want{
				false,
				"",
				&rpc.EndorsingRights{
					{
						Level:    100,
						Delegate: "some_delegate",
					},
				},
				&rpc.BakingRights{
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
					rpcClient: &test.RPCMock{
						EndorsingRightsErr: true,
					},
				},
			},
			want{
				true,
				"failed to get endorsing rights",
				&rpc.EndorsingRights{},
				&rpc.BakingRights{},
			},
		},
		{
			"handles failure to get baking rights",
			input{
				MissedOpportunityNotifier{
					rpcClient: &test.RPCMock{
						BakingRightsErr: true,
					},
				},
			},
			want{
				true,
				"failed to get baking rights",
				&rpc.EndorsingRights{},
				&rpc.BakingRights{},
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
