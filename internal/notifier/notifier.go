package notifier

import (
	"fmt"
	"time"

	"github.com/goat-systems/go-tezos/v3/rpc"
	log "github.com/sirupsen/logrus"
)

// Notifier -
type Notifier interface {
	Start()
}

// MissedOpportunityNotifier -
type MissedOpportunityNotifier struct {
	notifiers []ClientIFace
	rpcClient rpc.IFace
	baker     string
}

// MissedOpportunityNotifierInput -
type MissedOpportunityNotifierInput struct {
	Notifiers []ClientIFace
	RPCClient rpc.IFace
	Baker     string
}

// PayoutNotifierInput -
type PayoutNotifierInput struct {
	Notifiers []ClientIFace
}

// PayoutNotifier -
type PayoutNotifier struct {
	notifiers []ClientIFace
}

type rights struct {
	baking    rpc.BakingRights
	endorsing rpc.EndorsingRights
}

/*
NewMissedOpportunityNotifier -

A notification process that watches for missed endorsement and baking opportunities
and notifies you via SMS (twilio) or Email.
*/
func NewMissedOpportunityNotifier(input MissedOpportunityNotifierInput) Notifier {
	return &MissedOpportunityNotifier{input.Notifiers, input.RPCClient, input.Baker}
}

/*
NewPayoutNotifier -

A notification process that will automatically tweet, email, or text payout notifications.
*/
func NewPayoutNotifier(input PayoutNotifierInput) PayoutNotifier {
	return PayoutNotifier{
		input.Notifiers,
	}
}

// Notify -
func (p *PayoutNotifier) Notify(msg string) error {
	for _, notifier := range p.notifiers {
		if err := notifier.Send(msg); err != nil {
			return err
		}
	}

	return nil
}

func (m *MissedOpportunityNotifier) Start() {
	currentCycle := 0
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			block, err := m.rpcClient.Head()
			if err != nil {
				log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get current cycle")
			}

			if currentCycle < block.Metadata.Level.Cycle {
				brights, erights, err := m.getRights(block.Metadata.Level.Cycle)
				if err != nil {
					log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get rights")
				} else {
					m.rightsWorker(rights{
						baking:    *brights,
						endorsing: *erights,
					})
				}
			}
		}
	}()
}

func (m *MissedOpportunityNotifier) rightsWorker(r rights) {
	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			head, err := m.rpcClient.Head()
			if err != nil {
				log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get current cycle")
			}

			if head.Metadata.Level.Level > r.endorsing[len(r.endorsing)-1].Level && head.Metadata.Level.Level > r.baking[len(r.baking)-1].Level {
				ticker.Stop()
				break
			}

			for _, right := range r.endorsing {
				if head.Metadata.Level.Level == right.Level {
					if ok := m.isEndorsementSuccessful(head); !ok {
						if err := m.notify(fmt.Sprintf("[TZPAY]: Endorsement missed at level: %d", right.Level)); err != nil {
							log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to notify")
						}
					}
				}
			}

			for _, right := range r.baking {
				if head.Metadata.Level.Level == right.Level {
					if ok := m.isBakeSuccessful(head); !ok {
						if err := m.notify(fmt.Sprintf("[TZPAY]: Baking right missed at level: %d", right.Level)); err != nil {
							log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to notify")
						}
					}
				}
			}
		}
	}()
}

func (m *MissedOpportunityNotifier) isEndorsementSuccessful(block *rpc.Block) bool {
	for _, operations := range block.Operations {
		for _, op := range operations {
			for _, endorsement := range op.Contents.Organize().Endorsements {
				if endorsement.Metadata.Delegate == m.baker {
					return true
				}
			}
		}
	}

	return false
}

func (m *MissedOpportunityNotifier) isBakeSuccessful(block *rpc.Block) bool {
	if block.Metadata.Baker == m.baker {
		return true
	}

	return false
}

func (m *MissedOpportunityNotifier) getRights(cycle int) (*rpc.BakingRights, *rpc.EndorsingRights, error) {
	brights, err := m.rpcClient.BakingRights(rpc.BakingRightsInput{
		Cycle:       cycle,
		MaxPriority: 0,
		Delegate:    m.baker,
	})
	if err != nil {
		return &rpc.BakingRights{}, &rpc.EndorsingRights{}, err
	}

	erights, err := m.rpcClient.EndorsingRights(rpc.EndorsingRightsInput{
		Cycle:    cycle,
		Delegate: m.baker,
	})
	if err != nil {
		return &rpc.BakingRights{}, &rpc.EndorsingRights{}, err
	}

	return brights, erights, nil
}

func (m *MissedOpportunityNotifier) notify(msg string) error {
	for _, notifier := range m.notifiers {
		if err := notifier.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
