package httpmiddleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/c2fo/testify/assert"
	"github.com/muhammad-fakhri/go-libs/log"
	"github.com/sirupsen/logrus"
)

var (
	defCountry          = "id"
	defContextid        = "abcdefghijklmnopq"
	defExcludeHeaderKey = "X-Exclude-Key"
)

type requestBody struct {
	Name string `json:"name"`
}

func getMockServer(logger log.SLogger) *httptest.Server {
	logIngresssMiddleware := NewIngressLogMiddleware(logger)

	mux := http.NewServeMux()
	mux.Handle("/hello", logIngresssMiddleware.Enforce(http.HandlerFunc(hello)))
	mux.Handle("/multiplewrite", logIngresssMiddleware.Enforce(http.HandlerFunc(multipleWriteBody)))
	mux.Handle("/panic", logIngresssMiddleware.Enforce(http.HandlerFunc(panicHandler)))
	mux.Handle("/req-id", requestIDMiddleware(logIngresssMiddleware.Enforce(http.HandlerFunc(multipleWriteBody))))

	mockServer := httptest.NewServer(mux)
	return mockServer
}

func getMockServerWithConfig(logger log.SLogger, config *Config) *httptest.Server {
	logIngresssMiddleware := NewIngressLogMiddleware(logger, config)

	mux := http.NewServeMux()
	mux.Handle("/hello", logIngresssMiddleware.Enforce(http.HandlerFunc(hello)))
	mux.Handle("/exclude-options-success", logIngresssMiddleware.Enforce(http.HandlerFunc(excludeOptionsSuccessHandler)))
	mux.Handle("/exclude-options-error", logIngresssMiddleware.Enforce(http.HandlerFunc(excludeOptionsErrorHandler)))

	mockServer := httptest.NewServer(mux)
	return mockServer
}

func hello(writer http.ResponseWriter, request *http.Request) {
	// read body
	responseBodyBytes, _ := ioutil.ReadAll(request.Body)
	// modify header
	request.Header.Add("add", "new-value")

	time.Sleep(1 * time.Second)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseBodyBytes) // to match it with the request body
}

func multipleWriteBody(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("halow", request.Context().Value(log.ContextDataMapKey))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hello "))
	writer.Write([]byte("World"))
	writer.Write([]byte("!"))
}

func panicHandler(writer http.ResponseWriter, request *http.Request) {
	time.Sleep(111 * time.Millisecond)
	testPanic(nil)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("Hello "))
	writer.Write([]byte("World"))
	writer.Write([]byte("!"))
}

func excludeOptionsSuccessHandler(writer http.ResponseWriter, request *http.Request) {
	// read body
	responseBodyBytes, _ := ioutil.ReadAll(request.Body)
	// modify header
	request.Header.Add("add", "new-value")
	request.Header.Add(defExcludeHeaderKey, "exclude-header-value")

	time.Sleep(1 * time.Second)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseBodyBytes) // to match it with the request body
}

func excludeOptionsErrorHandler(writer http.ResponseWriter, request *http.Request) {
	// read body
	responseBodyBytes, _ := ioutil.ReadAll(request.Body)
	// modify header
	request.Header.Add("add", "new-value")
	request.Header.Add(defExcludeHeaderKey, "exclude-header-value")

	time.Sleep(1 * time.Second)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusInternalServerError)
	writer.Write(responseBodyBytes) // to match it with the request body
}

func testPanic(as *requestBody) {
	fmt.Println(as.Name)
}

// requestIDMiddleware assign request id middleware
func requestIDMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]string, 0)
		data[log.ContextCountryKey] = defCountry
		data[log.ContextIdKey] = defContextid

		ctx := context.WithValue(r.Context(), log.ContextDataMapKey, data)
		r2 := r.Clone(ctx)

		h.ServeHTTP(w, r2)
	})
}

func TestLogIngressMessage(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")
	mockServer := getMockServer(logger)
	defer mockServer.Close()

	reqBody, err := json.Marshal(&requestBody{
		Name: "muhammad-fakhri",
	})
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/hello", bytes.NewReader(reqBody))
	req.Header.Add("X-Country", "ID")
	req.Header.Add("Authorization", "Bearer abcdefghijkl")

	client := &http.Client{}
	_, err = client.Do(req)

	time.Sleep(100 * time.Millisecond)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))

	logMessage := extractLogMessage(t, hook.LastEntry().Data)

	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "ID", logMessage.ReqHeader.Get("X-Country"))
	assert.Empty(t, logMessage.ReqHeader.Get("Authorization"))
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))

	// check if response body contains the same body after the reqBody has been read twice
	assert.Equal(t, string(reqBody), logMessage.ResponseBody)
	assert.Equal(t, logMessage.ReqBody, logMessage.ResponseBody)
	assert.True(t, logMessage.TimeTakenInMS >= (1*time.Second).Milliseconds())
}

