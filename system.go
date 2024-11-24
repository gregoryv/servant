package servant

import (
	"net/http"

	"github.com/gregoryv/servant/htsec"
)

func NewSystem() *System {
	return &System{}
}

type System struct {
	*htsec.Secure
}

func NewRouter(sys *System, sec *htsec.Secure) http.HandlerFunc {
	return logware(
		authLayer(
			sec,
			newRouter(),
		),
	)
}
