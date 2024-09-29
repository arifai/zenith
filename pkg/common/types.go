package common

import (
	"errors"
	"fmt"
	"github.com/arifai/zenith/internal/types/response"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"math"
	"net/http"
	"strings"
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
		Entries    []T   `json:"entries"`
		Count      int64 `json:"count"`
		Page       int   `json:"page"`
		TotalPages int   `json:"total_pages"`
	}

	Pagination struct {
		Offset int    `form:"offset" validate:"omitempty"`
		Limit  int    `form:"limit" validate:"omitempty"`
		Search string `form:"search" validate:"omitempty"`
		Sort   string `form:"sort" validate:"omitempty"`
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
func NewEntries[T interface{}](entries []T, count int64, page, totalPages int) *EntriesModel[T] {
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

// NotFound is a handler function that responds with a '404 Not Found' status and a formatted message using JSON.
func (r Response) NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ResponseModel{
		Code:    http.StatusNotFound,
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  nil,
	})
}

func (p Pagination) GetOffset() int {
	return (p.getOffset() - 1) * p.GetLimit()
}

func (p Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}

	return p.Limit
}

func (p Pagination) getOffset() int {
	if p.Offset == 0 {
		p.Offset = 1
	}

	return p.Offset
}

func (p Pagination) GetPage(count int64) int {
	var page int
	if count == 0 {
		page = 0
	} else {
		page = p.getOffset()
	}

	return page
}

func (p Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at asc"
	}

	return validateSort(p.Sort)
}

func (p Pagination) GetTotalPages(count int64) int {
	limit := p.GetLimit()
	if limit == 0 {
		return 0
	}

	return int(math.Ceil(float64(count) / float64(limit)))
}

// Paginate applies pagination, sorting, and search functionality to a database query using GORM.
// It takes a pointer to a Pagination struct that contains limit, offset, search, and sort values.
// Additionally, it accepts a `column` parameter which specifies the database column to be used for the search.
//
// The function ensures that only alphanumeric column names are allowed to prevent SQL injection risks.
//
// Important: This function must be used within a GORM Scope to properly apply pagination, search, and sorting to the query.
func Paginate(paging *Pagination, column string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if paging.Search != "" {
			searchTerm := "%" + strings.ToLower(paging.Search) + "%"
			lowerColumn := strings.ToLower(column)

			if isAlphaNumeric(lowerColumn) {
				db = db.Where(lowerColumn+" ILIKE ? ", searchTerm)
			}
		}

		return db.Offset(paging.GetOffset()).Limit(paging.GetLimit()).Order(paging.GetSort())
	}
}

func isAlphaNumeric(s string) bool {
	for _, char := range s {
		if !((char >= 'a' && char <= 'z') || (char >= '0' && char <= '9')) {
			return false
		}
	}
	return true
}

// validateSort sanitizes and validates the provided sort parameter to prevent SQL injection vulnerabilities (CWE-89).
//
// This function ensures that the input sort string only contains valid field names ("id" or "created_at")
// and valid sort directions ("asc" or "desc"). If the input is invalid or contains unexpected values,
// it defaults to "created_at asc" to maintain safety.
//
// Security:
//   - This function mitigates potential SQL injection attacks (CWE-89) by strictly validating the input
//     against allowed fields and directions, ensuring only safe and expected values are used in SQL queries.
//
// Reference: https://cwe.mitre.org/data/definitions/89.html
func validateSort(sort string) string {
	allowedFields := map[string]bool{
		"id":         true,
		"created_at": true,
	}

	allowedDirections := map[string]bool{
		"asc":  true,
		"desc": true,
	}

	parts := strings.Fields(strings.TrimSpace(sort))
	if len(parts) != 2 {
		return "created_at asc"
	}

	field, direction := parts[0], parts[1]
	if !allowedFields[field] || !allowedDirections[direction] {
		return "created_at asc"
	}

	return field + " " + direction
}
