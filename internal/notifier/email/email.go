package email

import "gopkg.in/gomail.v2"

type IFace interface {
	Send(msg string) error
}

type Client struct {
	dialer *gomail.Dialer
	emails []string
}

func NewEmailClient(username, password string, emails ...string) IFace {
	return &Client{
		dialer: gomail.NewDialer("smtp.example.com", 587, username, password),
		emails: emails,
	}
}

func (c *Client) Send(msg string) error {
	return nil
}
