package servant

import (
	"fmt"
	"net/http"
)

func existingSession(r *http.Request) Session {
	ck, err := r.Cookie("token")
	if err != nil {
		return noSession
	}
	return sessions[ck.Value]
}

// token to name
var sessions = make(map[string]Session)

var noSession = Session{
	Name: "anonymous",
}

// Once authenticated the session contains the information from
// github.
type Session struct {
	Token string
	Name  string
	Email string
}

func (s *Session) String() string {
	return fmt.Sprintln(s.Name, s.Email)
}
