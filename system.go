package servant

import "net/http"

func NewSystem() *System {
	var sys System
	return &sys
}

type System struct{}

func (s *System) Handler() http.HandlerFunc {
	return logware(
		authLayer(
			newRouter(),
		),
	)
}
