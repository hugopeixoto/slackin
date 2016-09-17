package main

import "net/http"
import "net/url"
import "encoding/json"
import "errors"

type InviteData struct {
	Error string `json:"error"`
}

type PresenceData struct {
	Members []struct {
		Presence string `json:"presence"`
	} `json:"members"`
}

func Invite(name string, token string, email string) error {
	path := "https://" + name + ".slack.com/api/users.admin.invite"

	resp, err := http.PostForm(path, url.Values{
		"email":      {email},
		"token":      {token},
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

func TotalUsers(name string, token string) (int, int, error) {
	path := "https://" + name + ".slack.com/api/users.list"

	resp, err := http.PostForm(path, url.Values{
		"token":    {token},
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
