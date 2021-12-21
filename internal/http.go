package spotify

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	authCallbackURL = os.Getenv("CALLBACK_URL")
)

type authURL struct {
	URL string `json:"URL"`
}

func AuthURLBuilder(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", AuthURL, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", os.Getenv("SPOTIFY_ID"))
	q.Add("scope", buildScopes(ScopeUserReadPrivate, ScopeUserLibraryRead))
	q.Add("redirect_uri", authCallbackURL)
	q.Add("state", State)

	req.URL.RawQuery = q.Encode()
	aurl := authURL{URL: req.URL.String()}

	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.Encode(aurl)
}

func AuthCallback(w http.ResponseWriter, r *http.Request) {
	u, _ := url.ParseRequestURI(r.RequestURI)
	m, _ := url.ParseQuery(u.RawQuery)
	code := m["code"][0]

	res := requestAccessToken(code)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func UserTrack(w http.ResponseWriter, r *http.Request) {

}
