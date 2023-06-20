package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	logpkg "log"
	"net/http"
	"net/http/httputil"
	urlpkg "net/url"
	"strings"
	"time"

	"github.com/yulecd/pp-common/plog"

	"github.com/sirupsen/logrus"
)

type Request struct {
	ctx  context.Context
	opts Options
	cli  *http.Client
	req  *http.Request
	body io.Reader

	startTime time.Time
}

func NewRequest(req *http.Request) *Request {
	return &Request{
		req: req,
	}
}

func (r *Request) WithContext(context context.Context) *Request {
	r.ctx = context

	return r
}

// Log 从通用的请求上下文里获取日志对象
func (r *Request) Log() *plog.Entry {
	if s := plog.GetDefaultFieldEntry(r.ctx); s != nil {
		return s
	}

	return nil
}

// GetRequest returns http request ptr.
func (r *Request) GetRequest() *http.Request {
	return r.req
}

// Get sends request with get method.
func (r *Request) Get(uri string, opts ...Options) (*Response, error) {
	return r.Request("GET", uri, opts...)
}

// Post sends request with post method
func (r *Request) Post(uri string, opts ...Options) (*Response, error) {
	return r.Request("POST", uri, opts...)
}

// Put sends request with put method
func (r *Request) Put(uri string, opts ...Options) (*Response, error) {
	return r.Request("PUT", uri, opts...)
}

// Path sends request with patch method
func (r *Request) Patch(uri string, opts ...Options) (*Response, error) {
	return r.Request("PATCH", uri, opts...)
}

// Delete sends request with delete method
func (r *Request) Delete(uri string, opts ...Options) (*Response, error) {
	return r.Request("DELETE", uri, opts...)
}

// Options sends request with options method
func (r *Request) Options(uri string, opts ...Options) (*Response, error) {
	return r.Request("OPTIONS", uri, opts...)
}

// doRequest is last wrapper and execute request.Do
var doRequest = func(ctx context.Context, req *Request) (*Response, error) {
	return req.sendRequest()
}

// Request encapsulates http.request internal logic
func (r *Request) Request(method, uri string, opts ...Options) (*Response, error) {
	r.startTime = time.Now()
	defaultOpts := r.opts
	if len(opts) > 0 {
		r.opts.merge(opts[0])
	}

	if false == IsValidHttpUrl(uri) {
		baseUri := defaultOpts.BaseURI
		if r.opts.BaseURI != "" {
			baseUri = r.opts.BaseURI
		}
		if baseUri == "" {
			return nil, errors.New("invalid uri: empty uri")
		}
		uri = baseUri + uri
	}

	switch method {
	case http.MethodGet, http.MethodDelete:
		req, err := http.NewRequest(method, uri, nil)
		if err != nil {
			return nil, err
		}

		r.req = req
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions:
		r.parseBody()

		req, err := http.NewRequest(method, uri, r.body)
		if err != nil {
			return nil, err
		}

		r.req = req
		if r.opts.FormParams != nil {
			r.req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		} else if r.opts.JSON != nil {
			r.req.Header.Set("Content-Type", "application/json")
		}
	default:
		return nil, errors.New("invalid request method")
	}

	r.mergeDefaultOpts(defaultOpts)
	r.parseClient()
	r.parseQuery()
	r.parseHeaders()

	dump, err := httputil.DumpRequest(r.req, true)
	if r.Log() != nil {
		r.Log().WithFields(logrus.Fields{
			"host":   r.req.URL.String(),
			"method": r.req.Method,
		}).Infof("client req: %s", string(dump))
	}
	if r.opts.debug && err == nil {
		logpkg.Printf("\n%s", dump)
	}

	root := doRequest
	if len(r.opts.wrappers) > 0 {
		for i := len(r.opts.wrappers) - 1; i >= 0; i-- {
			root = r.opts.wrappers[i](root)
		}
	}

	return root(r.ctx, r)
}

func (r *Request) sendRequest() (resp *Response, err error) {
	call := func(i int) error {
		t, err := r.opts.backoff(r.ctx, r, i)
		if err != nil {
			return err
		}

		if t.Seconds() > 0 {
			time.Sleep(t)
		}

		_resp, err := r.cli.Do(r.req)

		resp = &Response{
			req:   r,
			hresp: _resp,
			hreq:  r.req,
			err:   err,
		}

		if err != nil {
			if r.opts.debug {
				fmt.Println(err)
			}
		}

		if r.opts.debug && _resp != nil {
			dump, err := httputil.DumpResponse(_resp, true)
			if err == nil {
				logpkg.Printf("\n%s", dump)
			}
		}

		return err
	}

	ch := make(chan error, r.opts.retries)

	for i := 0; i < r.opts.retries; i++ {
		go func() {
			ch <- call(i)
		}()

		select {
		case <-r.ctx.Done():
			return nil, errors.New(fmt.Sprintf("http.client request timeout, err: %v", err))
		case err = <-ch:
			if err == nil {
				return
			}

			retry, rerr := r.opts.retry(r.ctx, r, resp, i, err)
			if rerr != nil {
				return nil, rerr
			}

			if !retry {
				return nil, err
			}
		}
	}

	return
}

func (r *Request) mergeDefaultOpts(opts Options) {
	for k, v := range opts.Headers {
		if _, ok := r.opts.Headers[k]; !ok {
			r.opts.Headers[k] = v
		}
	}
}

// parseClient setups request client
func (r *Request) parseClient() {
	r.cli = &http.Client{
		Timeout: r.opts.timeout,
	}
}

// parseQuery parses request query
func (r *Request) parseQuery() {
	switch r.opts.Query.(type) {
	case string:
		str := r.opts.Query.(string)
		r.req.URL.RawQuery = str
	case map[string]interface{}:
		q := r.req.URL.Query()
		for k, v := range r.opts.Query.(map[string]interface{}) {
			if vv, ok := v.(string); ok {
				q.Set(k, vv)
				continue
			} else if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					q.Add(k, vvv)
				}
			} else {
				if r.Log() != nil {
					r.Log().Infof("query param not string %v", v)
				}
				q.Set(k, fmt.Sprintf("%v", v))
				continue
			}
		}
		r.req.URL.RawQuery = q.Encode()
	}
}

// parseHeaders parses request headers
func (r *Request) parseHeaders() {
	if r.opts.Headers != nil {
		for k, v := range r.opts.Headers {
			if vv, ok := v.(string); ok {
				r.req.Header.Set(k, vv)
				continue
			}
			if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					r.req.Header.Add(k, vvv)
				}
			}
		}
	}
}

// parseBody parses request body
func (r *Request) parseBody() {
	if r.opts.FormParams != nil {
		values := urlpkg.Values{}
		for k, v := range r.opts.FormParams {
			if vv, ok := v.(string); ok {
				values.Set(k, vv)
			} else if vv, ok := v.([]string); ok {
				for _, vvv := range vv {
					values.Add(k, vvv)
				}
			} else {
				if r.Log() != nil {
					r.Log().Infof("query param not string %v", v)
				}
				values.Set(k, fmt.Sprintf("%v", v))
				continue
			}
		}
		r.body = strings.NewReader(values.Encode())
		return
	}
	if r.opts.JSON != nil {
		b, err := json.Marshal(r.opts.JSON)
		if err == nil {
			r.body = bytes.NewReader(b)
			return
		}
	}

	return
}
