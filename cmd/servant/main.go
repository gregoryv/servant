package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gregoryv/servant"
	"github.com/gregoryv/servant/htsec"
	"github.com/gregoryv/servant/htsec/github"
	"github.com/gregoryv/servant/htsec/google"
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
		ReadUser: github.ReadUser(http.DefaultClient),
	})

	googleConf := &oauth2.Config{
		RedirectURL:  os.Getenv("OAUTH_GOOGLE_REDIRECT_URI"),
		ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENTID"),
		ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
		Scopes:       []string{"profile", "email"},
		Endpoint:     endpoints.Google,
	}
	sec.Use(&htsec.Auth{
		Config:   googleConf,
		ReadUser: google.ReadUser(googleConf),
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
