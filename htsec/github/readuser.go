package github

import (
	"encoding/json"
	"net/http"

	"github.com/gregoryv/servant/htsec"
	"golang.org/x/oauth2"
)

func ReadUser(c *http.Client) htsec.ReadUserFunc {
	return func(token *oauth2.Token) (*htsec.User, error) {
		r, _ := http.NewRequest(
			"GET", "https://api.github.com/user", nil,
		)
		r.Header.Set("Accept", "application/vnd.github.v3+json")
		r.Header.Set("Authorization", "token "+token.AccessToken)
		resp, err := c.Do(r)
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
}