func TestLogMessageResponseBodyWhenMultipleWrite(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")
	mockServer := getMockServer(logger)
	defer mockServer.Close()

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/multiplewrite", nil)
	req.Header.Add("X-Country", "ID")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("unexpected nil response")
	}
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)

	respBody, _ := ioutil.ReadAll(resp.Body)

	logMessage := extractLogMessage(t, hook.LastEntry().Data)

	assert.Nil(t, err)
	assert.Equal(t, "Hello World!", logMessage.ResponseBody)
	assert.Equal(t, "Hello World!", string(respBody))
}

// TestLogMessageResponsePanic to check if ingress log exists service panic
func TestLogMessageResponsePanic(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")
	mockServer := getMockServer(logger)
	defer mockServer.Close()

	reqBody, err := json.Marshal(&requestBody{
		Name: "muhammad-fakhri",
	})
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/panic", bytes.NewReader(reqBody))
	req.Header.Add("X-Country", "ID")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("unexpected nil response")
	}
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)

	respBody, _ := ioutil.ReadAll(resp.Body)

	logMessage := extractLogMessage(t, hook.LastEntry().Data)

	assert.Nil(t, err)
	// check if request is logged
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "ID", logMessage.ReqHeader.Get("X-Country"))
	assert.Empty(t, logMessage.ReqHeader.Get("Authorization"))
	assert.Equal(t, string(reqBody), logMessage.ReqBody)

	// response
	assert.Contains(t, logMessage.ResponseBody, "panic")
	assert.Contains(t, string(respBody), "panic")
}

// TestRequestIDUnchanged to check if request id on ingress log unchanged if exists before
func TestRequestIDUnchanged(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")
	mockServer := getMockServer(logger)
	defer mockServer.Close()

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/req-id", nil)
	// req.Header.Add("X-Country", "ID")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("unexpected nil response")
	}
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, defCountry, hook.LastEntry().Data["country"].(string))
	assert.Equal(t, defContextid, hook.LastEntry().Data["context_id"].(string))
}

func TestLogIngressMessageWithUserID(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")
	mockServer := getMockServer(logger)
	defer mockServer.Close()

	reqBody, err := json.Marshal(&requestBody{
		Name: "muhammad-fakhri",
	})
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/hello", bytes.NewReader(reqBody))
	req.Header.Add("X-Country", "ID")
	req.Header.Add("X-User-ID", "123456")
	req.Header.Add("Authorization", "Bearer abcdefghijkl")

	client := &http.Client{}
	_, err = client.Do(req)

	time.Sleep(100 * time.Millisecond)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))
	assert.Equal(t, "123456", hook.LastEntry().Data["user_id"].(string))

	logMessage := extractLogMessage(t, hook.LastEntry().Data)

	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "ID", logMessage.ReqHeader.Get("X-Country"))
	assert.Equal(t, "123456", logMessage.ReqHeader.Get("X-User-ID"))
	assert.Empty(t, logMessage.ReqHeader.Get("Authorization"))
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))

	// check if response body contains the same body after the reqBody has been read twice
	assert.Equal(t, string(reqBody), logMessage.ResponseBody)
	assert.Equal(t, logMessage.ReqBody, logMessage.ResponseBody)
	assert.True(t, logMessage.TimeTakenInMS >= (1*time.Second).Milliseconds())
}

