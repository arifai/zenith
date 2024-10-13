package common

import (
	"errors"
	"fmt"
	"github.com/arifai/zenith/internal/types/response"
	"github.com/arifai/zenith/pkg/errormessage"
	"github.com/arifai/zenith/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
		TraceID string         `json:"trace_id"`
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

	// Pagination struct handles fields required for paginating, searching, and sorting a collection of items.
	Pagination struct {
		Offset int    `form:"offset" validate:"omitempty"`
		Limit  int    `form:"limit" validate:"omitempty"`
		Search string `form:"search" validate:"omitempty"`
		Sort   string `form:"sort" validate:"omitempty"`
		Desc   bool   `form:"desc" validate:"omitempty"`
	}
)

// NewResponse initializes a new instance of the Response struct.
func NewResponse() *Response {
	return &Response{}
}

// extractTraceID extracts the trace ID from the context.
func (r Response) extractTraceID(c *gin.Context) string {
	span := trace.SpanFromContext(c.Request.Context())
	return span.SpanContext().TraceID().String()
}

// New sets the response format and sends a JSON response with HTTP code, message, errormessage, and result data.
func (r Response) New(c *gin.Context, code int, message string, errors []utils.IError, result interface{}) {
	c.JSON(code, ResponseModel{
		TraceID: r.extractTraceID(c),
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
		TraceID: r.extractTraceID(c),
		Message: "Successful",
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Created sets a JSON response with HTTP status 201, providing a message and the result data.
func (r Response) Created(c *gin.Context, message string, result interface{}) {
	c.JSON(http.StatusCreated, ResponseModel{
		TraceID: r.extractTraceID(c),
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Authorized sends an HTTP 202 Accepted response with the given authentication result.
func (r Response) Authorized(c *gin.Context, result *response.AccountAuthResponse) {
	c.JSON(http.StatusAccepted, ResponseModel{
		TraceID: r.extractTraceID(c),
		Message: "Authorized",
		Errors:  []utils.IError{},
		Result:  result,
	})
}

// Unauthorized sends an HTTP 401 Unauthorized response with a custom message and a list of errormessage.
func (r Response) Unauthorized(c *gin.Context, errors []utils.IError, message string) {
	c.JSON(http.StatusUnauthorized, ResponseModel{
		TraceID: r.extractTraceID(c),
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
		TraceID: r.extractTraceID(c),
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  errors,
		Result:  nil,
	})
}

// InternalServerError sends an HTTP 500 Internal Server Error response with a custom message and an empty list of errormessage.
func (r Response) InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ResponseModel{
		TraceID: r.extractTraceID(c),
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  nil,
	})
}

// NotFound is a handler function that responds with a '404 Not Found' status and a formatted message using JSON.
func (r Response) NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, ResponseModel{
		TraceID: r.extractTraceID(c),
		Message: utils.CapitalizeFirstLetter(message),
		Errors:  []utils.IError{},
		Result:  nil,
	})
}

// GetOffset returns the offset value for a paginated query.
// If the offset is not set, it defaults to 1. The returned value is calculated as (Offset - 1) * Limit.
func (p Pagination) GetOffset() int {
	if p.Offset == 0 {
		p.Offset = 1
	}

	return (p.Offset - 1) * p.GetLimit()
}

// GetLimit returns the limit value for a paginated query.
// If the limit is not set, it defaults to 10.
func (p Pagination) GetLimit() int {
	if p.Limit == 0 {
		p.Limit = 10
	}

	return p.Limit
}

// GetPage calculates and returns the current page number based on the given total item count.
func (p Pagination) GetPage(count int64) int {
	if count == 0 {
		return 0
	}

	limit := p.GetLimit()
	offset := p.GetOffset()

	return (offset / limit) + 1
}

// GetSort returns the sort field for a query. It defaults to "created_at" if no sort field is specified.
func (p Pagination) GetSort() string {
	if p.Sort == "" {
		p.Sort = "created_at"
	}

	return validateSort(p.Sort)
}

// GetTotalPages returns the total number of pages given the total item count.
func (p Pagination) GetTotalPages(count int64) int {
	limit := p.GetLimit()
	if limit == 0 {
		return 0
	}

	return int(math.Ceil(float64(count) / float64(limit)))
}

// Paginate applies pagination, sorting, and search functionality to a GORM database query.
// It uses values from the Pagination struct, including limit, offset, search, sort column, and a boolean (Desc) for sorting direction.
// The `column` parameter specifies the database column to be used for search functionality.
//
// Features:
// - Pagination: Applies limit and offset to control the number of records returned.
// - Sorting: Sorts results by the specified column and applies ascending or descending order based on the Desc flag.
// - Search: Filters results using the search term applied to the specified column.
//
// To mitigate SQL injection risks, only alphanumeric column names are allowed for sorting and searching.
//
// Important: This function should be used within a GORM Scope to apply pagination, sorting, and search functionality correctly.
func Paginate(paging *Pagination, column string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if paging.Search != "" {
			searchTerm := "%" + strings.ToLower(paging.Search) + "%"
			lowerColumn := strings.ToLower(column)

			if isAlphaNumeric(lowerColumn) {
				db = db.Where(lowerColumn+" ILIKE ? ", searchTerm)
			}
		}

		sort := paging.GetSort()

		return db.Offset(paging.GetOffset()).Limit(paging.GetLimit()).
			Order(clause.OrderByColumn{Column: clause.Column{Name: sort}, Desc: paging.Desc})
	}
}

// isAlphaNumeric checks if a given string consists only of alphanumeric characters (a-z, 0-9).
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
// This function ensures that the input sort string only contains valid field names (e.g., "id" or "created_at").
// If the input is invalid or contains unexpected values, the function defaults to "created_at" to ensure safety.
//
// Security:
//   - This function mitigates SQL injection attacks (CWE-89) by only allowing sorting on predefined fields
//     and rejecting any invalid input that could be used maliciously.
//
// Reference: https://cwe.mitre.org/data/definitions/89.html
func validateSort(sort string) string {
	allowedFields := map[string]bool{
		"id":         true,
		"created_at": true,
	}

	parts := strings.Fields(strings.TrimSpace(sort))
	if len(parts) != 2 {
		return "created_at"
	}

	field := parts[0]
	if !allowedFields[field] {
		return "created_at"
	}

	return field
}
