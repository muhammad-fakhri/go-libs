package uuid

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_GetUUICluster(t *testing.T) {
	userID := uint64(10)
	type fields struct {
		config *Config
		client *fakeHttpDoer
	}
	type args struct {
		filter *GetUUIClusterRequest
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		want         *GetUUIClusterResponse
		wantRespCode int
		wantErr      bool
		wantErrType  reflect.Type
	}{
		{
			name: "when operation API return HTTP error",
			fields: fields{
				config: &Config{
					baseURL:  "http://127.0.0.1",
					clientID: "ds_test",
				},
				client: &fakeHttpDoer{
					httpResp: nil,
					err:      errors.New("HTTP Error"),
				},
			},
			args: args{filter: &GetUUIClusterRequest{
				UserID:   userID,
				Version:  UUIClusterVersion,
				ClientID: "ds_test",
				Country:  "ID",
			}},
			want:         nil,
			wantRespCode: http.StatusInternalServerError,
			wantErr:      true,
			wantErrType:  reflect.TypeOf(errors.New("")),
		},
		{
			name: "when operation API return error",
			fields: fields{
				config: &Config{
					baseURL:  "http://127.0.0.1",
					clientID: "ds_test",
				},
				client: &fakeHttpDoer{
					resp:       `{"title":"400 Bad Request"}`,
					err:        nil,
					statusCode: http.StatusBadRequest,
				},
			},
			args: args{filter: &GetUUIClusterRequest{
				UserID:   userID,
				Version:  UUIClusterVersion,
				ClientID: "ds_test",
				Country:  "ID",
			}},
			want:         nil,
			wantRespCode: http.StatusBadRequest,
			wantErr:      true,
			wantErrType:  reflect.TypeOf(errors.New("")),
		},
		{
			name: "when operation API return success",
			fields: fields{
				config: &Config{
					baseURL:  "http://127.0.0.1",
					clientID: "ds_test",
				},
				client: &fakeHttpDoer{
					resp:       `{"uui_cluster":[],"user_id":10,"exceed_limit":false}`,
					err:        nil,
					statusCode: http.StatusOK,
				},
			},
			args: args{filter: &GetUUIClusterRequest{
				UserID:   userID,
				Version:  UUIClusterVersion,
				ClientID: "ds_test",
				Country:  "ID",
			}},
			want: &GetUUIClusterResponse{
				ExceedLimit: false,
				UUICluster:  []uint64{},
				UserID:      userID,
			},
			wantRespCode: http.StatusOK,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				config:     tt.fields.config,
				httpClient: tt.fields.client,
			}

			gotCode, got, err := c.GetUUICluster(*tt.args.filter)
			if !reflect.DeepEqual(gotCode, tt.wantRespCode) {
				t.Errorf("Client.GetUUICluster() = %v, want %v", gotCode, tt.wantRespCode)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetUUICluster() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if reflect.TypeOf(err) != tt.wantErrType {
				t.Errorf("Client.GetUUICluster() error = %v, wantErrType %v", err, tt.wantErrType)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetUUICluster() = %v, want %v", got, tt.want)
			}
		})
	}
}