func TestLogIngressMessageExcludeOptionsSuccess(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")

	config := &Config{
		ExcludeOpt: &ExcludeOption{
			SuccessResponseBody: true,
			RequestHeaderKeys:   []string{defExcludeHeaderKey},
		},
	}

	mockServer := getMockServerWithConfig(logger, config)
	defer mockServer.Close()

	reqBody, err := json.Marshal(&requestBody{
		Name: "long repetitive success response body need to be excluded",
	})
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/exclude-options-success", bytes.NewReader(reqBody))
	req.Header.Add("X-Country", "ID")
	req.Header.Add("Authorization", "Bearer abcdefghijkl")

	client := &http.Client{}
	_, err = client.Do(req)

	time.Sleep(100 * time.Millisecond)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))

	logMessage := extractLogMessage(t, hook.LastEntry().Data)

	assert.Equal(t, http.StatusOK, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "ID", logMessage.ReqHeader.Get("X-Country"))
	assert.Empty(t, logMessage.ReqHeader.Get("Authorization"))
	assert.Empty(t, logMessage.ReqHeader.Get(defExcludeHeaderKey))
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))

	// check if response body contains the same body after the reqBody has been read twice
	//assert.Equal(t, wipedMessage, logMessage.ResponseBody)
	assert.Equal(t, wipedMessage, logMessage.ResponseBody)
	assert.True(t, logMessage.TimeTakenInMS >= (1*time.Second).Milliseconds())
}

func TestLogIngressMessageExcludeOptionsError(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")

	config := &Config{
		ExcludeOpt: &ExcludeOption{
			SuccessResponseBody: true,
			RequestHeaderKeys:   []string{defExcludeHeaderKey},
		},
	}

	mockServer := getMockServerWithConfig(logger, config)
	defer mockServer.Close()

	reqBody, err := json.Marshal(&requestBody{
		Name: "error message",
	})
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/exclude-options-error", bytes.NewReader(reqBody))
	req.Header.Add("X-Country", "ID")
	req.Header.Add("Authorization", "Bearer abcdefghijkl")

	client := &http.Client{}
	_, err = client.Do(req)

	time.Sleep(100 * time.Millisecond)

	assert.Nil(t, err)
	assert.True(t, len(hook.LastEntry().Data["context_id"].(string)) > 0)
	assert.Equal(t, "ID", hook.LastEntry().Data["country"].(string))

	logMessage := extractLogMessage(t, hook.LastEntry().Data)

	assert.Equal(t, http.StatusInternalServerError, logMessage.ResponseCode)
	assert.Equal(t, http.MethodGet, logMessage.ReqMethod)
	assert.Equal(t, "ID", logMessage.ReqHeader.Get("X-Country"))
	assert.Empty(t, logMessage.ReqHeader.Get("Authorization"))
	assert.Empty(t, logMessage.ReqHeader.Get(defExcludeHeaderKey))
	assert.Equal(t, "application/json", logMessage.ResponseHeader.Get("Content-Type"))

	// check if response body contains the same body after the reqBody has been read twice
	assert.Equal(t, string(reqBody), logMessage.ResponseBody)
	assert.Equal(t, logMessage.ReqBody, logMessage.ResponseBody)
	assert.True(t, logMessage.TimeTakenInMS >= (1*time.Second).Milliseconds())
}

func extractLogMessage(t *testing.T, mssg logrus.Fields) *LogMessage {
	logMessage := &LogMessage{}

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

func extractHeader(t *testing.T, header string) http.Header {
	var result = make(http.Header)

	newHead := strings.ReplaceAll(header, "map", "")
	newHead = strings.ReplaceAll(newHead, "[", "")
	newHead = strings.ReplaceAll(newHead, "]", "")

	if newHead == "" {
		return result
	}

	messageComponents := strings.Split(newHead, " ")

	for i := 0; i < len(messageComponents); i++ {
		keyValueMessage := strings.Split(messageComponents[i], ":")

		result[keyValueMessage[0]] = []string{keyValueMessage[1]}
	}

	return result
}

func TestDisableLogIngressMessage(t *testing.T) {
	logger, hook := log.NewSLoggerWithTestHook("log-ingress-middleware")
	mockServer := getMockServerWithConfig(logger, &Config{DisableIngressLog: true})
	defer mockServer.Close()

	reqBody, err := json.Marshal(&requestBody{
		Name: "muhammad-fakhri",
	})
	if err != nil {
		t.Error(err)
	}

	req, _ := http.NewRequest(http.MethodGet, mockServer.URL+"/hello", bytes.NewReader(reqBody))
	req.Header.Add("X-Country", "ID")
	req.Header.Add("Authorization", "Bearer abcdefghijkl")

	client := &http.Client{}
	_, err = client.Do(req)

	time.Sleep(100 * time.Millisecond)

	assert.Nil(t, err)
	assert.Nil(t, hook.LastEntry())
}
