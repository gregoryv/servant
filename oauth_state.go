package servant

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

// newState returns a string RANDOM.SIGNATURE using som private
func newState() string {
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
	return random + "." + signature
}

func sign(random string) string {
	var privateKey = []byte("... load from file ...")
	hash := sha256.New()
	hash.Write([]byte(random))
	hash.Write(privateKey)
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

func verify(state string) error {
	parts := strings.Split(state, ".")
	if len(parts) != 2 {
		return fmt.Errorf("state: missing dot")
	}
	signature := sign(parts[0])
	if signature != parts[1] {
		return fmt.Errorf("state: invalid signature")
	}
	return nil
}
