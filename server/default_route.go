package server

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	"github.com/gin-gonic/gin"
)

const (
	debugTokenName       = "token"
	debugTokenHeaderName = "x-debug-token"
	debugTokenMd5        = "7c8165518e9da9064f55d2b0cd4cd902" 
	debugTokenSalt       = "richkeyu"
)

// 默认全局路由 一些公共功能
func registerDefaultRoute(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "%s", "pong")
	})

	// 空模板
	//stringTpl, _ := template.New("stringTpl").Parse("{{ . }}")
	//r.SetHTMLTemplate(stringTpl)

	// 路由列表
	debug := r.Group("/debug")
	debug.Use(DebugAuth)
	debug.GET("/route", func(c *gin.Context) {
		content := ""
		content += "**all global handler: \n"
		for _, handler := range r.Handlers {
			content += fmt.Sprintf("%s \n", runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name())
		}
		content += "\n**all route: \n"
		for _, v := range r.Routes() {
			content += fmt.Sprintf("%-6s %-25s --> %s \n", v.Method, v.Path, v.Handler)
		}

		c.String(http.StatusOK, "%s", content)
		//c.HTML(http.StatusOK, "stringTpl", content)
	})
}

// DebugAuth debug校验
// 为避免循环依赖默认路由的组件放在同包内
func DebugAuth(c *gin.Context) {
	token := c.Request.URL.Query().Get(debugTokenName)
	if len(token) == 0 {
		token = c.Request.Header.Get(debugTokenHeaderName)
	}

	md5Token := fmt.Sprintf("%x", md5.Sum([]byte(token+debugTokenSalt)))
	if md5Token != debugTokenMd5 {
		c.String(http.StatusUnauthorized, "StatusUnauthorized")
		c.Abort()
		return
	}

	c.Next()
}
