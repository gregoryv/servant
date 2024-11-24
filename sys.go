package servant

import "net/http"

func New() *System {
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
