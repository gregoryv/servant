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
func (s *Secure) NewState(use string) (string, error) {
	// see https://stackoverflow.com/questions/26132066/\
	//   what-is-the-purpose-of-the-state-parameter-in-oauth-authorization-request
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	// both random value and the signature must be usable in a url
	random := hex.EncodeToString(randomBytes)
	signature := s.Sign(random)
	return use + "." + random + "." + signature, nil
}

func (s *Secure) VerifyAndExchange(ctx context.Context, state, code string) (*oauth2.Token, *Contact, error) {
	if err := s.Verify(state); err != nil {
		return nil, nil, err
	}
	// which auth service was used
	auth, err := s.AuthService(s.ParseUse(state))
	if err != nil {
		return nil, nil, err
	}
	// get the token
	token, err := auth.Exchange(ctx, code)
	if err != nil {
		return nil, nil, err
	}

	// get user information from the Auth service
	contact, err := auth.Contact(token)
	if err != nil {
		return token, nil, err
	}
	return token, contact, err
}

func (s *Secure) Sign(random string) string {
	hash := sha256.New()
	hash.Write([]byte(random))
	hash.Write(s.PrivateKey)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// verify USE.RANDOM.SIGNATURE
func (s *Secure) Verify(state string) error {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return fmt.Errorf("state: invalid format")
	}
	signature := s.Sign(parts[1])
	if signature != parts[2] {
		return fmt.Errorf("state: invalid signature")
	}
	return nil
}

func (s *Secure) ParseUse(state string) string {
	i := strings.Index(state, ".")
	if i < 0 {
		return ""
	}
	return state[:i]
}
