package Notification

import (
	"fmt"
	"github.com/nlopes/slack"
)

//[]Link to DB
//[]Set up RedirctURL to make public.

type SlackNotifier struct {
	NotifierType string
}

func NewSlackNotifier() *SlackNotifier {

	slackNotifier := SlackNotifier{
		NotifierType: "Slack",
	}
	return &slackNotifier
}

func (notif *SlackNotifier) Notify(notification Notification) error {

	//include some kind of db call to get userID and SlackAuthto who'm we send to
	api := slack.New("xoxp-539969220354-621210039441-658880912709-46de6eb8aeacd6fc51b6582ce395eeac")

	//Matt UFX0ZCF4P
	//Grant UJ96615CZ
	//slackbot USLACKBOT
	userID := "USLACKBOT"

	_, _, channelID, err := api.OpenIMChannel(userID)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	_, _, err = api.PostMessage(channelID, slack.MsgOptionText(notification.Message, false))

	return nil
}
