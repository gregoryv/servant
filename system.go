package servant

import "github.com/gregoryv/servant/htsec"

func NewSystem() *System {
	return &System{}
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	security *htsec.Guard
}

func (sys *System) SetSecurity(v *htsec.Guard) { sys.security = v }
func (sys *System) Security() *htsec.Guard     { return sys.security }
