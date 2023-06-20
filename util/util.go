package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Bind(obj interface{}, data interface{}) error {
	if !containJsonTag(data) {
		return BindJson(obj, FieldToMap(data, "form"))
	}

	return BindJson(obj, data)
}

// BindJson 建议使用 bind方法 此方法进适用于来源数据没有json tag的时候
func BindJson(obj interface{}, data interface{}) error {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(jsonStr, obj); err != nil {
		return err
	}

	return nil
}

func containJsonTag(obj interface{}) bool {
	if reflect.ValueOf(obj).Kind() == reflect.Slice {
		return true
	}

	objT := reflect.TypeOf(obj)
	objV := reflect.ValueOf(obj)
	switch {
	case IsStruct(objT):
	case IsStructPtr(objT):
		objT = objT.Elem()
		objV = objV.Elem()
	default:
		return false
	}

	for i := 0; i < objT.NumField(); i++ {
		if objV.Field(i).Kind() == reflect.Struct {
			continue
		}
		// just judgement the first tag is contain json
		tag := reflect.StructTag(objT.Field(i).Tag)
		if tag.Get("json") != "" {
			return true
		}
		return false
	}

	return false
}

func FieldToMap(in interface{}, tagFlag string) map[string]interface{} {
	out := make(map[string]interface{})

	inT := reflect.TypeOf(in)
	inV := reflect.ValueOf(in)
	switch {
	case IsStruct(inT):
	case IsStructPtr(inT):
		inT = inT.Elem()
		inV = inV.Elem()
	default:
		return nil
	}

	for i := 0; i < inT.NumField(); i++ {
		if inV.Field(i).Kind() == reflect.Struct {
			out = mergeMap(out, FieldToMap(inV.Field(i).Interface(), tagFlag))
			continue
		}
		var field string
		tag := reflect.StructTag(inT.Field(i).Tag)
		if tag.Get(tagFlag) != "" {
			field = tag.Get(tagFlag)
		}
		if field == "" && tag.Get("json") != "" {
			field = tag.Get("json")
		}
		if field == "" && tag.Get("form") != "" {
			field = tag.Get("form")
		}

		if inV.Field(i).IsZero() {
			continue
		}
		// compatible support default value situation
		field = extractFieldSpec(field)
		out[field] = inV.Field(i).Interface()
	}

	return out
}

func extractFieldSpec(field string) string {
	var spec string
	index := strings.Index(field, ",")
	spec = field
	if index != -1 {
		spec = field[:index]
	}

	return spec
}

func IsStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}

func IsStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func mergeMap(dest, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		dest[k] = v
	}

	return dest
}

func DirExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// ChangeTimeToUTC 修改时区
func ChangeTimeToUTC() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	time.Local = loc
}

// DownloadFile 下载文件 存在会覆盖
func DownloadFile(url string, localPath string, fb func(length, downLen int64)) error {
	var (
		fsize   int64
		buf     = make([]byte, 32*1024)
		written int64
	)
	tmpFilePath := localPath + ".download"
	// 创建一个http client
	client := new(http.Client)
	//client.Timeout = time.Second * 60 //设置超时时间
	// get方法获取资源
	resp, err := client.Get(url)
	if err != nil {
		return err
	}

	// 读取服务器返回的文件大小
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 32)
	if err != nil {
		fmt.Println(err)
	}
	// 是否存在
	_, err = os.Stat(tmpFilePath)
	if !os.IsNotExist(err) {
		// 存在 删除
		err = os.Remove(tmpFilePath)
		if err != nil {
			return err
		}
	}
	_, err = os.Stat(localPath)
	if !os.IsNotExist(err) {
		// 存在 删除
		err = os.Remove(localPath)
		if err != nil {
			return err
		}
	}

	//创建文件
	file, err := os.Create(tmpFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	if resp.Body == nil {
		return errors.New("body is null")
	}
	defer resp.Body.Close()
	//下面是 io.copyBuffer() 的简化版本
	for {
		//读取bytes
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			//写入bytes
			nw, ew := file.Write(buf[0:nr])
			//数据长度大于0
			if nw > 0 {
				written += int64(nw)
			}
			//写入出错
			if ew != nil {
				err = ew
				break
			}
			//读取是数据长度不等于写入的数据长度
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		//没有错误了快使用 callback
		fb(fsize, written)
	}
	if err == nil {
		file.Close()
		err = os.Rename(tmpFilePath, localPath)
		return err
	}
	return nil
}
