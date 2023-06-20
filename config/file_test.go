package config

import (
	"io/ioutil"
	"testing"
)

func TestFileConfig(t *testing.T) {

	type RedisConfig struct {
		Name         string `yaml:"name"`
		Host         string `yaml:"host"`
		Port         int    `yaml:"port"`
		Db           int    `yaml:"db"`
		WriteTimeout string `yaml:"auth"`
	}

	fileName := "./app.yaml"
	content := `
app:
  http_port: 80
  read_timeout: 30
  write_timeout: 30

redis:
 pay_server:
  host: j.it603.com
  port: 49184
  auth: 
  db: 
  pool_size: 3
  min_idle_conns: 1
  connect_timeout: 1s
  idle_timeout: 3600s 
  write_timeout: 1200ms
  read_timeout: 100ms
`

	err := ioutil.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		t.Fatal()
	}
	//defer os.Remove(fileName)
	//app := server.Config{}
	c, err := NewConfig(fileName)
	if c == nil {
		t.Error("config is nil")
	}
	if err != nil {
		t.Errorf("new config error: %v", err)
	}
	//err = Load("redis", &redis.RedisConfMap)
	//if err != nil {
	//	t.Errorf("new config error: %v", err)
	//}
	//
	//fmt.Println(redis.RedisConfMap)

	//redisClient := redis.GetClient(`pay_server`)
	//if redisClient != nil {
	//	if _, err = redisClient.Ping(context.Background()).Result(); err != nil {
	//		fmt.Println(err.Error())
	//	}
	//
	//	val, err := redisClient.Get(context.Background(), "yf_test").Result()
	//	fmt.Println(val, err)
	//}
	//
	//time.Sleep(time.Second * 10)
}

//func TestLoad(t *testing.T) {
//	app := server.Config{}
//	_, err := New()
//	if err != nil {
//		t.Errorf("new config error: %v", err)
//	}
//	err = Load("app", &app)
//	if err != nil {
//		t.Errorf("new config error: %v", err)
//	}
//}
