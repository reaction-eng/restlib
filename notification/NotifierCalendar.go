// Copyright 2019 Reaction Engineering International. All rights reserved.
// Use of this source code is governed by the MIT license in the file LICENSE.txt.

package notification

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/reaction-eng/restlib/users"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

type CalendarNotifier struct {
	NotifierType string
	Serv         *calendar.Service
}

/*
All calendar events are created on whichever users calendar is linked to the server and the recipient for the event is added as an 'Attenddee'.
*/
func (notif *CalendarNotifier) Notify(notification Notification, user users.User) error {

	//create new Calendar event.
	newEvent := &calendar.Event{
		Summary:     "REI gEvent",
		Location:    "University of Utah",
		Description: notification.Message,
		Start: &calendar.EventDateTime{
			DateTime: "2019-06-13T15:00:00",
			TimeZone: "America/Denver",
		},
		End: &calendar.EventDateTime{
			DateTime: "2019-06-13T16:30:00",
			TimeZone: "America/Denver",
		},
		Recurrence: []string{"RRULE:FREQ=DAILY;COUNT=1"},
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{Email: user.Email()},
			//&calendar.EventAttendee{Email:"sbrin@example.com"},
		},
	}

	//Send calendar.
	_, err := notif.Serv.Events.Insert("primary", newEvent).Do()
	if err != nil {
		log.Println("Error in calendar Notify ->", err.Error())
	}

	return err
}

// Retrieve a token, saves the token, then returns the generated client.
func getClientCalendar(config *oauth2.Config) *http.Client {
	// The file tokenCalendar.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "tokenCalendar.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func NewCalendarNotifier() *CalendarNotifier {

	b, err := ioutil.ReadFile("Notification/credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved tokenCalendar.json.
	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClientCalendar(config)

	srv, err := calendar.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	newNotifier := CalendarNotifier{
		NotifierType: "Calendar",
		Serv:         srv,
	}

	return &newNotifier

}
