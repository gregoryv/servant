package servant

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

// Read more about security settings
// https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps#pattern-bff-cookie-security

func NewCookie(value string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName, // todo __Host-
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func SessionValid(r *http.Request) error {
	state, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}
	if _, found := sessions[state.Value]; !found {
		return fmt.Errorf("missing session")
	}
	return nil
}

const cookieName = "state"

func existingSession(r *http.Request) *Session {
	ck, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}
	return sessions[ck.Value]
}

func SetSession(key string, s *Session) {
	sessions[key] = s
	// todo save/restore sessions on restart
}

// state to session
var sessions = make(map[string]*Session)

// Once authorized the session contains the information.
type Session struct {
	Token *oauth2.Token
	Name  string
	Email string
}

func (s *Session) String() string {
	return fmt.Sprintln(s.Name, s.Email)
}
