package servant

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gregoryv/htsec"
	"github.com/gregoryv/htsec/github"
	"github.com/gregoryv/htsec/google"
)

func NewSystem() *System {
	s := &System{
		sec: htsec.NewSecurityDetail(
			github.Guard(),
			google.Guard(),
		),
		sessions: make(map[string]*Session),
	}
	s.sec.PrivateKey = []byte("my fixed private key")
	return s
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	sec *htsec.SecurityDetail

	sessions map[string]*Session
}

func (sys *System) Shutdown() (chan struct{}, func()) {
	c := make(chan struct{})
	return c, func() {
		debug.Print("########## shutting down")
		close(c)
	}
}

func (sys *System) Authorize(r *http.Request) (*Session, error) {
	slip, err := sys.sec.Authorize(r)
	if err != nil {
		return nil, err
	}
	s := Session{
		State: slip.State,
		Token: slip.Token,
		Name:  slip.Name,
		Email: slip.Email,
		dest:  slip.Destination(),
	}
	sys.SetSession(slip.State, &s)
	return &s, err
}

func (sys *System) GuardURL(name, dest string) (string, error) {
	return sys.sec.GuardURL(name, dest)
}

// ----------------------------------------

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

func (sys *System) SessionValid(r *http.Request) error {
	state, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}
	if _, found := sys.sessions[state.Value]; !found {
		return fmt.Errorf("missing session")
	}
	return nil
}

const cookieName = "state"

func (sys *System) ExistingSession(r *http.Request) *Session {
	ck, err := r.Cookie(cookieName)
	if err != nil {
		return nil
	}
	return sys.sessions[ck.Value]
}

func (sys *System) SetSession(key string, s *Session) {
	sys.sessions[key] = s
	// todo save/restore sessions on restart
}
