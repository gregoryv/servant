package servant

import "net/http"

func New() http.Handler {
	return logware(
		AuthLayer(
			Endpoints(),
		),
	)
}
