package redis

import (
	"context"
	"github.com/yulecd/pp-common/plog"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

type Client struct {
	*redis.Client
}

var redisClientMap = map[string]*Client{}
var createDbLock sync.Mutex
var dialTimeout time.Duration

type Config struct {
	Host         string `yaml:"host" json:"host"`
	Port         string `yaml:"port" json:"port"`
	Auth         string `yaml:"auth" json:"auth"`
	Db           int    `yaml:"db" json:"db"`
	PoolSize     int    `yaml:"pool_size" json:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns" json:"min_idle_conns"`

	IdleTimeout    time.Duration `yaml:"idle_timeout" json:"idleTimeout"`
	ConnectTimeout time.Duration `yaml:"connect_timeout" json:"connectTimeout"`
	WriteTimeout   time.Duration `yaml:"write_timeout" json:"writeTimeout"`
	ReadTimeout    time.Duration `yaml:"read_timeout" json:"read_timeout"` //单位ms
}

var RedisConfMap = map[string]Config{}

func GetClient(serverName string) *Client {
	client, _ := redisClientMap[serverName]
	if client != nil {
		return client
	}
	return nil
}

func InitRedisClient(serverName string, conf Config) {
	createDbLock.Lock()
	defer func() {
		createDbLock.Unlock()
	}()

	cnn := func(serverName string) (c *Client, connOk bool) {
		var option *redis.Options
		option = &redis.Options{
			Addr:         conf.Host + ":" + conf.Port,
			Password:     conf.Auth,
			DB:           conf.Db,
			DialTimeout:  dialTimeout,         // 建立新链接的超时时间
			PoolSize:     conf.PoolSize,       // 连接池
			MinIdleConns: conf.MinIdleConns,   // 空闲连接数
			IdleTimeout:  conf.IdleTimeout,    // 空闲连接数存活时间
			PoolTimeout:  conf.ConnectTimeout, // 客户端等待连接的时间
			WriteTimeout: conf.WriteTimeout,   // 写超时时间
			ReadTimeout:  conf.ReadTimeout,    // 读超时时间
		}

		plog.Infof(nil, `redis op %v`, map[string]interface{}{
			`Addr`:        conf.Host + ":" + conf.Port,
			`DB`:          conf.Db,
			`DialTimeout`: conf.ConnectTimeout,
		})
		client := redis.NewClient(option)
		client.AddHook(RedisLogger{})
		if _, err := client.Ping(context.Background()).Result(); err != nil {
			plog.Errorf(nil, "Redis尝试连接失败:%s", err.Error())
		} else {
			connOk = true
		}

		c = &Client{client}
		return
	}

	client, cOk := cnn(serverName)
	if !cOk {
		go func(con *Client) {
			ticker := time.NewTicker(time.Second * 5)
			defer ticker.Stop()
			for {
				<-ticker.C
				if _, err := con.Ping(context.Background()).Result(); err != nil {
					plog.Errorf(nil, "Redis数据库连接已断开:%s, 尝试重新连接", err.Error())
					cnn(serverName)
				}
			}
		}(client)

		redisClientMap[serverName] = nil
	} else {
		redisClientMap[serverName] = client
	}

	return
}
