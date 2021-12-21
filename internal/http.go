package spotify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
)

var (
	authCallbackURL = "http://localhost:8080/callback"
)

type authURL struct {
	URL string
}

func AuthURLBuilder(w http.ResponseWriter, r *http.Request) {
	req, err := http.NewRequest("GET", AuthURL, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}

	fmt.Println(os.Getenv("SPOTIFY_ID"))

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", os.Getenv("SPOTIFY_ID"))
	q.Add("scope", buildScopes(ScopeUserReadPrivate, ScopeUserLibraryRead))
	q.Add("redirect_uri", authCallbackURL)
	q.Add("state", State)

	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String()) // Good
	aurl := authURL{URL: req.URL.String()}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aurl) // Bad
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
