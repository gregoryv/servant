package github

import (
	"encoding/json"
	"net/http"

	"github.com/gregoryv/servant/htsec"
)

func Contact(c *http.Client) (*htsec.Contact, error) {
	r, _ := http.NewRequest(
		"GET", "https://api.github.com/user", nil,
	)
	r.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := c.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var u htsec.Contact
	if err := json.NewDecoder(resp.Body).Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}
