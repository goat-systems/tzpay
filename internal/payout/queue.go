package payout

import (
	"fmt"
	"sync"
	"time"

	"github.com/goat-systems/tzpay/v3/internal/notifier"
	"github.com/goat-systems/tzpay/v3/internal/print"
	"github.com/sirupsen/logrus"
)

type Queue struct {
	notifier       *notifier.PayoutNotifier
	payouts        []Payout
	mu             *sync.Mutex
	logger         *logrus.Logger
	tickerDuration time.Duration
}

func NewQueue(notifier *notifier.PayoutNotifier) Queue {
	return Queue{
		notifier:       notifier,
		mu:             &sync.Mutex{},
		tickerDuration: time.Minute,
		logger:         logrus.New(),
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
	q.logger.Info("Starting payout queue.")
	go func() {
		ticker := time.NewTicker(q.tickerDuration)
		for range ticker.C {
			q.logger.Info("Popping off payout queue.")
			payout, err := q.Front()
			if err != nil {
				q.logger.Info("Payout Queue is empty.")
				continue
			}

			err = q.Dequeue()
			if err != nil {
				q.logger.WithFields(logrus.Fields{"error": err.Error(), "cycle": payout.cycle}).Error("failed to dequeue payout in queue")
				continue
			}
			rewardsSplit, err := payout.Execute()
			if err != nil {
				q.logger.WithFields(logrus.Fields{"error": err.Error(), "cycle": payout.cycle}).Error("failed to execute payout in queue")
				q.Enqueue(payout)
				continue
			}

			q.logger.WithField("cycle", payout.cycle).Info("Payout successfully executed.")

			if q.notifier != nil {
				err = q.notifier.Notify(fmt.Sprintf("[TZPAY] payout for cycle %d: \n%s\n #tezos #pos", payout.cycle, rewardsSplit.OperationLink))
				if err != nil {
					q.logger.WithField("error", err.Error()).Error("Failed to notify.")
				}
			}

			err = print.JSON(rewardsSplit)
			if err != nil {
				q.logger.WithField("error", err.Error()).Fatal("Failed to print JSON report.")
			}

		}
	}()
}
