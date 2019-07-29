// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package Notification

import (
	"github.com/reaction-eng/restlib/users"
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

func (notif *SlackNotifier) Notify(notification Notification, user users.User) error {

	//include some kind of db call to get userID and SlackAuthto who'm we send to
	//For now it is REI's slack.
	api := slack.New("xoxp-539969220354-621210039441-658880912709-46de6eb8aeacd6fc51b6582ce395eeac")

	//Matt UFX0ZCF4P
	//Grant UJ96615CZ
	//slackbot USLACKBOT
	userID := "USLACKBOT"

	//Connecto to channel.
	_, _, channelID, err := api.OpenIMChannel(userID)

	if err != nil {
		fmt.Printf("%s\n", err)
	}

	//Write to channel.
	_, _, err = api.PostMessage(channelID, slack.MsgOptionText(notification.Message, false))

	return nil
}
