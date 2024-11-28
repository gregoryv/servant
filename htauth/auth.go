package htauth

import "golang.org/x/oauth2"

type Auth struct {
	// Used for the oauth2 flow
	*oauth2.Config

	// Used to read user information once authenticated
	Contact ContactFunc
}

type ContactFunc = func(token *oauth2.Token) (*Contact, error)
