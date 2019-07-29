// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package notification

import (
	"github.com/reaction-eng/restlib/users"
	"log"
)

type NotifierStruct struct {
	NotfierType string
}

type Notifier interface {
	Notify(notification Notification, user users.User) error
}

///////////////////////
type dummyNotifier struct {
	NotfierType string
}

func NewDummyNotifier() *dummyNotifier {

	newDummy := dummyNotifier{
		NotfierType: "Dummy, lol",
	}
	return &newDummy
}

func (notif *dummyNotifier) Notify(notification Notification, user users.User) error {

	log.Println("Sending message...->", notification.Message)

	return nil
}
