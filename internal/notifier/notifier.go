package notifier

import (
	"fmt"
	"time"

	gotezos "github.com/goat-systems/go-tezos/v3"
	log "github.com/sirupsen/logrus"
)

// Notifier -
type Notifier interface {
	Start()
}

// MissedOpportunityNotifier -
type MissedOpportunityNotifier struct {
	notifiers []ClientIFace
	gt        gotezos.IFace
	baker     string
}

// MissedOpportunityNotifierInput -
type MissedOpportunityNotifierInput struct {
	Notifiers []ClientIFace
	Gt        gotezos.IFace
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

/*
NewMissedOpportunityNotifier -

A notification process that watches for missed endorsement and baking opportunities
and notifies you via SMS (twilio) or Email.
*/
func NewMissedOpportunityNotifier(input MissedOpportunityNotifierInput) Notifier {
	return &MissedOpportunityNotifier{input.Notifiers, input.Gt, input.Baker}
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

// Start -
func (m *MissedOpportunityNotifier) Start() {
	ticker := time.NewTicker(time.Minute)
	block, err := m.gt.Head()
	if err != nil {
		log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get current cycle")
	}
	currentCycle := block.Metadata.Level.Cycle

	brights, erights, err := m.getRights(block.Metadata.Level.Cycle)
	if err != nil {
		log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get rights")
	}
	m.monitorRights(*erights, *brights)

	go func() {
		for range ticker.C {
			block, err := m.gt.Head()
			if err != nil {
				log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get current cycle")
			} else {
				if block.Metadata.Level.Cycle > currentCycle {
					brights, erights, err := m.getRights(block.Metadata.Level.Cycle)
					if err != nil {
						log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get rights")
					}
					m.monitorRights(*erights, *brights)
				}
			}
		}
	}()
}

func (m *MissedOpportunityNotifier) monitorRights(endorsementRights gotezos.EndorsingRights, bakingRights gotezos.BakingRights) {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for range ticker.C {
			head, err := m.gt.Head()
			if err != nil {
				log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to get current head")
			}

			if head.Metadata.Level.Level > endorsementRights[len(endorsementRights)-1].Level && head.Metadata.Level.Level > bakingRights[len(bakingRights)-1].Level {
				ticker.Stop()
				break
			}

			for _, right := range endorsementRights {
				if head.Metadata.Level.Level == right.Level {
					if ok := m.isEndorsementSuccessful(head); !ok {
						if err := m.notify(fmt.Sprintf("[TZPAY]: Endorsement missed at level: %d", right.Level)); err != nil {
							log.WithField("error", err.Error()).Error("MissedOpportunityNotifier failed to notify")
						}
					}
				}
			}

			for _, right := range bakingRights {
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

func (m *MissedOpportunityNotifier) isEndorsementSuccessful(block *gotezos.Block) bool {
	for _, operations := range block.Operations {
		for _, op := range operations {
			for _, endorsement := range op.Contents.Endorsements {
				if endorsement.Metadata.Delegate == m.baker {
					return true
				}
			}
		}
	}

	return false
}

func (m *MissedOpportunityNotifier) isBakeSuccessful(block *gotezos.Block) bool {
	if block.Metadata.Baker == m.baker {
		return true
	}

	return false
}

func (m *MissedOpportunityNotifier) getRights(cycle int) (*gotezos.BakingRights, *gotezos.EndorsingRights, error) {
	brights, err := m.gt.BakingRights(gotezos.BakingRightsInput{
		Cycle:       cycle,
		MaxPriority: 0,
		Delegate:    m.baker,
	})
	if err != nil {
		return &gotezos.BakingRights{}, &gotezos.EndorsingRights{}, err
	}

	erights, err := m.gt.EndorsingRights(gotezos.EndorsingRightsInput{
		Cycle:    cycle,
		Delegate: m.baker,
	})
	if err != nil {
		return &gotezos.BakingRights{}, &gotezos.EndorsingRights{}, err
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
