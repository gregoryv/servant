package servant

import (
	"github.com/gregoryv/servant/htsec"
	"github.com/gregoryv/servant/htsec/github"
	"github.com/gregoryv/servant/htsec/google"
)

func NewSystem() *System {
	return &System{
		security: htsec.NewGuard(
			github.Default(),
			google.Default(),
		),
	}
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	security *htsec.Guard
}

func (sys *System) SetSecurity(v *htsec.Guard) { sys.security = v }
func (sys *System) Security() *htsec.Guard     { return sys.security }
