package twilio

import (
	"github.com/pkg/errors"
	"github.com/sfreiberg/gotwilio"
)

// IFace is an interface to a client that sends SMS messages via twilio
type IFace interface {
	Send(msg string) error
}

// Client is a twilio client to send SMS messages
type Client struct {
	AccountSID string
	AuthToken  string
	From       string
	To         []string
	twilio     *gotwilio.Twilio
}

// New returns a new twilio IFace
func New(twilio Client) IFace {
	twilio.twilio = gotwilio.NewTwilioClient(twilio.AccountSID, twilio.AuthToken)
	return &twilio
}

// Send sends an SMS message
// TODO imporve error handling
func (c *Client) Send(msg string) error {
	for _, num := range c.To {
		_, _, err := c.twilio.SendSMS(c.From, num, msg, "", "")
		if err != nil {
			return errors.Wrapf(err, "failed to send message through twilio to '%s'", num)
		}
	}

	return nil
}
