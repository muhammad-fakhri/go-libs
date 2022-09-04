package log

import (
	"context"
	"net/http"
	"time"
)

const (
	FieldType           = "type"
	FieldURL            = "url_path"
	FieldReqMethod      = "method"
	FieldReqHeader      = "req_header"
	FieldReqBody        = "req_body"
	FieldResponseHeader = "rsp_header"
	FieldStatus         = "status"
	FieldResponseBody   = "rsp_body"
	FieldDurationMs     = "duration_ms"
	FieldReqTimestamp   = "req_timestamp"
)

type RequestResponse struct {
	Type             string                 `json:"type"`
	URLPath          string                 `json:"url_path"`
	RequestHeader    http.Header            `json:"req_header"`
	RequestBody      interface{}            `json:"req_body"`
	ResponseHeader   http.Header            `json:"rsp_header"`
	ResponseBody     interface{}            `json:"rsp_body"`
	Status           int                    `json:"status"`
	DurationMs       int64                  `json:"duration"`
	RequestTimestamp time.Time              `json:"req_timestamp"`
	DataMap          map[string]interface{} `json:"data_map"` // for additional data/custom field
}

func (r *RequestResponse) ToDataMap() map[string]interface{} {
	dataMap := make(map[string]interface{}, 0)
	dataMap[FieldType] = r.Type
	dataMap[FieldURL] = r.URLPath
	dataMap[FieldReqTimestamp] = r.RequestTimestamp.Unix()
	dataMap[FieldStatus] = r.Status
	dataMap[FieldDurationMs] = r.DurationMs

	if v := r.RequestHeader; v != nil {
		dataMap[FieldReqHeader] = v
	}

	if v := r.RequestBody; v != nil {
		dataMap[FieldReqBody] = v
	}

	if v := r.ResponseHeader; v != nil {
		dataMap[FieldResponseHeader] = v
	}

	if v := r.ResponseBody; v != nil {
		dataMap[FieldResponseBody] = v
	}

	if r.DataMap != nil {
		for key, value := range r.DataMap {
			dataMap[key] = value
		}
	}

	return dataMap
}

type CommonFields struct {
	ContextID string            `json:"context_id"`
	UserID    string            `json:"user_id"`
	Country   string            `json:"country"`
	EventID   string            `json:"event_id"`
	DataMap   map[string]string `json:"data_map"` // for additional data/custom field
}

func (c *CommonFields) ToDataMap() map[string]string {
	dataMap := make(map[string]string, 0)

	c.SetString(dataMap, ContextIdKey, c.ContextID)
	c.SetString(dataMap, ContextUserIdKey, c.UserID)
	c.SetString(dataMap, ContextCountryKey, c.Country)
	c.SetString(dataMap, ContextEventIdKey, c.EventID)

	if c.DataMap != nil {
		for key, value := range c.DataMap {
			dataMap[key] = value
		}
	}

	return dataMap
}

func (c *CommonFields) SetString(dataMap map[string]string, key string, value string) {
	if len(value) > 0 {
		dataMap[key] = value
	}
}

func (l *SLog) LogRequestResponse(ctx context.Context, data *RequestResponse, args ...interface{}) {
	l.InfoMap(ctx, data.ToDataMap(), args...)
}

func (l *SLog) SetContextData(ctx context.Context, data *CommonFields) (cctx context.Context) {
	return context.WithValue(ctx, ContextDataMapKey, data.ToDataMap())
}
