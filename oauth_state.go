package servant

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// newState returns a string use.RANDOM.SIGNATURE using som private
func newState(use string) string {
	// see https://stackoverflow.com/questions/26132066/\
	//   what-is-the-purpose-of-the-state-parameter-in-oauth-authorization-request
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		debug.Println("state", err)
	}
	// both random value and the signature must be usable in a url
	random := hex.EncodeToString(randomBytes)
	signature := sign(random)
	return use + "." + random + "." + signature
}

func sign(random string) string {
	var privateKey = []byte("... load from file ...")
	hash := sha256.New()
	hash.Write([]byte(random))
	hash.Write(privateKey)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

// verify USE.RANDOM.SIGNATURE
func verify(state string) error {
	parts := strings.Split(state, ".")
	if len(parts) != 3 {
		return fmt.Errorf("state: invalid format")
	}
	signature := sign(parts[1])
	if signature != parts[2] {
		return fmt.Errorf("state: invalid signature")
	}
	return nil
}
