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

func (sys *System) Authorize(ctx context.Context, r *http.Request) (*Session, error) {
	slip, err := sys.sec.Authorize(ctx, r)
	if err != nil {
		return nil, err
	}
	s := Session{
		State: slip.State,
		Token: slip.Token,
		Name:  slip.Contact.Name,
		Email: slip.Contact.Email,
		dest:  slip.Dest(),
	}
	SetSession(slip.State, &s)
	return &s, err
}

func (sys *System) GuardURL(name, dest string) (string, error) {
	return sys.sec.GuardURL(name, dest)
}
