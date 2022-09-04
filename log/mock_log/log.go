// Code generated by MockGen. DO NOT EDIT.
// Source: logger.go

// Package mock_log is a generated GoMock package.
package mock_log

import (
	context "context"
	http "net/http"
	reflect "reflect"

	log "github.com/muhammad-fakhri/go-libs/log"
	gomock "github.com/golang/mock/gomock"
	logrus "github.com/sirupsen/logrus"
)

// MockSLogger is a mock of SLogger interface.
type MockSLogger struct {
	ctrl     *gomock.Controller
	recorder *MockSLoggerMockRecorder
}

// MockSLoggerMockRecorder is the mock recorder for MockSLogger.
type MockSLoggerMockRecorder struct {
	mock *MockSLogger
}

// NewMockSLogger creates a new mock instance.
func NewMockSLogger(ctrl *gomock.Controller) *MockSLogger {
	mock := &MockSLogger{ctrl: ctrl}
	mock.recorder = &MockSLoggerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSLogger) EXPECT() *MockSLoggerMockRecorder {
	return m.recorder
}

// BuildContextDataAndSetValue mocks base method.
func (m *MockSLogger) BuildContextDataAndSetValue(country, contextID string) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildContextDataAndSetValue", country, contextID)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// BuildContextDataAndSetValue indicates an expected call of BuildContextDataAndSetValue.
func (mr *MockSLoggerMockRecorder) BuildContextDataAndSetValue(country, contextID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildContextDataAndSetValue", reflect.TypeOf((*MockSLogger)(nil).BuildContextDataAndSetValue), country, contextID)
}

// Debug mocks base method.
func (m *MockSLogger) Debug(ctx context.Context, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debug", varargs...)
}

// Debug indicates an expected call of Debug.
func (mr *MockSLoggerMockRecorder) Debug(ctx interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockSLogger)(nil).Debug), varargs...)
}

// Debugf mocks base method.
func (m *MockSLogger) Debugf(ctx context.Context, message string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, message}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Debugf", varargs...)
}

// Debugf indicates an expected call of Debugf.
func (mr *MockSLoggerMockRecorder) Debugf(ctx, message interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, message}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debugf", reflect.TypeOf((*MockSLogger)(nil).Debugf), varargs...)
}

// Error mocks base method.
func (m *MockSLogger) Error(ctx context.Context, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Error", varargs...)
}

// Error indicates an expected call of Error.
func (mr *MockSLoggerMockRecorder) Error(ctx interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Error", reflect.TypeOf((*MockSLogger)(nil).Error), varargs...)
}

// ErrorMap mocks base method.
func (m *MockSLogger) ErrorMap(ctx context.Context, dataMap map[string]interface{}, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dataMap}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "ErrorMap", varargs...)
}

// ErrorMap indicates an expected call of ErrorMap.
func (mr *MockSLoggerMockRecorder) ErrorMap(ctx, dataMap interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dataMap}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ErrorMap", reflect.TypeOf((*MockSLogger)(nil).ErrorMap), varargs...)
}

// Errorf mocks base method.
func (m *MockSLogger) Errorf(ctx context.Context, message string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, message}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Errorf", varargs...)
}

// Errorf indicates an expected call of Errorf.
func (mr *MockSLoggerMockRecorder) Errorf(ctx, message interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, message}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Errorf", reflect.TypeOf((*MockSLogger)(nil).Errorf), varargs...)
}

// Fatal mocks base method.
func (m *MockSLogger) Fatal(ctx context.Context, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatal", varargs...)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockSLoggerMockRecorder) Fatal(ctx interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*MockSLogger)(nil).Fatal), varargs...)
}

// Fatalf mocks base method.
func (m *MockSLogger) Fatalf(ctx context.Context, message string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, message}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Fatalf", varargs...)
}

// Fatalf indicates an expected call of Fatalf.
func (mr *MockSLoggerMockRecorder) Fatalf(ctx, message interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, message}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatalf", reflect.TypeOf((*MockSLogger)(nil).Fatalf), varargs...)
}

