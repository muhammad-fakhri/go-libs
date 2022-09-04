// naming package is logrus because in SetReportCaller will set caller to first function outside logrus package
package log

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

//go:generate mockgen -destination mock_log/log.go -package mock_log -source logger.go

type SLogger interface {
	BuildContextDataAndSetValue(country string, contextID string) (cctx context.Context)
	SetContextDataAndSetValue(r *http.Request, data map[string]string, country string, contextId string) *http.Request
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

	InfoMap(ctx context.Context, dataMap map[string]interface{}, args ...interface{})
	ErrorMap(ctx context.Context, dataMap map[string]interface{}, args ...interface{})

	LogRequestResponse(ctx context.Context, data *RequestResponse, args ...interface{})
	SetContextData(ctx context.Context, data *CommonFields) (cctx context.Context)

	SetLevel(level log.Level)
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
	ContextUserIdKey  = "user_id"
	ContextEventIdKey = "event_id"
)

type SLog struct {
	entry *log.Entry
}

// context key data added to map
type contextData struct {
	country   string
	contextId string
}

// map contains value need to be displayed
type fields log.Fields

const (
	maximumCallerDepth int = 25
	minimumCallerDepth int = 1
)

func NewSLogger(service string) SLogger {
	entry, _ := getEntryAndLogger(service)
	return &SLog{entry}
}

func getEntryAndLogger(service string) (*log.Entry, *log.Logger) {
	logger := log.New()

	logger.SetFormatter(&log.JSONFormatter{})
	entry := log.NewEntry(logger)
	entry = entry.WithField("service", service)
	return entry, logger
}

func (l *SLog) BuildContextDataAndSetValue(country string, contextId string) (ctx context.Context) {
	data := make(map[string]string, 0)
	data[ContextCountryKey] = country
	data[ContextIdKey] = contextId

	ctx = context.WithValue(context.Background(), ContextDataMapKey, data)

	return ctx
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

func (l *SLog) GetEntry() *log.Entry {
	return l.entry
}

func (l *SLog) Infof(ctx context.Context, message string, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.InfoLevel))).Infof(message, args...)
}

func (l *SLog) Errorf(ctx context.Context, message string, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.ErrorLevel))).Errorf(message, args...)
}

func (l *SLog) Warnf(ctx context.Context, message string, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.WarnLevel))).Warnf(message, args...)
}

func (l *SLog) Debugf(ctx context.Context, message string, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.DebugLevel))).Debugf(message, args...)
}

func (l *SLog) Fatalf(ctx context.Context, message string, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.FatalLevel))).Fatalf(message, args...)
}

func (l *SLog) Info(ctx context.Context, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.InfoLevel))).Info(args...)
}

func (l *SLog) Error(ctx context.Context, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.ErrorLevel))).Error(args...)
}

func (l *SLog) Warn(ctx context.Context, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.WarnLevel))).Warn(args...)
}

func (l *SLog) Debug(ctx context.Context, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.DebugLevel))).Debug(args...)
}

func (l *SLog) Fatal(ctx context.Context, args ...interface{}) {
	l.entry.WithFields(log.Fields(getDefaultData(ctx, log.FatalLevel))).Fatal(args...)
}

func (l *SLog) InfoMap(ctx context.Context, dataMap map[string]interface{}, args ...interface{}) {
	data := getDefaultData(ctx, log.InfoLevel)
	data.getFieldsFromDataMap(dataMap)
	l.entry.WithFields(log.Fields(data)).Info(args...)
}

func (l *SLog) ErrorMap(ctx context.Context, dataMap map[string]interface{}, args ...interface{}) {
	data := getDefaultData(ctx, log.ErrorLevel)
	data.getFieldsFromDataMap(dataMap)
	l.entry.WithFields(log.Fields(data)).Error(args...)
}

func (l *SLog) SetLevel(level log.Level) {
	l.entry.Logger.SetLevel(level)
}

func (f fields) getFieldsFromContext(ctx context.Context) {
	dataMap := ctx.Value(ContextDataMapKey)

	if dataMap != nil {
		if data, ok := dataMap.(map[string]string); ok {
			for key, value := range data {
				f[key] = value
			}
		}
	}
}

func (f fields) getFieldsFromDataMap(dataMap map[string]interface{}) {
	if dataMap != nil {
		for key, value := range dataMap {
			f[key] = value
		}
	}
}

func getDefaultData(ctx context.Context, logLevel log.Level) fields {
	data := fields(make(map[string]interface{}))
	if logLevel <= log.ErrorLevel { // only print func and field on level lower or equal error
		data.getCallStackTrace()
	}
	data.getFieldsFromContext(ctx)
	return data
}

// currently not supporting caller package named 'log'
func (f fields) getCallStackTrace() {
	f.getCaller(getFrame())
}

func (f fields) getCaller(caller *runtime.Frame) {
	if caller == nil {
		return
	}

	funcVal := caller.Function
	fileVal := fmt.Sprintf("%s:%d", caller.File, caller.Line)
	if funcVal != "" {
		f["func"] = funcVal
	}
	if fileVal != "" {
		f["file"] = fileVal
	}
}

// getCaller retrieves the name of the first non this package calling function
func getFrame() *runtime.Frame {
	// Restrict the lookback frames to avoid runaway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)
		// hard code here because package runtime Caller behaves differently between golang version
		thisPackageName := "log"
		lenPkgName := len(pkg)
		thisPkgNameIndex := strings.Index(pkg, "/"+thisPackageName)

		// If the caller isn't part of this package, we're done
		if thisPkgNameIndex == -1 || thisPkgNameIndex != lenPkgName-4 {
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
