package servant

import (
	"github.com/gregoryv/servant/htsec"
	"github.com/gregoryv/servant/htsec/github"
	"github.com/gregoryv/servant/htsec/google"
)

func NewSystem() *System {
	return &System{
		security: htsec.NewDetail(
			github.Default(),
			google.Default(),
		),
	}
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	security *htsec.Detail
}

func (sys *System) SetSecurity(v *htsec.Detail) { sys.security = v }
func (sys *System) Security() *htsec.Detail     { return sys.security }
