package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	spotify "github.com/noahschumacher/spotify/internal"
)

var (
	authCallbackURL = "http://localhost:8080/callback"
)

func buildScopes(s ...string) string {
	scopes := ""
	for _, i := range s {
		scopes += i + " "
	}
	return scopes
}

func authUser() {
	req, err := http.NewRequest("GET", spotify.AuthURL, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}

	q := req.URL.Query()
	q.Add("response_type", "code")
	q.Add("client_id", os.Getenv("SPOTIFY_ID"))
	q.Add("scope", buildScopes(spotify.ScopeUserReadPrivate, spotify.ScopeUserLibraryRead))
	q.Add("redirect_uri", authCallbackURL)
	q.Add("state", spotify.State)

	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())
}

func authCallback(w http.ResponseWriter, r *http.Request) {
	u, _ := url.ParseRequestURI(r.RequestURI)
	m, _ := url.ParseQuery(u.RawQuery)
	code := m["code"][0]
	fmt.Println("Code: ", code)

	requestAccessToken(code)
}

func buildAppAuthorizationString() string {
	com := fmt.Sprintf("%s:%s", os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
	return "Basic " + b64.StdEncoding.EncodeToString([]byte(com))
}

func requestAccessToken(code string) {

	req, err := http.NewRequest("POST", spotify.TokenURL, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}

	req.Header.Set("Authorization", buildAppAuthorizationString())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("code", code)
	q.Add("redirect_uri", authCallbackURL)
	q.Add("grant_type", "authorization_code")

	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var a AccessResponse
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		json.NewDecoder(res.Body).Decode(&a)
		a.print()
	}
}

type AccessResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (a AccessResponse) print() {
	fmt.Println("AccessToken:", a.AccessToken)
	fmt.Println("TokenType:", a.TokenType)
	fmt.Println("Scope:", a.Scope)
	fmt.Println("ExpiresIn:", a.ExpiresIn)
	fmt.Println("RefreshToken:", a.RefreshToken)
}

func refreshAccessToken(refreshToken string) {
	req, err := http.NewRequest("POST", spotify.TokenURL, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}

	req.Header.Set("Authorization", buildAppAuthorizationString())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")

	req.URL.RawQuery = q.Encode()
	fmt.Println(req.URL.String())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var a AccessResponse
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		json.NewDecoder(res.Body).Decode(&a)
		a.RefreshToken = refreshToken
		a.print()
	}
}

func main() {
	fmt.Println("Do it myself!")

	// refreshAccessToken("AQB2oQyt31Z_SqNHZ-SRPyI-d8XrF_dxq2QpRNQJD6xCpzP5Uv7WMfzgAmAT3Kt2ETHP1D1VkjHL4p9irzgy-_kFMh-0TUQugKrt260LjgKGo97gx9BfvgSohWRF0-VUie4")

	// User needs to be redirected to the outputed URL
	authUser()

	// Start the http server
	http.HandleFunc("/callback", authCallback)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error:%s", err)
	}
}
