package servant

func NewSystem() *System {
	return &System{}
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct{}
