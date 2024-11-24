package servant

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gregoryv/servant/htsec"
)

func Test_NewRouter_GET(t *testing.T) {
	// require
	sys := NewSystem()
	sec := htsec.NewSecure()
	h := NewRouter(sys, sec)
	cases := map[string]int{
		// public
		"/": 200,

		// private
		"/inside":   303,
		"/settings": 303,
	}
	for path, expCode := range cases {
		t.Run("GET "+path, func(t *testing.T) {
			// require
			r := httptest.NewRequest("GET", path, http.NoBody)
			// do
			resp := recordResp(h, r)
			// ensure
			err := statusCodeIs(resp.StatusCode, expCode)
			if err != nil {
				t.Error(err)
			}
		})
	}
}
