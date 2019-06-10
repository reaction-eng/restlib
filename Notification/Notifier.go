package Notification

import (
	"bitbucket.org/reidev/restlib/users"
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
