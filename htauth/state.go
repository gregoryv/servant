package htauth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"golang.org/x/oauth2"
)

// NewState returns a string use.RANDOM.SIGNATURE using som private
func (s *Guard) NewState(use string) (string, error) {
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

func (s *Guard) WhereIs(use string) (string, error) {
	svc, err := s.Gate(use)
	if err != nil {
		return "", err
	}
	state, err := s.NewState(use)
	if err != nil {
		return "", err
	}
	return svc.AuthCodeURL(state), nil
}
