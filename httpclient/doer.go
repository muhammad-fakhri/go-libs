package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/muhammad-fakhri/go-libs/log"
)

const (
	valueLogTypeEgressHttp = "egress_http"

	statusGeneralError = 1
)

type HttpDoer interface {
	Do(req *http.Request) ([]byte, error)
	DoV2(req *http.Request) ([]byte, int, error)
	DoRawResponse(req *http.Request) (*http.Response, error)
	DoRawResponseWithoutLogging(req *http.Request) (*HTTPCallInfo, error)
}

type httpDoer struct {
	client    *http.Client
	logger    log.SLogger
	logConfig *LogConfig
}

func NewHttpDo(timeout time.Duration, optionalLogger ...log.SLogger) HttpDoer {
	client := &http.Client{
		Timeout: timeout,
	}

	return NewHttpDoWithClient(client, optionalLogger...)
}

func NewHttpDoWithClient(client *http.Client, optionalLogger ...log.SLogger) HttpDoer {
	if len(optionalLogger) < 1 || optionalLogger[0] == nil {
		return newDoer(client, log.NewSLogger("httpclient"))
	}
	return newDoer(client, optionalLogger[0])
}

func NewHttpDoWithLogConfig(timeout time.Duration, logConf *LogConfig, optionalLogger ...log.SLogger) HttpDoer {
	client := &http.Client{
		Timeout: timeout,
	}

	logConfig := defaultLogConfig()
	if logConf != nil {
		logConfig = NewLogConfig(logConf)
	}

	if len(optionalLogger) < 1 || optionalLogger[0] == nil {
		return newDoerWithLogConfig(client, log.NewSLogger("httpclient"), logConfig)
	}

	return newDoerWithLogConfig(client, optionalLogger[0], logConfig)
}

type ClientParam struct {
	OptClient  *http.Client  //if undefined, create default httpclient with timeout:optTimeout
	OptLogConf *LogConfig    //if undefined, will log all message parts
	OptLogger  log.SLogger   //if undefined, will create internal logger instance
	OptTimeout time.Duration //if undefined, will use default value 5 s
}

func (p *ClientParam) GetTimeout() time.Duration {
	if p.OptTimeout == 0 {
		return time.Duration(5 * time.Second)
	}

	return p.OptTimeout
}

func (p *ClientParam) GetClient() *http.Client {
	if p.OptClient == nil {
		return &http.Client{
			Timeout: p.GetTimeout(),
		}
	}

	return p.OptClient
}

func (p *ClientParam) GetLogger() log.SLogger {
	if p.OptLogger == nil {
		return log.NewSLogger("httpclient")
	}

	return p.OptLogger
}

func (p *ClientParam) GetLogConfig() *LogConfig {
	if p.OptLogConf == nil {
		return defaultLogConfig()
	}

	return NewLogConfig(p.OptLogConf)
}

func NewHttpDoWithParam(param *ClientParam) HttpDoer {
	client := param.GetClient()
	logConfig := param.GetLogConfig()
	logger := param.GetLogger()

	return newDoerWithLogConfig(client, logger, logConfig)
}

func newDoer(client *http.Client, logger log.SLogger) HttpDoer {
	return &httpDoer{
		client:    client,
		logger:    logger,
		logConfig: defaultLogConfig(),
	}
}

func newDoerWithLogConfig(client *http.Client, logger log.SLogger, logConfig *LogConfig) HttpDoer {
	return &httpDoer{
		client:    client,
		logger:    logger,
		logConfig: logConfig,
	}
}

func (d *httpDoer) Do(req *http.Request) ([]byte, error) {
	resp, err := d.doApiCallAndLogging(req)
	if err != nil || resp == nil {
		return nil, err
	}

	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBytes, err
}

func (d *httpDoer) DoV2(req *http.Request) ([]byte, int, error) {
	resp, err := d.doApiCallAndLogging(req)
	if err != nil || resp == nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return respBytes, resp.StatusCode, err
}

func (d *httpDoer) DoRawResponse(req *http.Request) (resp *http.Response, err error) {
	resp, err = d.doApiCallAndLogging(req)
	return
}

