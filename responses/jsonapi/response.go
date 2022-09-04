/*
Package jsonapi implements JSON API Response to maintain standard response for each API.

The usage of this package is :

1. Return Success with Data interface
return jsonapi.Success().WithData(newCampaign)

2. Return Error with validation message based on Beego validation package
return jsonapi.BadRequest().WithValidationMessage(ValidationResult)

3. Return Error with some message
return jsonapi.InternalError().WithMessage("weeeeew")

4. Get status code only
jsonapi.InternalError().StatusCode
*/

package jsonapi

import (
	"fmt"
	"net/http"

	beegoContext "github.com/astaxie/beego/context"
)

type Responder interface {
	Success() *JsonAPIResponse
	Error() *JsonAPIResponse
	BadRequest() *JsonAPIResponse
	Unauthorized() *JsonAPIResponse
	Forbidden() *JsonAPIResponse
	NotFound() *JsonAPIResponse
	InternalError() *JsonAPIResponse
	AlreadyExist() *JsonAPIResponse
	CustomStatusCode(statusCode int) *JsonAPIResponse
	WithErrorMessage(message string) *JsonAPIResponse
	WithData(Data interface{}) *JsonAPIResponse
	WithValidationMessage(Result []ValidationMessage) *JsonAPIResponse
	WithErrorList(errors []string) *JsonAPIResponse
	WithErrorCodeAndMessage(message string, errCode ErrorCode) *JsonAPIResponse
	WithMultipleErrorCodeAndMessage(errors ...ErrorDetail) *JsonAPIResponse
}

type JsonAPIResponse struct {
	StatusCode int           `json:"-"`
	Data       interface{}   `json:"data,omitempty"`
	ErrorList  []ErrorDetail `json:"errors,omitempty"`
}

func NewJSONAPIResponse() Responder {
	return &JsonAPIResponse{StatusCode: 0, Data: nil, ErrorList: nil}
}

type ErrorDetail struct {
	Status string    `json:"status,omitempty"`
	Code   ErrorCode `json:"code,omitempty"`
	Source string    `json:"source,omitempty"`
	Title  string    `json:"title,omitempty"`
	Detail string    `json:"detail"`
}

func (e *ErrorDetail) Error() string {
	return fmt.Sprintf("Error: %s %s\n", e.Title, e.Detail)
}

type ValidationMessage struct {
	Message string
	Code    ErrorCode
}

func newResponseCode(statusCode int) *JsonAPIResponse {
	return &JsonAPIResponse{
		StatusCode: statusCode,
	}
}

// Will be deprecated
func (r *JsonAPIResponse) WithErrorMessage(message string) *JsonAPIResponse {
	r.ErrorList = append(r.ErrorList, ErrorDetail{Detail: message})
	return r
}

// Will be deprecated
func (r *JsonAPIResponse) WithErrorList(errors []string) *JsonAPIResponse {
	for _, element := range errors {
		r.ErrorList = append(r.ErrorList, ErrorDetail{Detail: element})
	}
	return r
}

func (r *JsonAPIResponse) WithData(Data interface{}) *JsonAPIResponse {
	r.Data = Data
	return r
}

func (r *JsonAPIResponse) WithValidationMessage(Result []ValidationMessage) *JsonAPIResponse {
	for _, elem := range Result {
		r.ErrorList = append(r.ErrorList, ErrorDetail{Detail: elem.Message, Code: elem.Code})
	}
	return r
}

func (r *JsonAPIResponse) Success() *JsonAPIResponse {
	return newResponseCode(http.StatusOK)
}

func (r *JsonAPIResponse) Error() *JsonAPIResponse {
	return newResponseCode(http.StatusOK)
}

func (r *JsonAPIResponse) BadRequest() *JsonAPIResponse {
	return newResponseCode(http.StatusBadRequest)
}

func (r *JsonAPIResponse) Unauthorized() *JsonAPIResponse {
	return newResponseCode(http.StatusUnauthorized)
}

func (r *JsonAPIResponse) Forbidden() *JsonAPIResponse {
	return newResponseCode(http.StatusForbidden)
}

func (r *JsonAPIResponse) NotFound() *JsonAPIResponse {
	return newResponseCode(http.StatusNotFound)
}

func (r *JsonAPIResponse) InternalError() *JsonAPIResponse {
	return newResponseCode(http.StatusInternalServerError)
}

func (r *JsonAPIResponse) AlreadyExist() *JsonAPIResponse {
	return newResponseCode(http.StatusConflict)
}

func (r *JsonAPIResponse) CustomStatusCode(statusCode int) *JsonAPIResponse {
	return newResponseCode(statusCode)
}

func (r *JsonAPIResponse) WithErrorCodeAndMessage(message string, errCode ErrorCode) *JsonAPIResponse {
	r.ErrorList = append(r.ErrorList, ErrorDetail{Detail: message, Code: errCode})
	return r
}

func (r *JsonAPIResponse) WithErrorCodeAndMessageAndTitle(message, title string, errCode ErrorCode) *JsonAPIResponse {
	r.ErrorList = append(r.ErrorList, ErrorDetail{Detail: message, Code: errCode, Title: title})
	return r
}

func (r *JsonAPIResponse) WithMultipleErrorCodeAndMessage(errs ...ErrorDetail) *JsonAPIResponse {
	r.ErrorList = append(r.ErrorList, errs...)
	return r
}

func (r *JsonAPIResponse) Render(ctx *beegoContext.Context) {
	ctx.Output.SetStatus(r.StatusCode)
	ctx.Output.JSON(JsonAPIResponse{0, r.Data, r.ErrorList}, false, false)
}
