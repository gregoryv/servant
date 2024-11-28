package servant

import "github.com/gregoryv/servant/htauth"

func NewSystem() *System {
	return &System{}
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	security *htauth.Guard
}

func (sys *System) SetSecurity(v *htauth.Guard) { sys.security = v }
func (sys *System) Security() *htauth.Guard     { return sys.security }
