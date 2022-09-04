package corsanywhere

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorsAnywhere(t *testing.T) {

	result := []byte("responded by handler")
	dummyOrigin := "http://blabla.com"

	var tts = []struct {
		caseName string
		request  func() *http.Request
		result   func(res *http.Response)
	}{
		{
			caseName: "when request method is not option",
			request: func() *http.Request {
				r, _ := http.NewRequest(http.MethodPost, "/", nil)

				return r
			},
			result: func(res *http.Response) {
				bod, err := ioutil.ReadAll(res.Body)
				assert.Nil(t, err)

				assert.Equal(t, result, bod)
				assert.Equal(t, "*", res.Header.Get("Access-Control-Allow-Origin"))
				assert.Equal(t, "true", res.Header.Get("Access-Control-Allow-Credentials"))
			},
		},
		{
			caseName: "when request method is an option",
			request: func() *http.Request {
				r, _ := http.NewRequest(http.MethodOptions, "/", nil)

				r.Header.Set("origin", dummyOrigin)
				r.Header.Set("Access-Control-Request-Method", http.MethodPost)
				r.Header.Set("Access-Control-Request-Headers", "X-User-ID")

				return r
			},
			result: func(res *http.Response) {
				assert.Equal(t, http.StatusOK, res.StatusCode)

				assert.Equal(t, dummyOrigin, res.Header.Get("Access-Control-Allow-Origin"))
				assert.Equal(t, http.MethodPost, res.Header.Get("Access-Control-Allow-Methods"))
				assert.Equal(t, "X-User-ID", res.Header.Get("Access-Control-Allow-Headers"))
				assert.Equal(t, "true", res.Header.Get("Access-Control-Allow-Credentials"))
			},
		},
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(result)
	})

	for _, tt := range tts {
		t.Log(tt.caseName)

		handler := CorsAnywhere(dummyHandler)
		rr := httptest.NewRecorder()
		req := tt.request()

		handler.ServeHTTP(rr, req)

		tt.result(rr.Result())
	}
}
