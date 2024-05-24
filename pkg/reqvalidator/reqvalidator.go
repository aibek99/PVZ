package reqvalidator

import (
	"log"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	validate            *validator.Validate
	phoneRegex          = `^(\+[1-9]\d{9,14}|\d{10,15})$`
	compiledPhoneRegexp *regexp.Regexp
)

func init() {
	validate = validator.New()

	compiledPhoneRegexp = regexp.MustCompile(phoneRegex)

	err := validate.RegisterValidation("phone", validatePhoneNumber)
	if err != nil {
		log.Printf("[reqvalidator][init] Unable to put validator for phonuNumber %v", err)
	}
}

// ValidateRequest is
func ValidateRequest(request interface{}) error {
	return validate.Struct(request)
}

func validatePhoneNumber(fl validator.FieldLevel) bool {
	phoneNumber := fl.Field().String()

	return compiledPhoneRegexp.MatchString(phoneNumber)
}
