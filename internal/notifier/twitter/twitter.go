package twitter

import (
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

// Client -
type Client struct {
	twc *twitter.Client
}

// NewClient -
func NewClient(consumerKey, consumerSecret, accessToken, accessSecret string) *Client {
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return &Client{
		twitter.NewClient(httpClient),
	}
}

// Send -
func (c *Client) Send(msg string) error {
	if _, _, err := c.twc.Statuses.Update(msg, nil); err != nil {
		return err
	}
	return nil
}
