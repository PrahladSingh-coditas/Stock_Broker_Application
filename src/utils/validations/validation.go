package validations

import (
	"fmt"
	"regexp"
	"stock_broker_application/src/constants"
	"stock_broker_application/src/models"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/go-playground/validator/v10"
)

var bffValidator *validator.Validate

func ValidatePasswordConstraints(password string) []string {
	var errors []string

	if len(password) < 8 {
		errors = append(errors, constants.ErrPasswordMinLength)
	}
	if !regexp.MustCompile(constants.LowercaseRegex).MatchString(password) {
		errors = append(errors, constants.ErrPasswordLowercase)
	}
	if !regexp.MustCompile(constants.UppercaseRegex).MatchString(password) {
		errors = append(errors, constants.ErrPasswordUppercase)
	}
	if !regexp.MustCompile(constants.DigitRegex).MatchString(password) {
		errors = append(errors, constants.ErrPasswordDigit)
	}
	if !regexp.MustCompile(constants.SpecialCharRegex).MatchString(password) {
		errors = append(errors, constants.ErrPasswordSpecialChar)
	}

	return errors
}

func FormatValidationErrors(err error) ([]models.ErrorMessage, string) {
	var validationErrors []models.ErrorMessage
	var validationErrorsStr string

	for _, err := range err.(validator.ValidationErrors) {
		var errorMsg string
		fieldName := err.Field()
		if err.Tag() == "required" {
			fieldName = strings.ToLower(fieldName)
			errorMsg = fmt.Sprintf(constants.ErrFieldRequired, fieldName)
		} else {
			switch err.Field() {
			case constants.FieldPassword:
				passwordErrors := ValidatePasswordConstraints(err.Value().(string))
				if len(passwordErrors) > 0 {
					for _, msg := range passwordErrors {
						validationErrors = append(validationErrors, models.ErrorMessage{
							Key:          err.Field(),
							ErrorMessage: msg,
							// ErrorMessage: msg + constants.ErrInvalidValue,
						})
					}
					continue
				}
			case constants.FieldConfirmPassword:
				if err.Tag() == "eqfield" {
					errorMsg = constants.ErrConfirmPasswordMatch
				}
			case constants.FieldPanCard:
				errorMsg = constants.ErrInvalidPanCard
			case constants.FieldPhoneNumber:
				errorMsg = constants.ErrInvalidPhoneNumber
			case constants.FieldEmail:
				errorMsg = constants.ErrInvalidEmail
			case constants.FieldOtp:
				switch err.Tag() {
				case "required":
					errorMsg = fmt.Sprintf(constants.ErrFieldRequired, "otp")
				case "otp":
					errorMsg = "otp must be exactly 4 digits and only numeric"
				default:
					errorMsg = fmt.Sprintf(constants.ErrInvalidValue, err.Field())
				}
			default:
				errorMsg = fmt.Sprintf(constants.ErrInvalidValue, err.Field())
			}
		}

		validationErrors = append(validationErrors, models.ErrorMessage{
			Key:          fieldName,
			ErrorMessage: errorMsg,
		})
		validationErrorsStr += fieldName + " is invalid; "
	}

	return validationErrors, validationErrorsStr
}

func panCardValidator(f1 validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(constants.PANCardRegex, f1.Field().String())
	return matched
}

func strongPasswordValidator(f1 validator.FieldLevel) bool {
	re := regexp2.MustCompile(constants.PasswordRegex, 0) // Compile regex with PCRE support
	matched, _ := re.MatchString(f1.Field().String())
	return matched
}

func IsEmailValid(f1 validator.FieldLevel) bool {
	email := f1.Field().String()

	// Split email into local part and domain part
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domainParts := strings.Split(parts[1], ".")

	if len(domainParts) < 2 || len(domainParts) > 3 {
		return false
	}

	for i := 0; i < len(domainParts)-1; i++ {
		for j := i + 1; j < len(domainParts); j++ {
			if domainParts[i] == domainParts[j] {
				return false
			}
		}
	}

	EmailRegex := regexp.MustCompile(constants.EmailRegex)

	return EmailRegex.MatchString(email)
}

func OtpValidator(f1 validator.FieldLevel) bool {
	matched, _ := regexp.MatchString(constants.OtpRegexp, f1.Field().String())
	return matched
}

func init() {
	bffValidator = validator.New()
	bffValidator.RegisterValidation("panCard", panCardValidator)
	bffValidator.RegisterValidation("strongPassword", strongPasswordValidator)
	bffValidator.RegisterValidation("Email", IsEmailValid)
	bffValidator.RegisterValidation("otp", OtpValidator)
}

func GetBFFValidator() *validator.Validate {
	return bffValidator
}
