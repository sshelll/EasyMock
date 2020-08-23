package easymock

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type easyResponse struct {
	body   interface{}
	seeker io.ReadSeeker
}

func (er *easyResponse) init() {
	if er.seeker == nil {
		switch d := er.body.(type) {
		case string:
			er.seeker = strings.NewReader(d)
		case []byte:
			er.seeker = bytes.NewReader(d)
		}
	}
}

func (er *easyResponse) Read(p []byte) (n int, err error) {
	n, err = er.seeker.Read(p)
	if err == io.EOF {
		_, _ = er.seeker.Seek(io.SeekStart, io.SeekStart)
	}
	return
}

func (er *easyResponse) Close() (err error) {
	_, err = er.seeker.Seek(io.SeekStart, io.SeekStart)
	return
}

func (er *easyResponse) Clone() *easyResponse {
	return &easyResponse{
		body:   er.body,
		seeker: er.seeker,
	}
}

func newEasyResponse(data interface{}) *easyResponse {
	eResp := &easyResponse{
		body: data,
	}
	eResp.init()
	return eResp
}

func newHttpResponseWithString(statusCode int, body string) *http.Response {
	resp := newHttpResponse(statusCode, body)
	resp.ContentLength = int64(len([]byte(body)))
	return resp
}

func newHttpResponseWithBytes(statusCode int, body []byte) *http.Response {
	resp := newHttpResponse(statusCode, body)
	resp.ContentLength = int64(len(body))
	return resp
}

func newHttpResponse(statusCode int, body interface{}) *http.Response {
	return &http.Response{
		Status:     strconv.Itoa(statusCode),
		StatusCode: statusCode,
		Header:     http.Header{},
		Body:       newEasyResponse(body),
	}
}