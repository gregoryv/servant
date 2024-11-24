package google

import (
	"context"

	"github.com/gregoryv/servant/htsec"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"google.golang.org/api/people/v1"
)

func ReadUser(config *oauth2.Config) htsec.ReadUserFunc {
	return func(token *oauth2.Token) (*htsec.User, error) {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token.AccessToken},
		)
		oauthClient := oauth2.NewClient(ctx, ts)

		service, err := people.NewService(ctx,
			option.WithHTTPClient(oauthClient),
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

		var u htsec.User
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
}
