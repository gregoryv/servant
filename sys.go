package servant

import "net/http"

func NewSystem() *System {
	sys := System{
		Handler: logware(
			authLayer(
				newRouter(),
			),
		),
	}
	return &sys
}

type System struct {
	http.Handler
}
