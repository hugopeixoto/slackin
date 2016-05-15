package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type Slackin struct {
	Slack Slack
}

type Settings struct {
	Listen string `json:"listen"`
	Name   string `json:"name"`
	Host   string `json:"host"`
	Key    string `json:"key"`
}

func (si *Slackin) ShowIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}

	RenderForm(si.Slack, w)
}

func (si *Slackin) RequestInvite(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/request-invite" {
		w.WriteHeader(404)
		return
	}

	r.ParseForm()

	email := r.PostForm.Get("email")

	log.Printf("[%v] requesting invite\n", email)
	err := si.Slack.Invite(email)

	if err == nil {
		log.Printf("[%v] OK\n", email)
		RenderSuccess(si.Slack, w)
	} else {
		log.Printf("[%v] error: %v\n", email, err)
		body := ``
		switch err.Error() {
		case "already_in_team":
			body = `Already in team`
		case "already_invited":
			body = `Already invited. Check your inbox`
		case "invalid_email":
			body = `Invalid email`
		default:
			body = `Internal server error. Please try again later`
		}
		RenderError(si.Slack, w, body)
	}
}

func ReadJSON(path string, obj interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()
	return json.NewDecoder(file).Decode(obj)
}

func main() {
	settings := Settings{}
	ReadJSON("config/settings.json", &settings)

	si := Slackin{
		Slack{
			settings.Host,
			settings.Key,
			settings.Name,
		},
	}

	http.HandleFunc("/", si.ShowIndex)
	http.HandleFunc("/request-invite", si.RequestInvite)

	log.Fatal(http.ListenAndServe(settings.Listen, nil))
}
