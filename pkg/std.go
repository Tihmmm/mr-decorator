package pkg

import (
	"golang.org/x/crypto/ssh/terminal"
	"log"
)

func ReadSecretStdinToString(promptMessage string, dest *string) error {
	if promptMessage == "" {
		promptMessage = "Enter the secret: "
	}

	log.Println(promptMessage)
	secretBytes, err := terminal.ReadPassword(0)
	if err != nil {
		return err
	}

	*dest = string(secretBytes)
	return nil
}
