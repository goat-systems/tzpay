package payout

import (
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/goat-systems/tzpay/v3/internal/tzkt"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func Test_Size(t *testing.T) {
	q := Queue{}
	assert.Equal(t, 0, q.Size())

	q = Queue{
		payouts: []Payout{
			{},
			{},
		},
	}
	assert.Equal(t, 2, q.Size())
}

func Test_Empty(t *testing.T) {
	q := Queue{}
	assert.True(t, q.Empty())
	q = Queue{
		payouts: []Payout{
			{},
		},
	}
	assert.False(t, q.Empty())
}

func Test_Start(t *testing.T) {
	type input struct {
		payouts []Payout
	}

	type want struct {
	}

	cases := []struct {
		name         string
		input        input
		successCount int
	}{
		{
			"is successful",
			input{
				payouts: []Payout{
					{
						cycle: 10,
						constructPayoutFunc: func() (tzkt.RewardsSplit, error) {
							return tzkt.RewardsSplit{Cycle: 10}, nil
						},
						applyFunc: func(delegators tzkt.Delegators) ([]string, error) {
							return []string{}, nil
						},
					},
					{
						cycle: 11,
						constructPayoutFunc: func() (tzkt.RewardsSplit, error) {
							return tzkt.RewardsSplit{Cycle: 11}, nil
						},
						applyFunc: func(delegators tzkt.Delegators) ([]string, error) {
							return []string{}, nil
						},
					},
					{
						cycle: 12,
						constructPayoutFunc: func() (tzkt.RewardsSplit, error) {
							return tzkt.RewardsSplit{Cycle: 12}, nil
						},
						applyFunc: func(delegators tzkt.Delegators) ([]string, error) {
							return []string{}, nil
						},
					},
				},
			},
			3,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewQueue(nil)
			queue.tickerDuration = time.Millisecond
			logger, hook := test.NewNullLogger()
			queue.logger = logger
			queue.Start()
			for _, payout := range tt.input.payouts {
				queue.Enqueue(payout)
			}
			time.Sleep(time.Second * 1)

			count := 0
			if len(hook.Entries) > 0 {
				for _, entry := range hook.Entries {
					if strings.Contains(entry.Message, "Payout successfully executed.") {
						count++
					}
				}
			}

			assert.Equal(t, tt.successCount, count)
		})
	}
}

func Test_Front(t *testing.T) {
	q := Queue{
		mu: &sync.Mutex{},
	}
	p, err := q.Front()
	assert.Error(t, err)
	assert.Equal(t, Payout{}, p)

	q = Queue{
		mu: &sync.Mutex{},
		payouts: []Payout{
			{
				cycle: 10,
			},
			{},
		},
	}
	p, err = q.Front()
	assert.Nil(t, err)
	assert.Equal(t, Payout{
		cycle: 10,
	}, p)
}

func Test_Dequeue(t *testing.T) {
	q := Queue{
		mu: &sync.Mutex{},
	}
	err := q.Dequeue()
	assert.Error(t, err)

	q = Queue{
		mu: &sync.Mutex{},
		payouts: []Payout{
			{
				cycle: 10,
			},
			{},
		},
	}
	err = q.Dequeue()
	assert.Nil(t, err)
}

func Test_Enqueue(t *testing.T) {
	q := Queue{
		mu: &sync.Mutex{},
	}
	q.Enqueue(Payout{})
	assert.Equal(t, 1, q.Size())
}
