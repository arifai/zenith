package common

import (
	"errors"
	"fmt"
	errmsg "github.com/arifai/go-modular-monolithic/internal/errors"
	"github.com/arifai/go-modular-monolithic/pkg/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type (
	// Response is a struct that represent api response
	Response struct{}

	// ResponseModel is a struct that represent api response model
	ResponseModel struct {
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Errors  []utils.IError `json:"errors"`
		Result  any            `json:"result"`
	}

	// AuthResponse is a struct that represent authorize response
	AuthResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	// EntriesModel is a struct that represent entries response
	EntriesModel[T interface{}] struct {
		Entries    []T `json:"entries"`
		Count      int `json:"count"`
		Page       int `json:"page"`
		TotalPages int `json:"total_pages"`
	}
)

// New is a function to create new response. You can customize the response code, message, and result
func (r Response) New(c *gin.Context, code int, message string, errors []utils.IError, result interface{}) {
	c.JSON(code, ResponseModel{
		Code:    code,
		Message: message,
		Errors:  errors,
		Result:  result,
	})
}

// NewEntries is a function to create new entries response
func NewEntries[T interface{}](entries []T, count, page, totalPages int) *EntriesModel[T] {
	return &EntriesModel[T]{
		Entries:    entries,
		Count:      count,
		Page:       page,
		TotalPages: totalPages,
	}
}

// Success is a function to create success response
func (r Response) Success(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, ResponseModel{
		Code:    http.StatusOK,
		Message: "Successful",
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Created is a function to create created response
func (r Response) Created(c *gin.Context, message string, result interface{}) {
	c.JSON(http.StatusCreated, ResponseModel{
		Code:    http.StatusCreated,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Authorized is a function to create authorized response
func (r Response) Authorized(c *gin.Context, result *AuthResponse) {
	c.JSON(http.StatusAccepted, ResponseModel{
		Code:    http.StatusAccepted,
		Message: "Authorized",
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Unauthorized is a function to create unauthorized response
func (r Response) Unauthorized(c *gin.Context, errors []utils.IError, message string) {
	c.JSON(http.StatusUnauthorized, ResponseModel{
		Code:    http.StatusUnauthorized,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  errors,
		Result:  nil,
	})
}

// NotFound is a function to create not found response
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ResponseModel{
		Code:    http.StatusNotFound,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  nil,
	})
}

// Error is a function to create error response
func (r Response) Error(c *gin.Context, errParam interface{}) {
	switch err := errParam.(type) {
	case string:
		fmt.Printf("String error: %v\n", err)
		r.BadRequest(c, []utils.IError{}, err)
	case []utils.IError:
		fmt.Printf("Validation errors: %v\n", err)
		r.BadRequest(c, err, errmsg.ErrBadRequestText)
	case error:
		if errors.Is(err, io.EOF) {
			fmt.Printf("EOF error: %v\n", err)
			r.BadRequest(c, []utils.IError{}, errmsg.ErrRequestBodyEmptyText)
		} else {
			fmt.Printf("Error: %v\n", err)
			r.BadRequest(c, []utils.IError{}, err.Error())
		}
	default:
		fmt.Printf("Unhandled error type: %T, value: %v\n", err, err)
		r.InternalServerError(c, errmsg.ErrParsingRequestDataText)
	}
}

// BadRequest is a function to create bad request response
func (r Response) BadRequest(c *gin.Context, errors []utils.IError, message string) {
	c.JSON(http.StatusBadRequest, ResponseModel{
		Code:    http.StatusBadRequest,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  errors,
		Result:  nil,
	})
}

// InternalServerError is a function to create internal server error response
func (r Response) InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ResponseModel{
		Code:    http.StatusInternalServerError,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  nil,
	})
}
