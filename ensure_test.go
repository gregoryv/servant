package servant

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func recordResp(h http.Handler, r *http.Request) *http.Response {
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w.Result()
}

func readBody(resp *http.Response) string {
	data, _ := ioutil.ReadAll(resp.Body)
	return string(data)
}

func ensure(t *testing.T, errors ...error) {
	t.Helper()
	for _, err := range errors {
		if err != nil {
			t.Error(err)
		}
	}
}

func statusCodeIs(got, exp int) error {
	if err := equals(got, exp); err != nil {
		return fmt.Errorf("status code %w", err)
	}
	return nil
}

func bodyContains(body string, phrases ...string) error {
	if err := contains(body, phrases...); err != nil {
		return fmt.Errorf("body %w", err)
	}
	return nil
}

func equals[T comparable](got, exp T) error {
	if got != exp {
		return fmt.Errorf("%v, expected %v", got, exp)
	}
	return nil
}

func contains(got string, expect ...string) error {
	var miss []string
	for _, exp := range expect {
		if !strings.Contains(got, exp) {
			miss = append(miss, exp)
		}
	}
	if len(miss) > 0 {
		return fmt.Errorf("missing %q", miss)
	}
	return nil
}
