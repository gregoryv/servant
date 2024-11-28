package htsec

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"

	"golang.org/x/oauth2"
)

func NewGuard() *Guard {
	s := Guard{
		PrivateKey: make([]byte, 32),
		src:        make(map[string]*Gate),
	}
	_, _ = rand.Read(s.PrivateKey)
	return &s
}

type Guard struct {
	PrivateKey []byte
	src        map[string]*Gate
}

// Include authorization service. It's name will be the significant
// part of the AuthURL, e.g. for http://example.com/ the name will be
// example.
func (s *Guard) Include(gates ...*Gate) {
	for _, gate := range gates {
		name := domainName(gate.Endpoint.AuthURL)
		s.src[name] = gate
	}
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
func (s *Guard) Names() []string {
	res := make([]string, 0, len(s.src))
	for name, _ := range s.src {
		res = append(res, name)
	}
	sort.Strings(res)
	return res
}

func (s *Guard) WhereIs(use string) (string, error) {
	svc, err := s.Gate(use)
	if err != nil {
		return "", err
	}
	state, err := s.newState(use)
	if err != nil {
		return "", err
	}
	return svc.AuthCodeURL(state), nil
}

// AuthService returns named service if included, error if not found.
func (s *Guard) Gate(name string) (*Gate, error) {
	a, found := s.src[name]
	if !found {
		err := fmt.Errorf("Secure.AuthService %v: %w", name, notFound)
		return nil, err
	}
	return a, nil
}

var notFound = fmt.Errorf("not found")

// NewState returns a string use.RANDOM.SIGNATURE using som private
func (s *Guard) newState(use string) (string, error) {
	// see https://stackoverflow.com/questions/26132066/\
	//   what-is-the-purpose-of-the-state-parameter-in-oauth-authorization-request
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// both random value and the signature must be usable in a url
	random := hex.EncodeToString(randomBytes)
	signature := s.sign(random)
	return use + "." + random + "." + signature, nil
}

func (s *Guard) Authorize(ctx context.Context, state, code string) (*oauth2.Token, *Contact, error) {
	if err := s.verify(state); err != nil {
		return nil, nil, err
	}
	// which auth service was used
	auth, err := s.Gate(s.parseUse(state))
	if err != nil {
		return nil, nil, err
	}
	// get the token
	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return nil, nil, err
	}

	client := auth.Client(ctx, token)
	// get user information from the Auth service
	contact, err := auth.Contact(client)
	if err != nil {
		return token, nil, err
	}
	return token, contact, err
}

// verify USE.RANDOM.SIGNATURE
func (s *Guard) verify(state string) error {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return fmt.Errorf("state: invalid format")
	}
	signature := s.sign(parts[1])
	if signature != parts[2] {
		return fmt.Errorf("state: invalid signature")
	}
	return nil
}

func (s *Guard) sign(random string) string {
	hash := sha256.New()
	hash.Write([]byte(random))
	hash.Write(s.PrivateKey)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func (s *Guard) parseUse(state string) string {
	i := strings.Index(state, ".")
	if i < 0 {
		return ""
	}
	return state[:i]
}
