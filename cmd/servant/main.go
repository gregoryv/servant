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
	guard := htsec.NewGuard()
	guard.Include(
		&htsec.Gate{
			Config: &oauth2.Config{
				RedirectURL:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
				ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
				ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
				Endpoint:     endpoints.GitHub,
			},
			Contact: github.Contact,
		},
		&htsec.Gate{
			Config: &oauth2.Config{
				RedirectURL:  os.Getenv("OAUTH_GOOGLE_REDIRECT_URI"),
				ClientID:     os.Getenv("OAUTH_GOOGLE_CLIENTID"),
				ClientSecret: os.Getenv("OAUTH_GOOGLE_SECRET"),
				Scopes:       []string{"profile", "email"},
				Endpoint:     endpoints.Google,
			},
			Contact: google.Contact,
		},
	)

	sys := servant.NewSystem()
	sys.SetSecurity(guard)

	srv := http.Server{
		Addr:    ":8100",
		Handler: servant.NewRouter(sys),
	}

	fmt.Println("listen", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
