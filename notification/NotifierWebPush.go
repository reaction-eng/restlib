// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package notification

import (
	"database/sql"
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/reaction-eng/restlib/preferences"
	"github.com/reaction-eng/restlib/users"

	//"encoding/json"
	//"github.com/SherClockHolmes/webpush-go"
	"log"
)

type WebPushNotifier struct {
	NotfierType     string
	PreferencesJson string
	db              *sql.DB
}

func NewWebPushNotifier(preferencesJsonFile string, db *sql.DB) *WebPushNotifier {

	newNotifier := WebPushNotifier{
		NotfierType:     "WebPush",
		PreferencesJson: preferencesJsonFile,
		db:              db,
	}
	return &newNotifier
}

func (notif *WebPushNotifier) Notify(notification Notification, user users.User) error {

	//Get user preferences.
	options, err := preferences.LoadOptionsGroup(notif.PreferencesJson)
	if err != nil {
		return err
	}
	sqlConnectiont, err := preferences.NewRepoMySql(notif.db, options)
	usrsPref, err := sqlConnectiont.GetPreferences(user)

	//Get user's subscription to send notification to.
	temp := usrsPref.Settings.SubGroup

	//RawSub will collect the Marshaled subscription.
	RawSub := ""

	//navigate to the webSubscription
	for _, v := range temp {
		thing := v.Settings

		RawSub = thing["webSubscription"]
		if RawSub != "" {
			break
		}
	}

	if err != nil {
		log.Println("Could not get user preferences ->", err.Error())
	}

	//if we still don't have a webSubscription don't continue.
	if RawSub == "" {
		log.Println("User doesn't have a subscription. Notification failed.")
		return nil
	}

	//cast Marshalled json to a webPushSubscription type.
	sub := webpush.Subscription{}
	err = json.Unmarshal([]byte(RawSub), &sub)
	if err != nil {
		log.Println("Error in unmarshalling Subscriptiong ->", err.Error())
	}

	//Send notification to subscription's endpoint using our VAPID keys.
	_, err = webpush.SendNotification(
		[]byte(notification.Message),
		&sub,
		&webpush.Options{
			TTL:             30,
			VAPIDPrivateKey: "isi0niTk-Ej6aR8UAp16WRi4Ghoi8gsx3GAQ71iHpkc",
			VAPIDPublicKey:  "BNh0sBvliZ70vkOqwr5ukFx1mYetOwidow2ELWLiY0hnLPi47VakHhEVkXM02_5j1l3AXkFex6c-JFCPEQHZChw"})

	if err != nil {
		log.Println("Error on notification ->", err.Error())
	}

	return nil
}
