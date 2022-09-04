package uuid

import (
	"net/http"
)

const (
	UUIClusterVersion = "v1"
)

type GetUUIClusterRequest struct {
	ClientID string `json:"client_id"`
	Version  string `json:"version"`
	Country  string `json:"country"`
	UserID   uint64 `json:"user_id"`
}

type GetUUIClusterResponse struct {
	UserID      uint64   `json:"user_id"`
	UUICluster  []uint64 `json:"uui_cluster"`
	ExceedLimit bool     `json:"exceed_limit"`
}

func (c *Client) GetUUICluster(reqBody GetUUIClusterRequest) (int, *GetUUIClusterResponse, error) {
	var result *GetUUIClusterResponse
	reqBody.ClientID = c.config.ClientID()
	reqBody.Version = UUIClusterVersion
	statusCode, err := c.commonRequest(http.MethodPost, c.config.GetUUIClusterPath(), reqBody, &result)
	if err != nil {
		return statusCode, nil, err
	}
	if result == nil {
		err = ErrNilResponse
		return http.StatusInternalServerError, nil, err
	}
	return statusCode, result, nil
}
