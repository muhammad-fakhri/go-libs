package httpmiddleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/muhammad-fakhri/go-libs/log"
)

// IngressLog represents concrete type of the middleware
type IngressLog struct {
	logger log.SLogger
	config *Config
}

type IngressLogger interface {
	Enforce(next http.Handler) http.Handler
	EnforceWithParams(next httprouter.Handle) httprouter.Handle
}

// LogMessage is a struct to keep the log message easier
type LogMessage struct {
	URL            string
	ReqMethod      string
	ReqHeader      http.Header
	ReqBody        string
	ResponseHeader http.Header
	ResponseCode   int
	ResponseBody   string
	TimeTakenInMS  int64
}

const (
	valueLogTypeIngress = "ingress_http"
)

type LogRequest struct {
	URL    string
	Method string
	Header http.Header
	Body   string
}

// NewIngressLogMiddleware is to initialize ingress log middleware object
func NewIngressLogMiddleware(logger log.SLogger, optionalConfig ...*Config) *IngressLog {
	var conf *Config
	if len(optionalConfig) == 0 || optionalConfig[0] == nil {
		conf = defaultConfig()
	} else {
		conf = NewConfig(optionalConfig[0])
	}

	return &IngressLog{
		logger: logger,
		config: conf,
	}
}

// Enforce is to apply log ingress middleware to the 'next' handler
func (i *IngressLog) Enforce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logReqMessage := buildLogRequest(r)

		newRequest := i.appendContextDataAndSetValue(r, i.logger)
		newWriter := NewResponseWriter(w)

		var (
			startTime       time.Time
			elapsedTimeInMS int64
		)

		defer func(ctx context.Context, reqmes *LogRequest, elapsedTimeInMS *int64, requestTimestamp *time.Time, writer *LogResponseWriter) {
			r := recover()
			if r != nil {
				fmt.Println("[ingress][panic] recovered from: ", r)
				debug.PrintStack()

				// default panic value
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(fmt.Sprintf("panic: %v.", r)))
			}

			i.log(newRequest.Context(), reqmes, *elapsedTimeInMS, *requestTimestamp, writer)

		}(newRequest.Context(), logReqMessage, &elapsedTimeInMS, &startTime, newWriter)

		startTime = time.Now()
		next.ServeHTTP(newWriter, newRequest)
		elapsedTimeInMS = time.Since(startTime).Milliseconds()

	})
}

func (i *IngressLog) log(ctx context.Context, request *LogRequest, timeTaken int64, requestTimestamp time.Time, rw *LogResponseWriter) {
	if i.config.DisableIngressLog || (i.config.LogFailedRequestOnly() && rw.statusCode == http.StatusOK) {
		// skip ingress log, rely on load balancer log or custom log instead
		return
	}

	data := &log.RequestResponse{
		Type:             valueLogTypeIngress,
		URLPath:          fmt.Sprintf("%s %s", request.Method, request.URL),
		RequestTimestamp: requestTimestamp,
		Status:           rw.statusCode,
		DurationMs:       timeTaken,
	}

	if i.config.LogResponseHeader() {
		header := rw.Header().Clone()
		header.Del("Authorization")
		data.ResponseHeader = header
	}

	if i.config.LogResponseBody() {
		if i.config.LogSuccessResponseBody() {
			data.ResponseBody = rw.bodyBuffer.String()
		} else {
			if rw.statusCode != http.StatusOK {
				data.ResponseBody = rw.bodyBuffer.String()
			} else {
				data.ResponseBody = wipedMessage
			}
		}
	}

	if i.config.LogRequestHeader() {
		header := request.Header.Clone()
		header.Del("Authorization")

		excludeRequestHeaderKeys := i.config.ExcludeOpt.RequestHeaderKeys
		if excludeRequestHeaderKeys != nil && len(excludeRequestHeaderKeys) > 0 {
			for _, headerKey := range excludeRequestHeaderKeys {
				header.Del(headerKey)
			}
		}

		data.RequestHeader = header
	}

	if i.config.LogRequestBody() {
		data.RequestBody = request.Body
	}

	i.logger.LogRequestResponse(ctx, data)
}

// Enforce is to apply log ingress middleware to the 'next' handler. Like http.HandlerFunc,
// but has a third parameter for the values of wildcards (variables), e.g: github.com/julienschmidt/httprouter
func (i *IngressLog) EnforceWithParams(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		logReqMessage := buildLogRequest(r)

		newRequest := i.appendContextDataAndSetValue(r, i.logger)
		newWriter := NewResponseWriter(w)

		var (
			startTime       time.Time
			elapsedTimeInMS int64
		)

		defer func(ctx context.Context, reqmes *LogRequest, elapsedTimeInMS *int64, requestTimestamp *time.Time, writer *LogResponseWriter) {
			r := recover()
			if r != nil {
				fmt.Println("[ingress][panic] recovered from: ", r)
				debug.PrintStack()

				// default panic value
				writer.WriteHeader(http.StatusInternalServerError)
				writer.Write([]byte(fmt.Sprintf("panic: %v.", r)))
			}

			i.log(newRequest.Context(), reqmes, *elapsedTimeInMS, *requestTimestamp, writer)

		}(newRequest.Context(), logReqMessage, &elapsedTimeInMS, &startTime, newWriter)

		startTime = time.Now()
		next(newWriter, newRequest, ps)
		elapsedTimeInMS = time.Since(startTime).Milliseconds()

	}
}

func (i *IngressLog) appendContextDataAndSetValue(r *http.Request, l log.SLogger) *http.Request {
	v := r.Context().Value(log.ContextDataMapKey)
	if v != nil {
		return r
	}

	data := &log.CommonFields{
		UserID: r.Header.Get(headerNameUserID),
	}

	if data.Country = strings.ToUpper(r.Header.Get(headerNameTenant)); data.Country == "" {
		data.Country = strings.ToUpper(r.Header.Get(headerNameCountry))
	}

	if data.ContextID = r.Header.Get(headerNameRequestID); data.ContextID == "" {
		data.ContextID = uuid.New().String()
	}

	url := r.URL.String()
	eventPrefix := i.config.GetEventPrefix()
	if strings.Contains(url, eventPrefix) {
		s := strings.Split(url, eventPrefix)
		if len(s) > 1 {
			data.EventID = strings.Split(s[1], URLSeparator)[0]
		}
	}

	return r.WithContext(l.SetContextData(r.Context(), data))
}

func buildLogRequest(r *http.Request) *LogRequest {
	return &LogRequest{
		URL:    r.URL.String(),
		Method: r.Method,
		Header: r.Header,
		Body:   getRequestBody(r),
	}
}

func getRequestBody(request *http.Request) string {
	if request.Body == nil {
		return "null"
	}

	requestBodyBytes, err := getBodyBytes(&request.Body)
	if err != nil {
		return "null"
	}

	return string(requestBodyBytes)
}

func getBodyBytes(body *io.ReadCloser) ([]byte, error) {
	responseBodyBytes, err := ioutil.ReadAll(*body)
	*body = ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes))
	return responseBodyBytes, err
}
