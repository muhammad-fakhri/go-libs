package httpclient

import (
	"net/http"
	"time"
)

const (
	ExcludeLog = true
	IncludeLog = false

	wipedMessage = "-"
)

// LogMessage is a struct to keep the log message easier
type LogMessage struct {
	URL            string
	ReqMethod      string
	ReqHeader      http.Header
	ReqBody        string
	ResponseHeader http.Header
	ResponseStatus string
	ResponseBody   string
	Duration       time.Duration
}

type LogConfig struct {
	ExcludeOpt *ExcludeOption
}

type ExcludeOption struct {
	RequestHeader  bool
	RequestBody    bool
	ResponseHeader bool
	ResponseBody   bool
}

func defaultLogConfig() *LogConfig {
	return &LogConfig{
		ExcludeOpt: &ExcludeOption{},
	}
}

func NewLogConfig(c *LogConfig) *LogConfig {
	if c == nil || c.ExcludeOpt == nil {
		c.ExcludeOpt = &ExcludeOption{}
	}

	return c
}

func (c *LogConfig) LogRequestHeader() bool {
	if c.ExcludeOpt == nil {
		return IncludeLog
	}

	return c.ExcludeOpt.RequestHeader == IncludeLog
}

func (c *LogConfig) LogRequestBody() bool {
	if c.ExcludeOpt == nil {
		return IncludeLog
	}

	return c.ExcludeOpt.RequestBody == IncludeLog
}

func (c *LogConfig) LogResponseHeader() bool {
	if c.ExcludeOpt == nil {
		return IncludeLog
	}

	return c.ExcludeOpt.ResponseHeader == IncludeLog
}

func (c *LogConfig) LogResponseBody() bool {
	if c.ExcludeOpt == nil {
		return IncludeLog
	}

	return c.ExcludeOpt.ResponseBody == IncludeLog
}
