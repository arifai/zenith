package common

import (
	"errors"
	"fmt"
	"github.com/arifai/zenith/internal/types/response"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type (
	// Response is a struct designed to handle various types of API responses, such as success, error, and authorization.
	Response struct{}

	// ResponseModel represents the structure of the response returned by the API.
	ResponseModel struct {
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Errors  []utils.IError `json:"errors"`
		Result  any            `json:"result"`
	}

	// EntriesModel represents a paginated collection of entries of a generic type T.
	EntriesModel[T interface{}] struct {
		Entries    []T `json:"entries"`
		Count      int `json:"count"`
		Page       int `json:"page"`
		TotalPages int `json:"total_pages"`
	}
)

// NewResponse initializes a new instance of the Response struct.
func NewResponse() *Response {
	return &Response{}
}

// New sets the response format and sends a JSON response with HTTP code, message, errormessage, and result data.
func (r Response) New(c *gin.Context, code int, message string, errors []utils.IError, result interface{}) {
	c.JSON(code, ResponseModel{
		Code:    code,
		Message: message,
		Errors:  errors,
		Result:  result,
	})
}

// NewEntries creates and returns a pointer to an EntriesModel with the given entries, count, page, and totalPages.
func NewEntries[T interface{}](entries []T, count, page, totalPages int) *EntriesModel[T] {
	return &EntriesModel[T]{
		Entries:    entries,
		Count:      count,
		Page:       page,
		TotalPages: totalPages,
	}
}

// Success sets a JSON success response with HTTP status 200 and the provided result data.
func (r Response) Success(c *gin.Context, result interface{}) {
	c.JSON(http.StatusOK, ResponseModel{
		Code:    http.StatusOK,
		Message: "Successful",
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Created sets a JSON response with HTTP status 201, providing a message and the result data.
func (r Response) Created(c *gin.Context, message string, result interface{}) {
	c.JSON(http.StatusCreated, ResponseModel{
		Code:    http.StatusCreated,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Authorized sends an HTTP 202 Accepted response with the given authentication result.
func (r Response) Authorized(c *gin.Context, result *response.AccountAuthResponse) {
	c.JSON(http.StatusAccepted, ResponseModel{
		Code:    http.StatusAccepted,
		Message: "Authorized",
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Unauthorized sends an HTTP 401 Unauthorized response with a custom message and a list of errormessage.
func (r Response) Unauthorized(c *gin.Context, errors []utils.IError, message string) {
	c.JSON(http.StatusUnauthorized, ResponseModel{
		Code:    http.StatusUnauthorized,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  errors,
		Result:  nil,
	})
}

// NotFound is a handler function that responds with a '404 Not Found' status and a formatted message using JSON.
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ResponseModel{
		Code:    http.StatusNotFound,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  nil,
	})
}

// Error handles different types of errormessage (string, []utils.IError, error) and responds with appropriate HTTP status.
func (r Response) Error(c *gin.Context, errParam interface{}) {
	switch err := errParam.(type) {
	case string:
		fmt.Printf("String error: %v\n", err)
		r.BadRequest(c, []utils.IError{}, err)
	case []utils.IError:
		fmt.Printf("Validation errormessage: %v\n", err)
		r.BadRequest(c, err, errormessage.ErrBadRequestText)
	case error:
		if errors.Is(err, io.EOF) {
			fmt.Printf("EOF error: %v\n", err)
			r.BadRequest(c, []utils.IError{}, errormessage.ErrRequestBodyEmptyText)
		} else {
			fmt.Printf("Error: %v\n", err)
			r.BadRequest(c, []utils.IError{}, err.Error())
		}
	default:
		fmt.Printf("Unhandled error type: %T, value: %v\n", err, err)
		r.InternalServerError(c, errormessage.ErrParsingRequestDataText)
	}
}

// BadRequest sends an HTTP 400 Bad Request response with a custom message and a list of errormessage.
func (r Response) BadRequest(c *gin.Context, errors []utils.IError, message string) {
	c.JSON(http.StatusBadRequest, ResponseModel{
		Code:    http.StatusBadRequest,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  errors,
		Result:  nil,
	})
}

// InternalServerError sends an HTTP 500 Internal Server Error response with a custom message and an empty list of errormessage.
func (r Response) InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ResponseModel{
		Code:    http.StatusInternalServerError,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  nil,
	})
}
