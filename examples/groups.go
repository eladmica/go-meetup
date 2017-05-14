package main

import (
	"log"

	"github.com/eladmica/go-meetup/meetup"
)

// Find and print the newest groups related to "Web Development".
func findGroups() {
	client := meetup.NewClient(nil)
	client.Authentication = meetup.NewKeyAuth("SECRET_KEY") // Authentication required for this endpoint.

	params := meetup.FindGroupsParams{
		Order: "newest",
		Text:  "Web Development",
	}

	groups, err := client.FindGroups(&params)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Newest %v groups:\n", params.Text)
	for _, group := range groups[:10] {
		log.Printf("%v. Link: %v\n", group.Name, group.Link)
	}
}
