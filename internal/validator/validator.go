package validator

import (
	"github.com/Tihmmm/mr-decorator/internal/models"
	jsonvalidator "github.com/go-playground/validator/v10"
	"log"
	"slices"
)

type Validator interface {
	IsValidAll(reqBody *models.MRRequest) bool
	IsValidStruct(reqBody *models.MRRequest) bool
	isValidArtifactFileName(fileName string) bool
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

func (v *RequestValidator) IsValidAll(reqBody *models.MRRequest) bool {
	return v.IsValidStruct(reqBody) && v.isValidArtifactFileName(reqBody.ArtifactFileName)
}

func (v *RequestValidator) IsValidStruct(reqBody *models.MRRequest) bool {
	err := v.jsonValidator.Struct(reqBody)
	if err == nil {
		return true
	}

	log.Printf("Error validating struct: %s\n", err)
	return false
}

func (v *RequestValidator) isValidArtifactFileName(fileName string) bool {
	if slices.Contains(v.validFilenames, fileName) {
		return true
	}

	log.Printf("Invalid artifact filename: %s\n", fileName)
	return false
}
