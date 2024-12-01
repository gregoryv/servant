package servant

import (
	"context"
	"net/http"

	"github.com/gregoryv/htsec"
	"github.com/gregoryv/htsec/github"
	"github.com/gregoryv/htsec/google"
)

func NewSystem() *System {
	s := &System{
		sec: htsec.NewDetail(
			github.Guard(),
			google.Guard(),
		),
	}
	s.sec.PrivateKey = []byte("my fixed private key")
	return s
}

// System carries domain logic which is exposed via a [http.Handler]
// using [NewRouter].
type System struct {
	sec *htsec.Detail
}

func (sys *System) SetSecurity(v *htsec.Detail) { sys.sec = v }
func (sys *System) Security() *htsec.Detail     { return sys.sec }

// todo return own Slip or maybe session
func (sys *System) Authorize(ctx context.Context, r *http.Request) (*htsec.Slip, error) {
	return sys.sec.Authorize(ctx, r)
}
