package servant

import "net/http"

func New() *System {
	sys := System{
		Handler: logware(
			AuthLayer(
				Endpoints(),
			),
		),
	}
	return &sys
}

type System struct {
	http.Handler
}
