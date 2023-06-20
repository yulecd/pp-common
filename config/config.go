package config

import (
	"bytes"
	"fmt"
	"log"
	"sync"

	"gopkg.in/yaml.v2"
)

const (
	AppEnvName = "APP_ENV"
	DevEnv     = "dev"
	TestEnv    = "test"
	PreEnv     = "pre"
	ProdEnv    = "prod"
)

type Callback func(name string)

var callbackTable map[string][]Callback
var callbackTableMutex sync.Mutex

// Reader ConfigReader reads config from any source
type Reader interface {
	Init() error  // read config to config.c
	Watch() error // sends config reread signal
}

type Config struct {
	reader Reader

	data       *bytes.Buffer
	c          map[string]interface{} // config content
	notifyList map[string][]Callback  //
	m          sync.Mutex
}

var config *Config

// Load is used for loading config with name to its own struct
// for example mysql component will call config.Load("mysql", configObjPointer) to load config to a config object
// // app.yaml
// mysql:
//  xzweb:
//    master:
//      - dsn1
//    slaver:
//      - dsn2
//
// // test.go
// type Mysql struct {
//    Master []string `yaml:"master"`
//    Slaver []string `yaml:"slaver"`
// }
//
// mysqlConfig := make(map[string]Mysql)
// Load("mysql", &mysqlConfig)
func Load(name string, configObj interface{}) error {
	// 检查初始化
	if config == nil {
		return fmt.Errorf("config module have not been inited lock")
	}

	config.m.Lock()
	defer config.m.Unlock()

	if config.c == nil {
		return fmt.Errorf("config module have not been inited")
	}

	v, ok := config.c[name]
	if !ok {
		return fmt.Errorf("no config with name=%s", name)
	}

	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("yaml marshal name=%s err=%w", name, err)
	}

	err = yaml.Unmarshal(data, configObj)
	if err != nil {
		return fmt.Errorf("yaml unmarshal name=%s err=%w", name, err)
	}

	return nil
}

// LoadWithCallback is like Load but can have a Callback func, it will be called when config changed.
func LoadWithCallback(name string, configObj interface{}, cb Callback) error {
	if err := Load(name, configObj); err != nil {
		return err
	}

	callbackTableMutex.Lock()
	defer callbackTableMutex.Unlock()

	if config.notifyList == nil {
		config.notifyList = make(map[string][]Callback)
	}

	if _, ok := config.notifyList[name]; !ok {
		config.notifyList[name] = make([]Callback, 0)
	}

	config.notifyList[name] = append(callbackTable[name], cb)
	return nil
}

func notifyCallback(names ...string) {
	callbackTableMutex.Lock()
	defer callbackTableMutex.Unlock()
	if len(names) > 0 {
		for _, name := range names {
			if callback, ok := config.notifyList[name]; ok {
				notify(callback, name)
			}
		}
		return
	}

	for name, cbArr := range config.notifyList {
		notify(cbArr, name)
	}
}

func notify(cbArr []Callback, name string) {
	for _, cb := range cbArr {
		cb(name)
	}
}

type Options struct {
	env    string
	reader Reader
}

type Option func(options *Options)

func New(opt ...Option) (*Config, error) {
	opts := &Options{}

	for _, o := range opt {
		o(opts)
	}

	config = &Config{
		data:       nil,
		c:          nil,
		notifyList: nil,
	}

	if opts.reader != nil {
		config.reader = opts.reader
	} else {
		config.reader = NewFileReader(opts.env)
	}

	if config.c == nil {
		config.c = make(map[string]interface{})
		config.data = &bytes.Buffer{}
	}
	if err := config.reader.Init(); err != nil {
		return nil, err
	}
	go func() {
		if err := config.reader.Watch(); err != nil {
			log.Println(err)
		}
	}()
	return config, nil
}

func NewConfig(env string) (*Config, error) {
	return New(WithEnv(env))
}

func WithEnv(env string) Option {
	return func(options *Options) {
		options.env = env
	}
}

func WithReader(reader Reader) Option {
	return func(options *Options) {
		options.reader = reader
	}
}
