package pkg

import (
	"crypto/rand"
	"github.com/alexedwards/argon2id"
	"log"
)

func GetArgonHash(secret string, params *argon2id.Params) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	if params == nil {
		params = argon2id.DefaultParams
	}

	return argon2id.CreateHash(secret, params)
}

func CheckArgonHash(inSecret, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(inSecret, hash)
	if err != nil {
		log.Printf("error comparing argon2 hash: %v\n", err)
		return false
	}

	return match
}
