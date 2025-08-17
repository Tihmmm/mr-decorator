package pkg

import (
	"log"

	"github.com/alexedwards/argon2id"
)

func GetArgonHash(secret string, params *argon2id.Params) (string, error) {
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
