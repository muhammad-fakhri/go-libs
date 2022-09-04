package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/muhammad-fakhri/go-libs/log"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var (
	MockServer *httptest.Server
	ReturnJson = `{"userId":1,"id":1,"title":"delectus aut autem","completed":false}`
	ReturnXml  = `<dateInformation><friendlyDate>May 28, 2013</friendlyDate><unixTime>1369739047</unixTime><monthNum>May</monthNum><dayOfWeek>Tuesday</dayOfWeek><yearNum>2013</yearNum></dateInformation>`
)

func getMockServer() *httptest.Server {
	if MockServer == nil {
		MockServer = httptest.NewServer(getHandlers())
	}
	return MockServer
}

func getHandlers() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/return/200/json", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintln(writer, ReturnJson)
	})

	mux.HandleFunc("/return/200/xml", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/xml")
		writer.WriteHeader(http.StatusOK)
		fmt.Fprintln(writer, ReturnXml)
	})

	mux.HandleFunc("/return/204", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusNoContent)
	})

	mux.HandleFunc("/return/500", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
	})

	mux.HandleFunc("/return/200/timeout", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		time.Sleep(10 * time.Millisecond)
		fmt.Fprintln(writer, ReturnJson)
	})

	return mux
}

type LogHTTPMessage struct {
	URL            string
	ReqMethod      string
	ReqHeader      http.Header
	ReqBody        string
	ResponseHeader http.Header
	ResponseCode   int
	ResponseBody   string
	TimeTakenInMS  int64
}

func extractLogMessage(mssg logrus.Fields) *LogHTTPMessage {
	logMessage := &LogHTTPMessage{}

	urlPath := strings.Split(mssg[log.FieldURL].(string), " ")
	logMessage.URL = urlPath[1]
	logMessage.ReqMethod = urlPath[0]
	logMessage.ResponseCode = mssg[log.FieldStatus].(int)
	logMessage.TimeTakenInMS = mssg[log.FieldDurationMs].(int64)
	logMessage.ReqHeader = mssg[log.FieldReqHeader].(http.Header)
	logMessage.ReqBody = mssg[log.FieldReqBody].(string)
	logMessage.ResponseHeader = mssg[log.FieldResponseHeader].(http.Header)
	logMessage.ResponseBody = mssg[log.FieldResponseBody].(string)
	return logMessage
}

func TestHttpDoer_Do_WithGetMethod_WithHeaderXCountry_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Country", "ID")

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
	assert.True(t, logMessage.TimeTakenInMS <= (1*time.Second).Milliseconds())
}

func TestHttpDoer_Do_WithGetMethod_WithHeaderXTenant_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Tenant", "ID")

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
	assert.True(t, logMessage.TimeTakenInMS <= (1*time.Second).Milliseconds())
	assert.Equal(t, "ID", logMessage.ReqHeader.Get("X-Tenant"))
}

func TestHttpDoer_Do_WithGetMethod_WithoutHeaderXTenantAndXCountry_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
	assert.True(t, logMessage.TimeTakenInMS <= (1*time.Second).Milliseconds())
	assert.Equal(t, "", logMessage.ReqHeader.Get("X-Tenant"))
}

func TestHttpDoer_Do_WithGetMethod_WithoutHeaderXTenantAndXCountry_ReturnXml(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/xml"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "application/xml", logMessage.ResponseHeader.Get("Content-Type"))
	assert.True(t, logMessage.TimeTakenInMS <= (1*time.Second).Milliseconds())
}

func TestHttpDoer_Do_WithPostMethod_WithRequestId_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodPost, logMessage.ReqMethod)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
	assert.True(t, logMessage.TimeTakenInMS <= (1*time.Second).Milliseconds())
}

func TestHttpDoer_Do_WithPostMethod_ReturnNoContent(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/204"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusNoContent, logMessage.ResponseCode)
	assert.Equal(t, http.MethodPost, logMessage.ReqMethod)
	assert.Equal(t, "", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_Do_WithPostMethod_ReturnNotFound(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/404"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusNotFound, logMessage.ResponseCode)
	assert.Equal(t, http.MethodPost, logMessage.ReqMethod)
	assert.Equal(t, "text/plain; charset=utf-8", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_Do_WithPostMethod_ReturnInternalServerError(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/500"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusInternalServerError, logMessage.ResponseCode)
}

func TestHttpDoer_DoV2_WithGetMethod_WithHeaderXCountry_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Country", "ID")

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))
}

func TestHttpDoer_DoV2_WithGetMethod_WithHeaderXTenant_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Tenant", "ID")

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
}

func TestHttpDoer_DoV2_WithGetMethod_WithoutHeaderXTenantAndXCountry_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_DoV2_WithGetMethod_WithoutHeaderXTenantAndXCountry_ReturnXml(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/xml"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, "application/xml", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_DoV2_WithPostMethod_WithRequestId_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))
}

