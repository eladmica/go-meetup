package meetup

import (
	"net/http"
)

const (
	baseURL = "https://api.meetup.com/"
)

// Client is used to communicate with the Meetup API
type Client struct {
	// Private API key used for authentication
	Key string

	// HTTP client used to communicate with the API.
	client *http.Client
}

// NewClient returns a new Meetup API client
func NewClient(key string) *Client {
	return &Client{
		Key:    key,
		client: &http.Client{},
	}
}
