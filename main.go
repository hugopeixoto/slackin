package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type SlackInfo struct {
	Name        string `json:"name"`
	Token       string `json:"token"`
	DisplayName string `json:"display_name"`
	Hostname    string `json:"hostname"`
}

type Settings struct {
	Listen string      `json:"listen"`
	Slacks []SlackInfo `json:"slacks"`
}

func SlackScope(settings Settings, fn func(SlackInfo, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, s := range settings.Slacks {
			if s.Hostname == r.Host {
				fn(s, w, r)
				return
			}
		}

		w.WriteHeader(404)
		w.Write([]byte(`nope`))
	}
}

func ShowIndex(si SlackInfo, w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		return
	}

	RenderForm(si, w)
}

func RequestInvite(si SlackInfo, w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/request-invite" {
		w.WriteHeader(404)
		return
	}

	r.ParseForm()

	email := r.PostForm.Get("email")

	log.Printf("[%v] requesting invite\n", email)
	err := Invite(si.Name, si.Token, email)

	if err == nil {
		log.Printf("[%v] OK\n", email)
		RenderSuccess(si, w)
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
		RenderError(si, w, body)
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

	http.HandleFunc("/", SlackScope(settings, ShowIndex))
	http.HandleFunc("/request-invite", SlackScope(settings, RequestInvite))

	log.Fatal(http.ListenAndServe(settings.Listen, nil))
}
