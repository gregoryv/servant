package htauth

import (
	"crypto/rand"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

func NewSecure() *Secure {
	s := Secure{
		PrivateKey: make([]byte, 32),
		src:        make(map[string]*Auth),
	}
	_, _ = rand.Read(s.PrivateKey)
	return &s
}

type Secure struct {
	PrivateKey []byte
	src        map[string]*Auth
}

// Include authorization service. It's name will be the significant
// part of the AuthURL, e.g. for http://example.com/ the name will be
// example.
func (s *Secure) Include(a *Auth) {
	name := domainName(a.Endpoint.AuthURL)
	s.src[name] = a
}

func domainName(uri string) string {
	v, err := url.Parse(uri)
	if err != nil {
		panic(err.Error())
	}
	parts := strings.Split(v.Hostname(), ".")
	return parts[len(parts)-2]
}

// Names returns included authorization service names.
func (s *Secure) Names() []string {
	res := make([]string, 0, len(s.src))
	for name, _ := range s.src {
		res = append(res, name)
	}
	sort.Strings(res)
	return res
}

// AuthService returns named service if included, error if not found.
func (s *Secure) AuthService(name string) (*Auth, error) {
	a, found := s.src[name]
	if !found {
		err := fmt.Errorf("Secure.AuthService %v: %w", name, notFound)
		return nil, err
	}
	return a, nil
}

var notFound = fmt.Errorf("not found")

// ----------------------------------------
