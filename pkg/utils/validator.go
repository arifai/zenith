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

// IError is a struct to represent error
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

// SetupTranslation is a function to set up translation
func SetupTranslation() {
	if err := enlocale.RegisterDefaultTranslations(validate, trans); err != nil {
		fmt.Printf("Error registering translation: %v", err)
	}
}

// ValidateBody is a function to validate body
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

// ValidateQuery is a function to validate query
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

// validateStruct is a function to validate struct
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
