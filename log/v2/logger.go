// naming package is logrus because in SetReportCaller will set caller to first function outside logrus package
package log

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type SLogger interface {
	SetLevel(level log.Level)

	BuildContextDataAndSetValue(country string, contextId string) (ctx context.Context)
	AppendContextDataAndSetValue(r *http.Request, country string, contextId string) *http.Request
	SetContextDataAndSetValue(r *http.Request, data map[string]string, country string, contextId string) *http.Request

	CreateResponseWrapper(rw http.ResponseWriter) *loggingResponseWriter

	GetEntry() *log.Entry

	Infof(ctx context.Context, message string, args ...interface{})
	Errorf(ctx context.Context, message string, args ...interface{})
	Warnf(ctx context.Context, message string, args ...interface{})
	Debugf(ctx context.Context, message string, args ...interface{})
	Fatalf(ctx context.Context, message string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Debug(ctx context.Context, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})

	LogRequest(ctx context.Context, r *http.Request)
	LogResponse(ctx context.Context, rw *loggingResponseWriter)
}

// safe typing https://golang.org/pkg/context/#WithValue
type contextDataMapKeyType string

// add key here for future request based value
var (
	// data map to contain values
	ContextDataMapKey contextDataMapKeyType = "value"

	// context key data added to map
	ContextCountryKey = "country"
	ContextIdKey      = "context_id"
	PathKey           = "url_path"
	RequestKey        = "request"
	ResponseKey       = "response"
	ResponseCodeKey   = "response_code"
)

type SLog struct {
	entry *log.Entry
}

type LogParams struct {
	fields log.Fields
}

// context key data added to map
type contextData struct {
	country   string
	contextId string
}

const (
	maximumCallerDepth int = 25
	knownLogFrames     int = 3
)

var (
	minimumCallerDepth = 1

	// Used for caller information initialisation
	callerInitOnce sync.Once

	// qualified package name, cached at first use
	thisPackageName string
)

func NewSLogger(service string) SLogger {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	entry := log.NewEntry(logger)
	entry = entry.WithField("service", service)
	return &SLog{entry}
}

func NewSLoggerWithLevel(service string, level log.Level) SLogger {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})
	logger.SetLevel(level)
	entry := log.NewEntry(logger)
	entry = entry.WithField("service", service)
	return &SLog{entry}
}

func (l *SLog) SetLevel(level log.Level) {
	l.entry.Logger.SetLevel(level)
}

func (l *SLog) getContextData(ctx context.Context) *contextData {
	dataMap := ctx.Value(ContextDataMapKey)
	var result *contextData

	if dataMap != nil {
		if data, ok := dataMap.(map[string]string); ok {
			result = &contextData{
				country:   data[ContextCountryKey],
				contextId: data[ContextIdKey],
			}
		}
	}

	return result
}

func (l *SLog) BuildContextDataAndSetValue(country string, contextId string) (ctx context.Context) {
	data := make(map[string]string, 0)
	data[ContextCountryKey] = country
	data[ContextIdKey] = contextId

	ctx = context.WithValue(context.Background(), ContextDataMapKey, data)

	return ctx
}

func (l *SLog) AppendContextDataAndSetValue(r *http.Request, country string, contextId string) *http.Request {
	data := make(map[string]string, 0)
	data[ContextCountryKey] = country
	data[ContextIdKey] = contextId

	ctx := context.WithValue(r.Context(), ContextDataMapKey, data)

	return r.WithContext(ctx)
}

func (l *SLog) SetContextDataAndSetValue(r *http.Request, data map[string]string, country string, contextId string) *http.Request {
	if data == nil {
		data = make(map[string]string, 0)
	}
	data[ContextCountryKey] = country
	data[ContextIdKey] = contextId

	ctx := context.WithValue(r.Context(), ContextDataMapKey, data)

	return r.WithContext(ctx)
}

func (l *SLog) CreateResponseWrapper(rw http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: rw,
	}
}

func (l *SLog) GetEntry() *log.Entry {
	return l.entry
}

func (l *SLog) Infof(ctx context.Context, message string, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.InfoLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Infof(message, args...)
}

func (l *SLog) Warnf(ctx context.Context, message string, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.WarnLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Warningf(message, args...)
}

