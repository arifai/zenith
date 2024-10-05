package utils

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
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
	Value string `json:"reason"`
}

var (
	validate = validator.New()
	uni      = ut.New(en.New(), en.New())
	trans, _ = uni.GetTranslator("en")
)

// SetupTranslation registers default English translations for the validator package.
func SetupTranslation() {
	if err := enlocale.RegisterDefaultTranslations(validate, trans); err != nil {
		fmt.Printf("Error registering translation: %v", err)
	}
}

// ValidateBody validates the JSON body of an HTTP request and binds it to a struct.
func ValidateBody[T any](ctx *gin.Context) (*T, interface{}) {
	body := new(T)
	if err := ctx.ShouldBind(body); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, []IError{}
		}
		return nil, []IError{{Value: err.Error()}}
	}

	if errs := validateStruct(body); errs != nil {
		return nil, errs
	}

	return body, nil
}

// ValidateQuery binds query parameters to a struct and validates them.
func ValidateQuery[T any](ctx *gin.Context) (*T, interface{}) {
	body := new(T)
	if err := ctx.ShouldBindQuery(body); err != nil {
		return nil, []IError{{Value: err.Error()}}
	}

	if errs := validateStruct(body); errs != nil {
		return nil, errs
	}

	return body, nil
}

// validateStruct validates a given struct and returns a slice of IError if validation errors are present.
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
				errMsg := err.Translate(trans)

				if field, found := reflect.TypeOf(body).Elem().FieldByName(err.StructField()); found {
					customMsg := getCustomReason(field, err.Tag())
					if customMsg != "" {
						errMsg = customMsg
					}
				}

				errs = append(errs, IError{
					Field: err.Field(),
					Value: errMsg,
				})
			}
		} else {
			errs = append(errs, IError{Value: err.Error()})
		}
	}

	return errs
}

// getCustomReason retrieves a custom reason based on the validation tag from the struct field's "reason" tag.
func getCustomReason(field reflect.StructField, validationTag string) string {
	reasonTag := field.Tag.Get("reason")
	rules := strings.Split(reasonTag, ";")

	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) == 2 && parts[0] == validationTag {
			return parts[1]
		}
	}

	return ""
}
