package log

import (
	"net/http"
	"testing"
	"time"
)

var (
	sampleMap          = map[string]interface{}{"url": "be-service/users", "method": "GET", "status": 200, FieldType: "ingress"}
	requestWithContext *http.Request
)

func BenchmarkSLog_InfoWithMap(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		logger.InfoMap(sampleContext, sampleMap)
	}
}

func BenchmarkSLog_RequestResponse(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ctx := logger.SetContextData(sampleContext, &CommonFields{
			UserID:  "1200909",
			Country: "ID",
		})
		logger.LogRequestResponse(ctx, &RequestResponse{
			Type:             "ingress_http",
			URLPath:          "GET path/of/the/method",
			RequestTimestamp: time.Now(),
			RequestBody:      sampleMap,
			Status:           200,
			ResponseBody:     sampleObjects,
			DurationMs:       1000,
		})
	}
}
