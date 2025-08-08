package pkg

import (
	"log"

	"golang.org/x/term"
)

func ReadSecretStdinToString(promptMessage string, dest *string) error {
	if promptMessage == "" {
		promptMessage = "Enter the secret: "
	}

	log.Println(promptMessage)
	secretBytes, err := term.ReadPassword(0)
	if err != nil {
		return err
	}

	*dest = string(secretBytes)
	return nil
}
