package htauth

import (
	"net/http"

	"golang.org/x/oauth2"
)

type AuthService struct {
	// Used for the oauth2 flow
	*oauth2.Config

	// Used to read user information once authenticated
	Contact ContactFunc
}

type ContactFunc = func(client *http.Client) (*Contact, error)
