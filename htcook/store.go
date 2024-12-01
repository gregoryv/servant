package htcook

import (
	"fmt"
	"net/http"
)

func NewStore() *Store {
	return &Store{
		sessions: make(map[string]*Session),
	}
}

type Store struct {
	sessions map[string]*Session
}

func (t *Store) SetSession(key string, s *Session) {
	t.sessions[key] = s
}

func (t *Store) SessionValid(r *http.Request) error {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return err
	}
	if _, found := t.sessions[c.Value]; !found {
		return fmt.Errorf("missing session")
	}
	return nil
}

const cookieName = "state"
