package servant

import "net/http"

func NewSystem() *System {
	var sys System
	return &sys
}

type System struct{}

func NewRouter(sys *System) http.HandlerFunc {
	return logware(
		authLayer(
			newRouter(),
		),
	)
}
