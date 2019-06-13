package Notification

import (
	"bitbucket.org/reidev/restlib/users"
	"context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"io/ioutil"
	"log"
	"net/http"
)

type CalendarNotifier struct {
	NotifierType string
	Serv         *calendar.Service
}

func (notif *CalendarNotifier) Notify(notification Notification, user users.User) error {

	newEvent := &calendar.Event{
		Summary:     "REI Event",
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

	_, err := notif.Serv.Events.Insert("primary", newEvent).Do()
	if err != nil {
		log.Println("Error in calendar Notify ->", err.Error())
	}

	//getting all events
	//t := time.Now().Format(time.RFC3339)
	//events, err := notif.Serv.Events.List("primary").ShowDeleted(false).
	//	SingleEvents(true).TimeMin(t).MaxResults(10).OrderBy("startTime").Do()
	//if err != nil {
	//	log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	//}
	//log.Println("Upcoming events:")
	//if len(events.Items) == 0 {
	//	log.Println("No upcoming events found.")
	//} else {
	//	for _, item := range events.Items {
	//		date := item.Start.DateTime
	//		if date == "" {
	//			date = item.Start.Date
	//		}
	//		log.Printf("%v (%v)\n", item.Summary, date)
	//	}
	//}
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
