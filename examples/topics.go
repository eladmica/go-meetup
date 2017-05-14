package main

import (
	"log"

	"github.com/eladmica/go-meetup/meetup"
)

// Find and print topics related to "Software".
func findTopics() {
	client := meetup.NewClient(nil)
	client.Authentication = meetup.NewKeyAuth("SECRET_KEY") // Authentication required for this endpoint.

	title := "Software"
	topics, err := client.FindTopics(title)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Topics related to \"%v\":\n", title)
	for _, topic := range topics[:10] {
		log.Printf("Topic: %v\nMembers: %v\nDescription: %v\n", topic.Name, topic.MemberCount, topic.Description)
	}
}
