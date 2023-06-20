// http server 初始化

package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/yulecd/pp-common/config"

	"github.com/gin-gonic/gin"
)

// Config http server配置
type Config struct {
	Name         string        `yaml:"name"`
	HttpPort     int           `yaml:"http_port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// ServiceConfig http server配置
var ServiceConfig *Config

// 封装 默认http server
type server struct {
	*http.Server
}

// NewServer 生成一个http server
func NewServer(handler http.Handler) *server {
	loadConfig()

	if os.Getenv(config.AppEnvName) == config.ProdEnv || os.Getenv(config.AppEnvName) == config.PreEnv || os.Getenv(config.AppEnvName) == config.TestEnv {
		gin.SetMode(gin.ReleaseMode)
	}

	if ServiceConfig == nil || ServiceConfig.HttpPort == 0 {
		panic("ServiceConfig need set value")
	}

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", ServiceConfig.HttpPort),
		Handler:        handler,
		ReadTimeout:    ServiceConfig.ReadTimeout * time.Second,
		WriteTimeout:   ServiceConfig.WriteTimeout * time.Second,
		MaxHeaderBytes: 1 << 20, // 1M
	}

	registerDefaultRoute(handler.(*gin.Engine))

	return &server{
		s,
	}
}

// Run 启动服务
func (s *server) Run() {
	go func() {
		fmt.Printf("Server starting at http://127.0.0.1:%d \n", ServiceConfig.HttpPort)
		if err := s.ListenAndServe(); err != nil {
			fmt.Printf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	fmt.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(nil, "Server shutdown:", err)
	}

	fmt.Println("Server exiting")
}

// 根据配置中心规则加载默认配置
func loadConfig() {
	if err := config.Load("app", &ServiceConfig); err != nil {
		panic(fmt.Sprintf("load server config failed, err: %v", err))
	}
}
