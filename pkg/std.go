package pkg

import (
	"golang.org/x/crypto/ssh/terminal"
	"log"
)

func ReadSecretStdinToString(dest *string) error {
	log.Println("Enter the secret: ")
	secretBytes, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}

	*dest = string(secretBytes)
	return nil
}
