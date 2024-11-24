package servant

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_newRouter_GET(t *testing.T) {
	// require
	h := newRouter()
	cases := []string{
		"/",
		"/inside",
		"/settings",
	}
	for _, path := range cases {
		t.Run("GET "+path, func(t *testing.T) {
			// require
			r := httptest.NewRequest("GET", path, http.NoBody)
			// do
			resp := recordResp(h, r)
			// ensure
			if err := statusCodeIs(resp.StatusCode, 200); err != nil {
				t.Error(err)
			}
		})
	}
}
