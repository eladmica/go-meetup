package meetup

import (
	"fmt"
)

type Group struct {
	ID                   int             `json:"id"`
	Name                 string          `json:"name"`
	Link                 string          `json:"link"`
	URLName              string          `json:"urlname"`
	Description          string          `json:"description"`
	Created              int             `json:"created"`
	City                 string          `json:"city"`
	Country              string          `json:"country"`
	LocalizedCountryName string          `json:"localized_country_name"`
	State                string          `json:"state"`
	JoinMode             string          `json:"join_mode"`
	Visibility           string          `json:"visibility"`
	Lat                  float64         `json:"lat"`
	Lon                  float64         `json:"lon"`
	Members              int             `json:"members"`
	Who                  string          `json:"who"`
	Timezone             string          `json:"timezone"`
	WelcomeMessage       string          `json:"welcome_message"`
	NextEvent            *GroupNextEvent `json:"next_event"`
	Organizer            *GroupOrganizer `json:"organizer"`
	Photo                *GroupPhoto     `json:"group_photo"`
	KeyPhoto             *GroupKeyPhoto  `json:"key_photo"`
	Category             *GroupCategory  `json:"category"`
}

type GroupNextEvent struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	YesRSVPCount int    `json:"yes_rsvp_count"`
	Time         int64  `json:"time"`
	UTCOffset    int    `json:"utc_offset"`
}

type GroupOrganizer struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Bio   string `json:"bio"`
	Photo struct {
		ID          int    `json:"id"`
		HighResLink string `json:"highres_link"`
		PhotoLink   string `json:"photo_link"`
		ThumbLink   string `json:"thumb_link"`
		Type        string `json:"type"`
		BaseURL     string `json:"base_url"`
	} `json:"photo"`
}

type GroupPhoto struct {
	ID          int    `json:"id"`
	HighResLink string `json:"highres_link"`
	PhotoLink   string `json:"photo_link"`
	ThumbLink   string `json:"thumb_link"`
	Type        string `json:"type"`
	BaseURL     string `json:"base_url"`
}

type GroupKeyPhoto struct {
	ID          int    `json:"id"`
	HighResLink string `json:"highres_link"`
	PhotoLink   string `json:"photo_link"`
	ThumbLink   string `json:"thumb_link"`
	Type        string `json:"type"`
	BaseURL     string `json:"base_url"`
}

type GroupCategory struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortname"`
	SortName  string `json:"sort_name"`
}

func (c *Client) GetGroup(urlName string) (*Group, error) {
	url := fmt.Sprintf("%v/%v", c.BaseURL, urlName)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var group *Group
	err = c.Do(req, &group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (c *Client) GetSimilarGroups(urlName string) ([]*Group, error) {
	url := fmt.Sprintf("%v/%v/similar_groups", c.BaseURL, urlName)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var groups []*Group
	err = c.Do(req, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *Client) FindGroups() ([]*Group, error) {
	url := fmt.Sprintf("%v/find/groups", c.BaseURL)

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var groups []*Group
	err = c.Do(req, &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}
