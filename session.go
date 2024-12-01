package servant

import (
	"fmt"

	"golang.org/x/oauth2"
)

// Once authorized the session contains the information.
type Session struct {
	State string
	Token *oauth2.Token
	Name  string
	Email string
	dest  string
}

func (s *Session) String() string {
	return fmt.Sprintln(s.Name, s.Email)
}
