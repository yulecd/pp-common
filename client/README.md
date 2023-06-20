# http 客户端

实现与trace组件组合传递链路信息

融合多种协议的client

```
client 是调用的入口
http_client 是http协议的实现
options是client的配置项 部分会作用到request上
wrapper是client的中间件 会一次调用返回
request 是一个api请求的http实现
response 是一个api的响应的http实现
```

## 使用

### 基本使用

```golang
# 使用默认client
resp, err := DefaultClient.Get("http://xxxx/xxx/xxx")
resp, err := DefaultClient.Get("/xxx/xxx", client.Option{
    BaseURI: "http://xxxx",
    Query: map[string]interface{}{
        "k1": "v1",
    }
})
resp, err := DefaultClient.Post("/xxx/xxx", client.Option{
    BaseURI: "http://xxxx",
    JSON: struct {
        Key1 string `json:"key1"`,
    },
})
# 针对某个客户端设置一些选项
resp, err := NewClient(
    Timeout(time.Second*10),
    Retries(2),
).Get("http://xxxx/xxx/xxx")
# 全局更改默认选项
client.WithDefaultRetries(2)
client.WithDefaultTimeout(time.Second * 10)
```

### Query Params

- **query map**

```golang
resp, err := DefaultClient.Get('/api/user', gozzle.Options{
    Query: map[string]interface{}{
        "key1": "value1",
        "key2": []string{"value1", "value2"},
    },
})
if err != nil {
    log.error("call user service err:", err.Error())
}
```

- **query string**

```golang
resp, err := DefaultClient.Get('/api/user2', gozzle.Options{
    Query: "key1=value2&key2=value21&key2=value22",
})
if err != nil {
    log.error("call user service err:", err.Error())
}
```

### Post Data

- **post form**

```golang
resp, err := DefaultClient.Post("/api/order", gozzle.Options{
    Headers: map[string]interface{}{
        "Content-Type": "application/x-www-form-urlencode",
    },
    FormParams: map[string]interface{}{
        "key1": "value1",
        "key2": []string{"value1", "value2"},
    },
})
if err != nil {
    log.error("call user service err:", err.Error())
}
```

- **post json**

```golang
resp, err := DefaultClient.Post("/api/user", gozzle.Options{
    Headers: map[string]interface{}{
        "Content-Type": "application/json",
    },
    JSON: struct {
        Key1 string `json:"key1"`,
        Key2 []string `json:"key2"`,
    }{"value1", []string{"value1", "value2"}},
})
if err != nil {
    log.error("call user service err:", err.Error())
}
```

### Request Header

```golang
resp, err := DefaultClient.Post("/api/anything", gozzle.Options{
    Headers: map[string]interface{}{
        "X-Virtual-Env": "feature-xxx",
    },
})
if err != nil {
    log.error("call user service err:", err.Error())
}
```

### Timeout

```golang
resp, err := DefaultClient.Post("/api/anything", gozzle.Options{
    Timeout: 10,
})
if err != nil {
    log.error("call user service err:", err.Error())
}
```

### Response

```golang
resp, err := DefaultClient.Post("/api/anything", gozzle.Options{
    Headers: map[string]interface{}{
        "X-Virtual-Env": "feature-xxx",
    },
})
if err != nil {
    log.error("call user service err:", err.Error())
}

body, err := resp.GetBody()
if err != nil {
    log.error("xxx")
}

contents := body.GetContents()

status := resp.GetStatusCode()

errMsg := resp.GetReasonPhrase()
```

## 支持的功能

- http正常的restful格式请求
- `retry` 机制
- `wrapper` 中间件机制