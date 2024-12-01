package servant

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gregoryv/htsec"
	"golang.org/x/oauth2"
)

// Read more about security settings
// https://datatracker.ietf.org/doc/html/draft-ietf-oauth-browser-based-apps#pattern-bff-cookie-security

func newCookie(state string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName, // todo __Host-
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(15 * time.Minute),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
}

func sessionValid(r *http.Request) error {
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

func newSession(slip *htsec.Slip) {
	// cache the session
	session := Session{
		Token: slip.Token,
		Name:  slip.Contact.Name,
		Email: slip.Contact.Email,
	}
	sessions[slip.State] = &session
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
