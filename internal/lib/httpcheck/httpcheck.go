package httpcheck

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Checker struct {
	handler http.Handler
}

func New(handler http.Handler) *Checker {
	return &Checker{handler: handler}
}

func (c *Checker) Test(t *testing.T, method, path string) *Request {
	return &Request{
		t:       t,
		handler: c.handler,
		method:  method,
		path:    path,
		headers: make(http.Header),
	}
}

type Request struct {
	t       *testing.T
	handler http.Handler
	method  string
	path    string
	headers http.Header
	body    []byte
}

func (r *Request) WithHeader(key, value string) *Request {
	r.headers.Add(key, value)
	return r
}

func (r *Request) WithBody(body []byte) *Request {
	r.body = body
	return r
}

func (r *Request) Check() *Response {
	var bodyReader io.Reader
	if r.body != nil {
		bodyReader = bytes.NewReader(r.body)
	}

	req := httptest.NewRequest(r.method, r.path, bodyReader)
	req.Header = r.headers
	rec := httptest.NewRecorder()
	r.handler.ServeHTTP(rec, req)

	return &Response{
		t:        r.t,
		recorder: rec,
	}
}

type Response struct {
	t        *testing.T
	recorder *httptest.ResponseRecorder
}

func (r *Response) HasStatus(code int) *Response {
	r.t.Helper()
	assert.Equal(r.t, code, r.recorder.Code)
	return r
}

func (r *Response) HasJSON(want any) *Response {
	r.t.Helper()
	wantBytes, err := json.Marshal(want)
	if err != nil {
		r.t.Fatalf("failed to marshal expected JSON: %v", err)
	}
	assert.JSONEq(r.t, string(wantBytes), r.recorder.Body.String())
	return r
}

func (r *Response) HasHeader(key, value string) *Response {
	r.t.Helper()
	assert.Equal(r.t, value, r.recorder.Header().Get(key))
	return r
}

func (r *Response) HasString(want string) *Response {
	r.t.Helper()
	assert.Equal(r.t, want, r.recorder.Body.String())
	return r
}