func (d *httpDoer) DoRawResponseWithoutLogging(request *http.Request) (resp *HTTPCallInfo, err error) {

	requestStartTime := time.Now()
	response, httpErr := d.client.Do(request)
	requestDuration := time.Since(requestStartTime)

	resp = &HTTPCallInfo{
		Request:          request,
		Response:         response,
		Duration:         requestDuration,
		RequestTimestamp: requestStartTime,
		HttpError:        httpErr,
	}

	return resp, nil
}

func (d *httpDoer) doApiCallAndLogging(request *http.Request) (*http.Response, error) {
	getRequestBody := request.GetBody
	requestStartTime := time.Now()
	response, httpErr := d.client.Do(request)
	requestDuration := time.Since(requestStartTime)

	//Re-attach Request Body after read by http.client.do()
	if request.Body != nil {
		var err error
		request.Body, err = getRequestBody()
		if err != nil {
			return nil, err
		}
	}

	context, err := d.buildContextFromRequest(request)
	if err != nil {
		return nil, err
	}

	err = d.logApiCall(context, request, response, requestDuration, requestStartTime, httpErr)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (d *httpDoer) buildContextFromRequest(request *http.Request) (context.Context, error) {
	var country string
	if country = request.Header.Get("X-Country"); country == "" {
		country = request.Header.Get("X-Tenant")
	}

	requestID := ""
	requestBody, err := getRequestBodyJSON(request)
	if err != nil {
		return request.Context(), err
	}

	var body map[string]interface{}
	json.Unmarshal([]byte(requestBody), &body)
	if value, exists := body["request_id"]; exists {
		requestID = value.(string)
	} else {
		requestID = uuid.New().String()
	}

	return d.logger.BuildContextDataAndSetValue(country, requestID), nil
}

func (d *httpDoer) logApiCall(context context.Context, request *http.Request, response *http.Response, duration time.Duration, requestTimestamp time.Time, httpError error) error {
	data := &log.RequestResponse{}

	requestBody, err := getRequestBodyJSON(request)
	if err != nil {
		return err
	}

	data.Type = valueLogTypeEgressHttp
	data.URLPath = fmt.Sprintf("%s %s", request.Method, request.URL.String())
	data.RequestTimestamp = requestTimestamp
	data.DurationMs = duration.Milliseconds()

	if httpError != nil {
		data.Status = statusGeneralError
		d.logger.ErrorMap(context, data.ToDataMap(), httpError.Error())
		return httpError
	}

	responseBody, err := getResponseBodyJSON(response)
	if err != nil {
		return err
	}

	if response != nil {
		data.Status = response.StatusCode

		if d.logConfig.LogResponseHeader() {
			header := response.Header.Clone()
			header.Del("Authorization")
			data.ResponseHeader = response.Header
		}

		if d.logConfig.LogResponseBody() {
			data.ResponseBody = responseBody
		} else {
			data.ResponseBody = wipedMessage
		}
	}

	if d.logConfig.LogRequestHeader() {
		header := request.Header.Clone()
		header.Del("Authorization")
		data.RequestHeader = request.Header
	}

	if d.logConfig.LogRequestBody() {
		data.RequestBody = requestBody
	} else {
		data.RequestBody = wipedMessage
	}

	d.logger.LogRequestResponse(context, data)
	return nil
}

func getRequestBodyJSON(request *http.Request) (string, error) {
	if request.Body == nil {
		return "null", nil
	}

	requestBodyBytes, err := getBodyBytes(&request.Body)
	if err != nil {
		return "null", err
	}
	return string(requestBodyBytes), nil
}

func getResponseBodyJSON(response *http.Response) (string, error) {
	if response == nil || response.Body == nil || response.Body == http.NoBody {
		return "null", nil
	}

	responseBodyBytes, err := getBodyBytes(&response.Body)
	if err != nil {
		return "null", err
	}
	return string(responseBodyBytes), nil
}

func getBodyBytes(body *io.ReadCloser) ([]byte, error) {
	responseBodyBytes, err := ioutil.ReadAll(*body)
	*body = ioutil.NopCloser(bytes.NewBuffer(responseBodyBytes))
	return responseBodyBytes, err
}
