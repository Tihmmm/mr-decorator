package validator

import (
	"github.com/Tihmmm/mr-decorator/internal/models"
	jsonvalidator "github.com/go-playground/validator/v10"
	"log"
	"slices"
)

type Validator interface {
	Validate(reqBody *models.MRRequest) bool
}

type RequestValidator struct {
	jsonValidator  *jsonvalidator.Validate
	validFilenames []string
}

func NewValidator() Validator {
	return &RequestValidator{
		jsonValidator:  jsonvalidator.New(),
		validFilenames: []string{models.FprFn, models.CyclonedxJsonFn, models.DependencyCheckJsonFn},
	}
}

func (v *RequestValidator) Validate(reqBody *models.MRRequest) bool {
	return v.ValidateStruct(reqBody) && v.ValidateArtifactFileName(reqBody.ArtifactFileName)
}

func (v *RequestValidator) ValidateStruct(reqBody *models.MRRequest) bool {
	err := v.jsonValidator.Struct(reqBody)
	if err == nil {
		return true
	}

	log.Printf("Error validating request body: %s\n", err)
	return false
}

func (v *RequestValidator) ValidateArtifactFileName(fileName string) bool {
	if slices.Contains(v.validFilenames, fileName) {
		return true
	}

	log.Printf("Invalid artifact filename: %s\n", fileName)
	return false
}
