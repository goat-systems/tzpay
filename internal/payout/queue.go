package payout

import (
	"fmt"
	"sync"
	"time"

	"github.com/goat-systems/tzpay/v2/internal/notifier"
	"github.com/goat-systems/tzpay/v2/internal/print"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	notifier notifier.PayoutNotifier
	payouts  []Payout
	mu       *sync.Mutex
}

func NewQueue(notifier notifier.PayoutNotifier) Queue {
	return Queue{
		notifier: notifier,
	}
}

func (q *Queue) Enqueue(p Payout) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.payouts = append(q.payouts, p)
}

func (q *Queue) Dequeue() error {
	if len(q.payouts) > 0 {
		q.mu.Lock()
		defer q.mu.Unlock()
		q.payouts = q.payouts[1:]
		return nil
	}
	return fmt.Errorf("Pop Error: Queue is empty")
}

func (q *Queue) Front() (Payout, error) {
	if len(q.payouts) > 0 {
		q.mu.Lock()
		defer q.mu.Unlock()
		return q.payouts[0], nil
	}
	return Payout{}, fmt.Errorf("Peep Error: Queue is empty")
}

func (q *Queue) Size() int {
	return len(q.payouts)
}

func (q *Queue) Empty() bool {
	return len(q.payouts) == 0
}

func (q *Queue) Start() {
	logrus.Info("Starting payout queue.")
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			logrus.Info("Popping off payout queue.")
			payout, err := q.Front()
			if err != nil {
				logrus.Info("Payout Queue is empty.")
				continue
			}

			err = q.Dequeue()
			if err != nil {
				logrus.WithField("error", err.Error()).Error("failed to dequeue payout in queue")
				continue
			}
			rewardsSplit, err := payout.Execute()
			if err != nil {
				logrus.WithField("error", err.Error()).Error("failed to execute payout in queue")
				q.Enqueue(payout)
				continue
			}

			err = q.notifier.Notify(fmt.Sprintf("[TZPAY] payout for cycle %d: \n%s\n #tezos #pos", payout.cycle, rewardsSplit.OperationLink))
			if err != nil {
				logrus.WithField("error", err.Error()).Error("Failed to notify.")
			}

			err = print.JSON(rewardsSplit)
			if err != nil {
				logrus.WithField("error", err.Error()).Fatal("Failed to print JSON report.")
			}

		}
	}()
}
