package google

import (
	"context"
	"net/http"

	"github.com/gregoryv/servant/htauth"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func Contact(c *http.Client) (*htauth.Contact, error) {
	ctx := context.Background()
	service, err := people.NewService(ctx,
		option.WithHTTPClient(c),
	)
	if err != nil {
		return nil, err
	}

	profile, err := service.People.Get("people/me").PersonFields(
		"names,emailAddresses",
	).Do()
	if err != nil {
		return nil, err
	}

	var u htauth.Contact
	if len(profile.EmailAddresses) > 0 {
		u.Email = profile.EmailAddresses[0].Value
	}
	if len(profile.Names) > 0 {
		n := profile.Names[0]
		u.Name = n.GivenName + " " + n.FamilyName
		if u.Name == "" {
			u.Name = n.DisplayName
		}
	}
	return &u, nil
}
