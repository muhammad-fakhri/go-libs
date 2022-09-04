package uuid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/muhammad-fakhri/go-libs/httpclient"
)

type Client struct {
	config     *Config
	httpClient httpclient.HttpDoer
}

func NewClient(config *Config, httpClient httpclient.HttpDoer) Clienter {
	return &Client{
		config:     config,
		httpClient: httpClient,
	}
}

type Clienter interface {
	GetUUICluster(reqBody GetUUIClusterRequest) (int, *GetUUIClusterResponse, error)
}

type DefaultErrorResponse struct {
	Title string `json:"title"`
}

func (c *Client) getErrorResponse(resp []byte) (error, error) {
	var errResp DefaultErrorResponse
	err := json.Unmarshal(resp, &errResp)
	if err != nil {
		log.Printf("failed to parse response %+v, %+v", err, string(resp))
		return nil, err
	}

	return fmt.Errorf(errResp.Title), nil
}

func (c *Client) commonRequest(method, path string, request, response interface{}) (int, error) {
	f, err := json.Marshal(request)
	url := c.config.baseURL + path
	req, err := http.NewRequest(method, url, bytes.NewReader(f))
	if err != nil {
		log.Println("failed to create request", err)
		return http.StatusInternalServerError, err
	}

	respByte, statusCode, err := c.httpClient.DoV2(req)
	if err != nil {
		log.Println("failed to do request response body", err)
		return http.StatusInternalServerError, err
	}

	if statusCode != http.StatusOK {
		errResp, err := c.getErrorResponse(respByte)
		if err != nil {
			log.Println("failed to parse err response", err)
			return http.StatusInternalServerError, err
		}
		return statusCode, errResp
	}

	err = json.Unmarshal(respByte, &response)
	if err != nil {
		log.Printf("failed to parse response to json %+v", err)
		return http.StatusInternalServerError, err
	}

	return statusCode, nil
}
