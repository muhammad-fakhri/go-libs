package uuid

import "net/http"

type fakeHttpDoer struct {
	httpResp   *http.Response
	resp       string
	statusCode int
	err        error
}

func (d *fakeHttpDoer) Do(request *http.Request) ([]byte, error) {
	return []byte(d.resp), d.err
}

func (d *fakeHttpDoer) DoV2(request *http.Request) ([]byte, int, error) {
	return []byte(d.resp), d.statusCode, d.err
}

func (d *fakeHttpDoer) DoRawResponse(request *http.Request) (*http.Response, error) {
	return d.httpResp, d.err
}
