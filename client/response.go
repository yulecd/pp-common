package client

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	jsoniter "github.com/json-iterator/go"
)

// Response api响应
type Response struct {
	req      *Request
	hresp    *http.Response
	hreq     *http.Request
	err      error
	copyBody []byte
}

type ResponseBody []byte

func NewResponse(req *http.Request, resp *http.Response, err error) *Response {
	return &Response{
		hresp: resp,
		hreq:  req,
		err:   err,
	}
}

// String returns string struct of response body
func (r ResponseBody) String() string {
	return string(r)
}

// Read returns specified length bytes of response body
func (r ResponseBody) Read(length int) []byte {
	if length > len(r) {
		length = len(r)
	}

	return r[:length]
}

// GetContents returns string value of response body
func (r ResponseBody) GetContents() string {
	return string(r)
}

// GetRequest returns request object
func (r *Response) GetHttpRequest() *http.Request {
	return r.hreq
}

// GetBody returns response body
func (r *Response) GetBody() (ResponseBody, error) {
	defer r.hresp.Body.Close()
	var body []byte
	var err error

	if r.copyBody != nil {
		// 重复读取
		body = r.copyBody
	} else {
		// 第一次读取
		body, err = ioutil.ReadAll(r.hresp.Body)
		if err != nil {
			return nil, err
		}
		if body != nil {
			r.copyBody = body
		}
	}

	if r.req.Log() != nil {
		r.req.Log().WithFields(logrus.Fields{
			"host": r.req.req.URL.String(),
			"cost": float64(time.Now().Sub(r.req.startTime).Microseconds()) / 1000,
		}).Infof("client resp: %s", string(body))
	}

	return ResponseBody(body), nil
}

// GetStatusCode returns response code
func (r *Response) GetStatusCode() int {
	return r.hresp.StatusCode
}

// GetReasonPhrase returns response reason phrase
func (r *Response) GetReasonPhrase() string {
	status := r.hresp.Status
	arr := strings.Split(status, " ")

	return arr[1]
}

// IsOk returns true if statusCode is 200.
func (r *Response) IsOk() bool {
	if r.GetStatusCode() == http.StatusOK {
		return true
	}

	return false
}

// IsTimeout check request is timeout
func (r *Response) IsTimeout() bool {
	if r.err == nil {
		return false
	}
	netErr, ok := r.err.(net.Error)
	if !ok {
		return false
	}
	if netErr.Timeout() {
		return true
	}
	return false
}

// GetHeaders returns response header
func (r *Response) GetHeaders() map[string][]string {
	return r.hresp.Header
}

// GetHeader returns specified response header
func (r *Response) GetHeader(name string) []string {
	headers := r.GetHeaders()
	for k, v := range headers {
		if strings.ToLower(name) == strings.ToLower(k) {
			return v
		}
	}
	return nil
}

// ParseBody 按格式解析响应结果 不检查code
func (r *Response) ParseBody(obj interface{}) (resp Resp, err error) {
	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		err = fmt.Errorf("ParseBody param need pointer")
		return
	}

	if r.GetStatusCode() != http.StatusOK {
		err = fmt.Errorf("ParseBody http status error: %d %s", r.GetStatusCode(), r.GetReasonPhrase())
		return
	}

	respBody, err := r.GetBody()
	if err != nil {
		err = fmt.Errorf("ParseBody get body error: %w", err)
		return
	}

	err = jsoniter.Unmarshal([]byte(respBody.String()), &resp)
	if err != nil {
		err = fmt.Errorf("ParseBody unmarshal error: %w", err)
		return
	}

	dataStr, err := jsoniter.Marshal(resp.Data)
	if err != nil {
		err = fmt.Errorf("ParseBody marshal data error: %w", err)
		return
	}
	err = jsoniter.Unmarshal(dataStr, obj)
	if err != nil {
		err = fmt.Errorf("ParseBody unmarshal data error: %w", err)
		return
	}
	return resp, nil
}

// MustParseBody 按格式解析响应结果 检查code 如果不等于1 返回error
func (r *Response) MustParseBody(obj interface{}) (resp Resp, err error) {
	resp, err = r.ParseBody(obj)
	if err != nil {
		return resp, err
	}
	if !resp.IsSuccess() {
		return resp, fmt.Errorf("call api fail: %d %s", resp.Code, resp.Message)
	}
	return resp, nil
}
