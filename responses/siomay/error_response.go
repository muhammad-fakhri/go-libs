package responses

import (
	"context"
	beegoContext "github.com/astaxie/beego/context"
	"net/http"
)

type ErrorResponse struct {
	Context    context.Context
	Message    string
	StatusCode int
}

func newErrorResponse(statusCode int) ErrorResponse {
	return ErrorResponse{
		Context:    context.TODO(),
		StatusCode: statusCode,
		Message:    "",
	}
}

func (r ErrorResponse) WithContext(ctx context.Context) ErrorResponse {
	r.Context = ctx
	return r
}

func (r ErrorResponse) WithMessage(message string) ErrorResponse {
	r.Message = message
	return r
}

func BadRequest() ErrorResponse {
	return newErrorResponse(http.StatusBadRequest)
}

func Unauthorized() ErrorResponse {
	return newErrorResponse(http.StatusUnauthorized)
}

func Forbidden() ErrorResponse {
	return newErrorResponse(http.StatusForbidden)
}

func NotFound() ErrorResponse {
	return newErrorResponse(http.StatusNotFound)
}

func InternalError() ErrorResponse {
	return newErrorResponse(http.StatusInternalServerError)
}

func AlreadyExist() ErrorResponse {
	return newErrorResponse(http.StatusConflict)
}

func (r ErrorResponse) Error() string {
	return r.Message
}

type errorMessage struct {
	Message string `json:"message,omitempty"`
}

func (r ErrorResponse) Render(ctx *beegoContext.Context) {
	ctx.Output.SetStatus(r.StatusCode)

	clientMessage := ""
	if r.StatusCode == http.StatusBadRequest || r.StatusCode == http.StatusUnauthorized || r.StatusCode == http.StatusForbidden ||
		r.StatusCode == http.StatusNotFound || r.StatusCode == http.StatusInternalServerError || r.StatusCode == http.StatusConflict {
		clientMessage = r.Message
	}
	ctx.Output.JSON(errorMessage{Message: clientMessage}, false, false)
}
