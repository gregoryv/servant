package htcook

import (
	"fmt"

	"golang.org/x/oauth2"
)

// Once authorized the session contains the information.
type Session struct {
	Token *oauth2.Token
	Name  string
	Email string
}

func (s *Session) String() string {
	return fmt.Sprintln(s.Name, s.Email)
}
