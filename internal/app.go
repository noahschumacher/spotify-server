package spotify

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func buildScopes(s ...string) string {
	scopes := ""
	for _, i := range s {
		scopes += i + " "
	}
	return scopes
}

func buildAppAuthorizationString() string {
	com := fmt.Sprintf("%s:%s", os.Getenv("SPOTIFY_ID"), os.Getenv("SPOTIFY_SECRET"))
	return "Basic " + b64.StdEncoding.EncodeToString([]byte(com))
}

type AccessResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func requestAccessToken(code string) *AccessResponse {

	req, err := http.NewRequest("POST", TokenURL, nil)
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

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	var a AccessResponse
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		json.NewDecoder(res.Body).Decode(&a)
	}
	return &a
}

func refreshAccessToken(refreshToken string) *AccessResponse {
	req, err := http.NewRequest("POST", TokenURL, nil)
	if err != nil {
		log.Fatalf("Error creating request")
	}

	req.Header.Set("Authorization", buildAppAuthorizationString())
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("refresh_token", refreshToken)
	q.Add("grant_type", "refresh_token")

	req.URL.RawQuery = q.Encode()

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
	}
	return &a
}
