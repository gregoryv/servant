package htcook

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
