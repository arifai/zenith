package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	locale "github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enlocale "github.com/go-playground/validator/v10/translations/en"
	"io"
	"reflect"
	"strings"
)

// IError represents a structured error with field and description.
type IError struct {
	Field string `json:"field"`
	Value string `json:"description"`
}

var (
	validate = validator.New()
	en       = locale.New()
	uni      = ut.New(en, en)
	trans, _ = uni.GetTranslator("en")
)

// SetupTranslation registers default English translations for the validator package.
// Logs to console if registration fails.
func SetupTranslation() {
	if err := enlocale.RegisterDefaultTranslations(validate, trans); err != nil {
		fmt.Printf("Error registering translation: %v", err)
	}
}

// ValidateBody validates the JSON body of an HTTP request and binds it to a struct.
func ValidateBody[T any](ctx *gin.Context) (*T, interface{}) {
	body := new(T)
	if err := ctx.ShouldBindJSON(body); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, []IError{}
		}
		return nil, []IError{{Value: CapitalizeFirstLetter(err.Error())}}
	}

	if errs := validateStruct(body); errs != nil {
		return nil, errs
	}

	return body, nil
}

// ValidateQuery binds query parameters to a struct and validates them.
// ctx is the Gin context containing the request data.
// Returns the parsed struct (body) or a slice of IError if validation fails.
func ValidateQuery[T any](ctx *gin.Context) (*T, interface{}) {
	body := new(T)
	if err := ctx.ShouldBindQuery(body); err != nil {
		return nil, []IError{{Value: CapitalizeFirstLetter(err.Error())}}
	}

	if errs := validateStruct(*body); errs != nil {
		return nil, errs
	}

	return body, nil

}

// validateStruct validates a given struct and returns a slice of IError if validation errormessage are present.
func validateStruct(body interface{}) []IError {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	var errs []IError
	if err := validate.Struct(body); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, err := range validationErrors {
				el := IError{
					Field: err.Field(),
					Value: CapitalizeFirstLetter(err.Translate(trans)),
				}
				errs = append(errs, el)
			}
		} else {
			errs = append(errs, IError{Value: CapitalizeFirstLetter(err.Error())})
		}
	}

	return errs
}