func (l *SLog) Errorf(ctx context.Context, message string, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.ErrorLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Errorf(message, args...)
}

func (l *SLog) Debugf(ctx context.Context, message string, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.DebugLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Debugf(message, args...)
}

func (l *SLog) Fatalf(ctx context.Context, message string, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.FatalLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Fatalf(message, args...)
}

func (l *SLog) Info(ctx context.Context, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.InfoLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Info(args...)
}

func (l *SLog) Warn(ctx context.Context, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.WarnLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Warning(args...)
}

func (l *SLog) Error(ctx context.Context, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.ErrorLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Error(args...)
}

func (l *SLog) Debug(ctx context.Context, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.DebugLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Debug(args...)
}

func (l *SLog) Fatal(ctx context.Context, args ...interface{}) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.FatalLevel)
	lp.injectContextDataMap(ctx)
	l.entry.WithFields(lp.fields).Fatal(args...)
}

func (l *SLog) LogRequest(ctx context.Context, r *http.Request) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.InfoLevel)
	lp.injectContextDataMap(ctx).injectURLPath(ctx, r).injectRequestBody(ctx, r)
	l.entry.WithFields(lp.fields).Info("Request Body")
}

func (l *SLog) LogResponse(ctx context.Context, rw *loggingResponseWriter) {
	lp := LogParams{fields: log.Fields{}}
	lp.setCallStackTrace(log.InfoLevel)
	lp.injectContextDataMap(ctx).injectResponseBody(ctx, rw)
	l.entry.WithFields(lp.fields).Info("Response Body")
}

func (lp *LogParams) setCallStackTrace(logLevel log.Level) {
	if logLevel <= log.ErrorLevel {
		lp.setCaller(getCaller())
	}
}

func (lp *LogParams) setCaller(caller *runtime.Frame) {
	if caller == nil {
		return
	}

	funcVal := caller.Function
	fileVal := fmt.Sprintf("%s:%d", caller.File, caller.Line)
	if funcVal != "" {
		lp.fields["func"] = funcVal
	}
	if fileVal != "" {
		lp.fields["file"] = fileVal
	}
}

func (lp *LogParams) injectContextDataMap(ctx context.Context) *LogParams {
	dataMap := ctx.Value(ContextDataMapKey)

	if dataMap != nil {
		if data, ok := dataMap.(map[string]string); ok {
			for key, value := range data {
				lp.fields[key] = value
			}
		}
	}

	return lp
}

func (lp *LogParams) injectURLPath(ctx context.Context, r *http.Request) *LogParams {
	lp.fields[PathKey] = r.Host + r.URL.Path
	return lp
}

func (lp *LogParams) injectRequestBody(ctx context.Context, r *http.Request) *LogParams {
	buf, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	lp.fields[RequestKey] = fmt.Sprintf("%q", r.Body)
	return lp
}

func (lp *LogParams) injectResponseBody(ctx context.Context, rw *loggingResponseWriter) *LogParams {
	lp.fields[ResponseCodeKey] = rw.status
	lp.fields[ResponseKey] = rw.body
	return lp
}

type loggingResponseWriter struct {
	status int
	body   string
	http.ResponseWriter
}

func (w *loggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *loggingResponseWriter) Write(body []byte) (int, error) {
	w.body = string(body)
	return w.ResponseWriter.Write(body)
}

// getCaller retrieves the name of the first non this package calling function
func getCaller() *runtime.Frame {

	// cache this package's fully-qualified name
	callerInitOnce.Do(func() {
		// get log package name (this package)
		pcs := make([]uintptr, 2)
		frames := runtime.CallersFrames(pcs[:runtime.Callers(1, pcs)])
		frame, _ := frames.Next()
		thisPackageName = getPackageName(frame.Function)

		// now that we have the cache, we can skip a minimum count of known functions
		// XXX this is dubious, the number of frames may vary
		minimumCallerDepth = knownLogFrames
	})

	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != thisPackageName {
			return &f
		}
	}

	// if we got here, we failed to find the caller's context
	return nil
}

// getPackageName reduces a fully qualified function name to the package name
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		} else {
			break
		}
	}

	return f
}
