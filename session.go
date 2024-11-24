package servant

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gregoryv/servant/htsec"
	"golang.org/x/oauth2"
)

// Read more about security settings
// https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps#pattern-bff-cookie-security

func newCookie(t *oauth2.Token) *http.Cookie {
	return &http.Cookie{
		Name:     "token",
		Value:    t.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
	}
}

func sessionValid(r *http.Request) error {
	token, err := r.Cookie("token")
	if err != nil {
		return err
	}
	if _, found := sessions[token.Value]; !found {
		return fmt.Errorf("missing session")
	}
	return nil
}

func existingSession(r *http.Request) Session {
	ck, err := r.Cookie("token")
	if err != nil {
		return noSession
	}
	return sessions[ck.Value]
}

func newSession(t *oauth2.Token, u *htsec.User) {
	// cache the session
	session := Session{
		Token: t.AccessToken,
		Name:  u.Name,
		Email: u.Email,
	}
	sessions[t.AccessToken] = session
	debug.Println(session.String())
}

// token to name
var sessions = make(map[string]Session)

var noSession = Session{
	Name: "anonymous",
}

// Once authenticated the session contains the information.
type Session struct {
	Token string
	Name  string
	Email string
}

func (s *Session) String() string {
	return fmt.Sprintln(s.Name, s.Email)
}
