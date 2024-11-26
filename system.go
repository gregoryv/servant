package servant

import "github.com/gregoryv/servant/htsec"

func NewSystem() *System {
	return &System{}
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	security *htsec.Secure
}

func (sys *System) SetSecurity(v *htsec.Secure) { sys.security = v }
func (sys *System) Security() *htsec.Secure     { return sys.security }