func TestHttpDoer_DoV2_WithPostMethod_ReturnNoContent(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/204"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusNoContent, logMessage.ResponseCode)
}

func TestHttpDoer_DoV2_WithPostMethod_ReturnNotFound(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/404"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusNotFound, logMessage.ResponseCode)
}

func TestHttpDoer_DoV2_WithPostMethod_ReturnInternalServerError(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/500"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, _, err := httpClient.DoV2(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusInternalServerError, logMessage.ResponseCode)
}

func TestHttpDoer_DoRawResponse_WithGetMethod_WithHeaderXCountry_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Country", "ID")

	_, err := httpClient.DoRawResponse(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
}

func TestHttpDoer_DoRawResponse_WithGetMethod_WithHeaderXTenant_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Tenant", "ID")

	_, err := httpClient.DoRawResponse(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_DoRawResponse_WithPostMethod_WithRequestId_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.DoRawResponse(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodPost, logMessage.ReqMethod)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_DoRawResponse_WithPostMethod_ReturnNoContent(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/204"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.DoRawResponse(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusNoContent, logMessage.ResponseCode)
	assert.Equal(t, http.MethodPost, logMessage.ReqMethod)
	assert.Equal(t, "", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_DoRawResponse_WithPostMethod_ReturnNotFound(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/404"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.DoRawResponse(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusNotFound, logMessage.ResponseCode)
	assert.Equal(t, http.MethodPost, logMessage.ReqMethod)
}

func TestHttpDoer_DoRawResponse_WithPostMethod_ReturnInternalServerError(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/500"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	httpClient := newDoer(getMockServer().Client(), logger)
	requestMethod := http.MethodPost
	req, _ := http.NewRequest(requestMethod, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.DoRawResponse(req)

	assert.Nil(t, err)
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	logMessage := extractLogMessage(hook.LastEntry().Data)
	assert.Equal(t, http.StatusInternalServerError, logMessage.ResponseCode)
	assert.Equal(t, http.MethodPost, logMessage.ReqMethod)
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))
}

func TestHttpDoer_DoRawResponse_WithPostMethod_ReturnErrorRequestTimeout(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/timeout"
	payload := struct {
		RequestId string `json:"request_id"`
	}{
		RequestId: "mock_request_id",
	}
	payloadBytes, _ := json.Marshal(payload)

	logger, hook := log.NewSLoggerWithTestHook("httpClient")
	client := getMockServer().Client()
	client.Timeout = 1 * time.Millisecond
	httpClient := newDoer(client, logger)
	req, _ := http.NewRequest(http.MethodPost, apiUrl, bytes.NewReader(payloadBytes))

	_, err := httpClient.DoRawResponse(req)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "(Client.Timeout exceeded while awaiting headers)")
	assert.Equal(t, payload.RequestId, hook.LastEntry().Data["context_id"].(string))
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
}

func TestHttpDoer_Do_WithNoLogger_AndReturnsNoError(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	httpClient := NewHttpDo(10 * time.Second)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Country", "ID")

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
}

func TestHttpDoer_Do_WithLoggerEqualsToNil_AndReturnsNoError(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	httpClient := NewHttpDo(10*time.Second, nil)
	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Country", "ID")

	_, err := httpClient.Do(req)

	assert.Nil(t, err)
}

func TestHttpDoer_WithLogConfig_DoRawResponse_WithGetMethod_WithHeaderXCountry_ReturnJson(t *testing.T) {
	apiUrl := getMockServer().URL + "/return/200/json"
	logger, hook := log.NewSLoggerWithTestHook("httpClient")

	logExcludeOpt := &ExcludeOption{
		RequestHeader:  true,
		RequestBody:    true,
		ResponseHeader: false,
		ResponseBody:   false,
	}
	logConfig := &LogConfig{
		ExcludeOpt: logExcludeOpt,
	}

	httpClient := NewHttpDoWithLogConfig(10*time.Second, logConfig, logger)

	requestMethod := http.MethodGet
	req, _ := http.NewRequest(requestMethod, apiUrl, nil)
	req.Header.Add("X-Country", "ID")

	_, err := httpClient.DoRawResponse(req)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, fmt.Sprintf("%s", logrus.InfoLevel), fmt.Sprintf("%s", hook.LastEntry().Level))

	assert.Equal(t, wipedMessage, hook.LastEntry().Data[log.FieldReqBody].(string))
	assert.NotEqual(t, wipedMessage, hook.LastEntry().Data[log.FieldResponseBody].(string))
	assert.Nil(t, hook.LastEntry().Data[log.FieldReqHeader])
	assert.NotNil(t, hook.LastEntry().Data[log.FieldResponseHeader])
}
