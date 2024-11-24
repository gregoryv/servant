package servant

import (
	"fmt"
	"net/http"
	"time"

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

func existingSession(r *http.Request) Session {
	ck, err := r.Cookie("token")
	if err != nil {
		return noSession
	}
	return sessions[ck.Value]
}

func newSession(t *oauth2.Token, u *user) {
	// cache the session
	session := Session{
		Token: t.AccessToken,
		Name:  u.Name,
		Email: u.Email,
	}
	sessions[t.AccessToken] = session
	debug.Println(session.String())
}

type user struct {
	Email string
	Name  string
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
