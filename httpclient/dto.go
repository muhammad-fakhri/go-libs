package httpclient

import (
	"net/http"
	"time"
)

// HTTPCallInfo store HTTP Call Information for logging purpose
type HTTPCallInfo struct {
	Request          *http.Request
	Response         *http.Response
	Duration         time.Duration
	RequestTimestamp time.Time
	HttpError        error
}
