package main

import (
	"log"

	"github.com/eladmica/go-meetup/meetup"
)

// Print events hosted by the NY Tech Meetup group.
func groupEvents() {
	client := meetup.NewClient(nil)

	events, err := client.GetEvents("ny-tech", nil)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("NY Tech Meetup group events:")
	for _, event := range events {
		log.Println(event.Name)
	}
}

// Print the top 5 most recommended events near Vancouver, including their respective hosting groups.
func recommendedEvents() {
	client := meetup.NewClient(nil)
	client.Authentication = meetup.NewKeyAuth("SECRET_KEY") // Authentication required for this endpoint.

	params := meetup.GetRecommendedEventsParams{
		Lat: 49.2827,
		Lon: 123.1207,
	}

	events, err := client.GetRecommendedEvents(&params)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Top 5 recommended events nearby:")
	for _, event := range events[:5] {
		group := "unknown"
		if event.Group != nil {
			group = event.Group.Name
		}
		log.Printf("%v, held by: %v\n", event.Name, group)
	}
}
