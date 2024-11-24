package htsec

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/oauth2"
)

func NewSecure() *Secure {
	s := Secure{
		src: map[string]*Auth{},
	}
	return &s
}

type Secure struct {
	src map[string]*Auth
}

func (s *Secure) Use(a *Auth) {
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

func (s *Secure) Names() []string {
	res := make([]string, 0, len(s.src))
	for name, _ := range s.src {
		res = append(res, name)
	}
	sort.Strings(res)
	return res
}

func (s *Secure) AuthService(name string) (*Auth, error) {
	a, found := s.src[name]
	if !found {
		return nil, fmt.Errorf("AuthService %v: %w", name, notFound)
	}
	return a, nil
}

var notFound = fmt.Errorf("not found")

// ----------------------------------------

type Auth struct {
	*oauth2.Config
	ReadUser ReadUserFunc
}

type User struct {
	Email string
	Name  string
}

type ReadUserFunc = func(accessToken string) (*User, error)
