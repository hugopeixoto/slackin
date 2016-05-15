package main

import "net/http"
import "net/url"
import "encoding/json"
import "errors"

type Slack struct {
	Host  string
	Token string
	Name  string
}

type InviteData struct {
	Error string `json:"error"`
}

func (s Slack) Invite(email string) error {
	path := "https://" + s.Host + ".slack.com/api/users.admin.invite"

	resp, err := http.PostForm(path, url.Values{
		"email":      {email},
		"token":      {s.Token},
		"set_active": {"true"},
	})

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var data InviteData
	json.NewDecoder(resp.Body).Decode(&data)

	if data.Error != "" {
		return errors.New(data.Error)
	}

	return nil
}

type PresenceData struct {
	Members []struct {
		Presence string `json:"presence"`
	} `json:"members"`
}

func (s *Slack) TotalUsers() (int, int, error) {
	path := "https://" + s.Host + ".slack.com/api/users.list"

	resp, err := http.PostForm(path, url.Values{
		"token":    {s.Token},
		"presence": {"1"},
	})

	if err != nil {
		return 0, 0, err
	}

	defer resp.Body.Close()

	var data PresenceData

	json.NewDecoder(resp.Body).Decode(&data)

	active := 0
	for _, u := range data.Members {
		if u.Presence == "active" {
			active++
		}
	}

	return len(data.Members), active, nil
}
