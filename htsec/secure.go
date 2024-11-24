package htsec

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/endpoints"
)

func NewSecure() *Secure {
	return &Secure{
		src: map[string]*Auth{
			"github": &Auth{
				Config:   githubOauth,
				ReadUser: readGithubUser,
			},
		},
	}
}

type Secure struct {
	src map[string]*Auth
}

func (s *Secure) Names() []string {
	res := make([]string, 0, len(s.src))
	for name, _ := range s.src {
		res = append(res, name)
	}
	sort.Strings(res)
	return res
}

func (s *Secure) AuthService(name string) (*Auth, error) {
	a, found := s.src[name]
	if !found {
		return nil, fmt.Errorf("AuthService %v: %w", name, notFound)
	}
	return a, nil
}

var notFound = fmt.Errorf("not found")

// ----------------------------------------

type Auth struct {
	*oauth2.Config
	ReadUser func(token *oauth2.Token) (*User, error)
}

var githubOauth = &oauth2.Config{
	RedirectURL:  os.Getenv("OAUTH_GITHUB_REDIRECT_URI"),
	ClientID:     os.Getenv("OAUTH_GITHUB_CLIENTID"),
	ClientSecret: os.Getenv("OAUTH_GITHUB_SECRET"),
	Endpoint:     endpoints.GitHub,
}

func readGithubUser(token *oauth2.Token) (*User, error) {
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	r.Header.Set("Authorization", "token "+token.AccessToken)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return nil, err
	}
	var u User
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
