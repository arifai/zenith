package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Response is a struct that represent api response
type Response struct{}

// ResponseModel is a struct that represent api response model
type ResponseModel struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  any    `json:"result"`
}

// AuthorizeModel is a struct that represent authorize response
type AuthorizeModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// EntriesModel is a struct that represent entries response
type EntriesModel struct {
	Entries   []any `json:"entries"`
	Count     int   `json:"count"`
	Page      int   `json:"page"`
	TotalPage int   `json:"total_page"`
}

// New is a function to create new response
func (r Response) New(c *gin.Context, code int, message string, result any) {
	c.JSON(code, ResponseModel{
		Code:    code,
		Message: message,
		Result:  result,
	})
}

// Success is a function to create success response
func (r Response) Success(c *gin.Context, result any) {
	c.JSON(http.StatusOK, ResponseModel{
		Code:    http.StatusOK,
		Message: "",
		Result:  result,
	})
}

// Entries is a function to create entries response
func (r Response) Entries(c *gin.Context, entries EntriesModel) {
	r.Success(c, entries)
}

// Created is a function to create created response
func (r Response) Created(c *gin.Context, message string, result any) {
	c.JSON(http.StatusCreated, ResponseModel{
		Code:    http.StatusCreated,
		Message: message,
		Result:  result,
	})
}

// Authorized is a function to create authorized response
func (r Response) Authorized(c *gin.Context, result AuthorizeModel) {
	c.JSON(http.StatusAccepted, ResponseModel{
		Code:    http.StatusAccepted,
		Message: "Authorized",
		Result:  result,
	})
}

// BadRequest is a function to create bad request response
func (r Response) BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ResponseModel{
		Code:    http.StatusBadRequest,
		Message: message,
		Result:  nil,
	})
}

// InternalServerError is a function to create internal server error response
func (r Response) InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, ResponseModel{
		Code:    http.StatusInternalServerError,
		Message: message,
		Result:  nil,
	})
}
