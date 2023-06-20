package middleware

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/yulecd/pp-common/plog"
	"github.com/yulecd/pp-common/server"

	"github.com/gin-gonic/gin"
)

type logData struct {
	Method   string      `json:"method"`
	Path     string      `json:"path"`
	Cost     int64       `json:"cost"`
	Query    string      `json:"query"`
	Header   http.Header `json:"header"`
	Params   gin.Params  `json:"params"`
	Body     string      `json:"body"`
	Response string      `json:"response"`
}

type requestLogBody struct {
	io.ReadCloser
	body *bytes.Buffer
}

func (r requestLogBody) Read(p []byte) (n int, err error) {
	n, err = r.ReadCloser.Read(p)
	if err != nil {
		return
	}
	r.body.Write(p)
	return
}

type respLogWriter struct {
	gin.ResponseWriter
	resp *bytes.Buffer
}

func (r respLogWriter) Write(b []byte) (int, error) {
	r.resp.Write(b)
	return r.ResponseWriter.Write(b)
}

// LogReqResp 记录请求响应日志
func LogReqResp(c *gin.Context) {
	// 生产环境先打开详细日志
	//if os.Getenv("APP_ENV") != config.TestEnv && os.Getenv("APP_ENV") != config.DevEnv {
	//	c.Next()
	//	return
	//}

	debug := c.Query("common-debug")
	if len(debug) == 0 && c.Request != nil && c.Request.Header != nil {
		debug = c.Request.Header.Get("common-debug")
	}

	rw := &respLogWriter{
		resp:           bytes.NewBufferString(""),
		ResponseWriter: c.Writer,
	}
	c.Writer = rw

	b, _ := c.GetRawData()
	c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	data := make(map[string]interface{})
	data["client_ip"] = c.ClientIP()
	data["method"] = c.Request.Method
	data["path"] = c.Request.URL.Path
	data["query"] = c.Request.URL.RawQuery
	data["params"] = c.Params
	if len(string(b)) < 500 {
		data["body"] = string(b)
	} else {
		data["body"] = "content too long"
	}

	begin := time.Now()
	c.Next()
	end := time.Now()

	data["status"] = rw.Status()
	data["cost"] = float64(end.Sub(begin).Microseconds()) / 1000

	if len(rw.resp.String()) < 500 {
		data["response"] = rw.resp.String()
	} else {
		data["response"] = "content too long"
	}

	// 详细调试信息
	if len(debug) > 0 {
		data["header"] = c.Request.Header
		data["response_header"] = rw.Header()
	}

	plog.GetDefaultFieldEntry(server.NewContext(context.Background(), c)).WithFields(data).Info("route")
}
