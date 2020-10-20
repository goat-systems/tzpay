package payout

import (
	"sync"
	"testing"

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
