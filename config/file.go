package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v2"
)

const devConfigFile = "config/app_dev.yaml"
const configFileFmt = "config/app_%s.yaml"

type ReaderFile struct {
	filePath string
}

func NewFileReader(env string) Reader {
	var filePath string
	switch {
	case env == "":
		filePath = devConfigFile
	case strings.HasSuffix(env, ".yaml") || strings.HasSuffix(env, ".yml"):
		filePath = env
	default:
		filePath = fmt.Sprintf(configFileFmt, env)
	}
	return &ReaderFile{
		filePath: filePath,
	}
}

func (r *ReaderFile) Init() error {
	config.m.Lock()
	defer config.m.Unlock()
	content, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		return err
	}
	if config.c == nil {
		config.c = make(map[string]interface{})
	}

	return yaml.Unmarshal(content, &config.c)
}

func (r *ReaderFile) Watch() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	// defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					err := r.Init()
					if err != nil {
						log.Println("file=%s load err=%s", r.filePath, err.Error())
						continue
					}

					notifyCallback()
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				log.Println("config watcher err=%s", err.Error())
			}
		}
	}()

	err = watcher.Add(r.filePath)
	if err != nil {
		return err
	}

	return nil
}
