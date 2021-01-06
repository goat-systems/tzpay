package notifier

import "github.com/pkg/errors"

// ClientIFace is an interface to a client that sends messages (twitter, email, twilio)
type ClientIFace interface {
	Send(msg string) error
}

// MockClient mocks twilio.IFace
type MockClient struct {
	WantSendErr bool
}

// Send satisfies twilio.IFace
func (m *MockClient) Send(msg string) error {
	if m.WantSendErr {
		return errors.New("failed to send message")
	}

	return nil
}
