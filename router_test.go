package servant

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_NewRouter_GET(t *testing.T) {
	// require
	sys := NewSystem()

	h := NewRouter(sys)
	cases := map[string]int{
		// public
		"/": 200,

		// private
		"/inside":   200,
		"/settings": 200,
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
