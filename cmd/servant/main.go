package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gregoryv/servant"
	"github.com/gregoryv/servant/htsec"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func main() {
	sys := servant.NewSystem()
	sec := htsec.NewSecure()
	sec.Use(&htsec.Auth{
		Config: &oauth2.Config{
			RedirectURL:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
			ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
			ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
			Endpoint:     endpoints.GitHub,
		},
		ReadUser: readGithubUser,
	})
	srv := http.Server{
		Addr:    ":8100",
		Handler: servant.NewRouter(sys, sec),
	}

	fmt.Println("listen", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func readGithubUser(token *oauth2.Token) (*htsec.User, error) {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", "token "+token.AccessToken)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var u htsec.User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
