package servant

import "github.com/gregoryv/servant/htauth"

func NewSystem() *System {
	return &System{}
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	security *htauth.Secure
}

func (sys *System) SetSecurity(v *htauth.Secure) { sys.security = v }
func (sys *System) Security() *htauth.Secure     { return sys.security }