// GetEntry mocks base method.
func (m *MockSLogger) GetEntry() *logrus.Entry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntry")
	ret0, _ := ret[0].(*logrus.Entry)
	return ret0
}

// GetEntry indicates an expected call of GetEntry.
func (mr *MockSLoggerMockRecorder) GetEntry() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntry", reflect.TypeOf((*MockSLogger)(nil).GetEntry))
}

// Info mocks base method.
func (m *MockSLogger) Info(ctx context.Context, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Info", varargs...)
}

// Info indicates an expected call of Info.
func (mr *MockSLoggerMockRecorder) Info(ctx interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockSLogger)(nil).Info), varargs...)
}

// InfoMap mocks base method.
func (m *MockSLogger) InfoMap(ctx context.Context, dataMap map[string]interface{}, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dataMap}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "InfoMap", varargs...)
}

// InfoMap indicates an expected call of InfoMap.
func (mr *MockSLoggerMockRecorder) InfoMap(ctx, dataMap interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dataMap}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InfoMap", reflect.TypeOf((*MockSLogger)(nil).InfoMap), varargs...)
}

// Infof mocks base method.
func (m *MockSLogger) Infof(ctx context.Context, message string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, message}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Infof", varargs...)
}

// Infof indicates an expected call of Infof.
func (mr *MockSLoggerMockRecorder) Infof(ctx, message interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, message}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Infof", reflect.TypeOf((*MockSLogger)(nil).Infof), varargs...)
}

// LogRequestResponse mocks base method.
func (m *MockSLogger) LogRequestResponse(ctx context.Context, data *log.RequestResponse, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, data}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "LogRequestResponse", varargs...)
}

// LogRequestResponse indicates an expected call of LogRequestResponse.
func (mr *MockSLoggerMockRecorder) LogRequestResponse(ctx, data interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, data}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogRequestResponse", reflect.TypeOf((*MockSLogger)(nil).LogRequestResponse), varargs...)
}

// SetContextData mocks base method.
func (m *MockSLogger) SetContextData(ctx context.Context, data *log.CommonFields) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetContextData", ctx, data)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// SetContextData indicates an expected call of SetContextData.
func (mr *MockSLoggerMockRecorder) SetContextData(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContextData", reflect.TypeOf((*MockSLogger)(nil).SetContextData), ctx, data)
}

// SetContextDataAndSetValue mocks base method.
func (m *MockSLogger) SetContextDataAndSetValue(r *http.Request, data map[string]string, country, contextId string) *http.Request {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetContextDataAndSetValue", r, data, country, contextId)
	ret0, _ := ret[0].(*http.Request)
	return ret0
}

// SetContextDataAndSetValue indicates an expected call of SetContextDataAndSetValue.
func (mr *MockSLoggerMockRecorder) SetContextDataAndSetValue(r, data, country, contextId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetContextDataAndSetValue", reflect.TypeOf((*MockSLogger)(nil).SetContextDataAndSetValue), r, data, country, contextId)
}

// SetLevel mocks base method.
func (m *MockSLogger) SetLevel(level logrus.Level) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLevel", level)
}

// SetLevel indicates an expected call of SetLevel.
func (mr *MockSLoggerMockRecorder) SetLevel(level interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLevel", reflect.TypeOf((*MockSLogger)(nil).SetLevel), level)
}

// Warn mocks base method.
func (m *MockSLogger) Warn(ctx context.Context, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warn", varargs...)
}

// Warn indicates an expected call of Warn.
func (mr *MockSLoggerMockRecorder) Warn(ctx interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockSLogger)(nil).Warn), varargs...)
}

// Warnf mocks base method.
func (m *MockSLogger) Warnf(ctx context.Context, message string, args ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, message}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "Warnf", varargs...)
}

// Warnf indicates an expected call of Warnf.
func (mr *MockSLoggerMockRecorder) Warnf(ctx, message interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, message}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warnf", reflect.TypeOf((*MockSLogger)(nil).Warnf), varargs...)
}