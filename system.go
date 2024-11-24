package servant

import (
	"net/http"

	"github.com/gregoryv/servant/htsec"
)

func NewSystem() *System {
	var sys System
	return &sys
}

type System struct{}

func NewRouter(sys *System) http.HandlerFunc {
	sec := htsec.NewSecure()
	return logware(
		authLayer(
			sec,
			newRouter(),
		),
	)
}
