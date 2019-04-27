package reddit

import (
	"fmt"
	"strconv"

	"github.com/turnage/graw/reddit"
)

// Bot is a wrapper for a reddit session to post payout results to
type Bot struct {
	sub     string
	title   string
	session reddit.Bot
}

// NewRedditSession takes in account name and password and returns a Bot containing the reddit session
func NewRedditSession(agentFile, sub, title string) (*Bot, error) {
	reddit, err := reddit.NewBotFromAgentFile(agentFile, 0)
	if err != nil {
		return nil, err
	}

	return &Bot{sub: sub, title: title, session: reddit}, nil
}

// Post posts a tzscan link to the ophash
func (bot *Bot) Post(ophash string, cycle int) error {
	ophash = ophash[1 : len(ophash)-2]
	link := "https://tzscan.io/" + ophash
	title := bot.title + fmt.Sprintf(" Payout for Cycle "+strconv.Itoa(cycle))
	err := bot.session.PostLink(bot.sub, title, link)
	if err != nil {
		return err
	}

	return nil
}
