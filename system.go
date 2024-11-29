package servant

import (
	"github.com/gregoryv/htsec"
	"github.com/gregoryv/htsec/github"
	"github.com/gregoryv/htsec/google"
)

func NewSystem() *System {
	return &System{
		security: htsec.NewDetail(
			github.Guard(),
			google.Guard(),
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
