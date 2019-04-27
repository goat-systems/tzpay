package twitter

import (
	"fmt"
	"strconv"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
)

// Bot is a wrapper for a twitter session to post payout results to
type Bot struct {
	title   string
	session *twitter.Client
}

// NewTwitterSession creates a new twitter bot based off a twitter.yml file located in the current path or path passed.
// can pass an optional title for the bot
func NewTwitterSession(path string, title string) (*Bot, error) {
	bot := Bot{title: title}
	viper.SetConfigName("twitter")
	if path != "" {
		viper.AddConfigPath(path)
	}
	viper.AddConfigPath("./")
	err := viper.ReadInConfig()
	if err != nil {
		return &bot, fmt.Errorf("could not find twitter.yml: %v", err)
	}

	key := viper.GetString("consumerKey")
	keySecret := viper.GetString("consumerKeySecret")
	access := viper.GetString("accessToken")
	accessSecret := viper.GetString("accessTokenSecret")
	if key == "" || access == "" || keySecret == "" || accessSecret == "" {
		return &bot, fmt.Errorf("could not read key or access token")
	}

	config := oauth1.NewConfig(key, keySecret)
	token := oauth1.NewToken(access, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	bot.session = client
	return &bot, nil
}

// Post posts a tzscan link to the ophash
func (bot *Bot) Post(ophash string, cycle int) error {
	ophash = ophash[1 : len(ophash)-2]
	link := "https://tzscan.io/" + ophash
	title := bot.title + fmt.Sprintf(" Payout for Cycle "+strconv.Itoa(cycle)+":")
	_, _, err := bot.session.Statuses.Update(title+" "+link, nil)
	if err != nil {
		return err
	}

	return nil
}
