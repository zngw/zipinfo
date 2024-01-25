# 获取IP归属地接口

本程序集合了 淘宝、ip.useragentinfo.com、网易126.net、百度、太平洋、ip-api.com等免费IP归属地Api接口，ip138收费IP归属地Api接口，将各网页提供的数据以json格式时实或异步返回。  

## 1.直接访问

返回当前IP归属地信息

```url
http://url:port
```

## 2. 带ip参数

返回查询ip归属地信息

```url
http://url:port?ip=123.123.123.123
```

## 3. 带回调地址和透穿参数

查询结果以post形式回调到url地址中，这里的url需要使用urlEncode编码

```url
http://url:port/?ip=36.18.115.152&url=http%3A%2F%2F192.168.0.1%3A8080%2Fipcb&param=123
```

## 4. 启动用户验证功能

在原参数上加上了user和key，user为user.json文件配置的，key为所有url参数按字母排序加上user.json配置的对应的key的小写md5值.

如果存在url，url则是非urlEncode编码前的值做md5

```url
http://url:port/?ip=36.18.115.152&url=http%3A%2F%2F192.168.0.1%3A8080%2Fipcb&param=id123&user=tmp&key=dec33e5b288d5618de7d4de51acd0b62
```

key计算=`md5(36.18.115.152id123http://192.168.0.1:8080/ipcbtmpf1ae62e1c76c5150f9b0d7e17db95dac)`


# 配置说明

## 1. 基本配置，config.json

如果config.json配置文件不存在时，默认配置

```json
{
  "port": 8080,
  "redis": "redis://:123456@127.0.0.1:6379",
  "cacheDay": 30,
  "free": true,
  "third":
  [
    {
      "name": "zipinfo",
      "free": true,
      "url": "http://127.0.0.1:80",
      "user": "tmp",
      "token": "f1ae62e1c76c5150f9b0d7e17db95dac"
    },
    {
      "name": "ip138",
      "free": false,
      "url": "https://api.ipshudi.com/ip/?ip=%s&datatype=jsonp",
      "token": "ip138提供的token"
    },
    {
      "name": "baidu",
      "free": true,
      "url": "https://opendata.baidu.com/api.php?query=%s&co=&resource_id=6006&oe=utf8"
    },
    {
      "name": "taobao",
      "free": true,
      "url": "https://ip.taobao.com/outGetIpInfo?ip=%s&accessKey=alibaba-inc"
    },
    {
      "name": "useragentinfo",
      "free": true,
      "url": "https://ip.useragentinfo.com/json?ip=%s"
    },
    {
      "name": "ipapi",
      "free": true,
      "url": "http://ip-api.com/json/%s?lang=zh-CN"
    },
    {
      "name": "pconline",
      "free": true,
      "url": "https://whois.pconline.com.cn/ipJson.jsp?ip=%s&json=true"
    }
  ]
}
```

* port： 监听端口
* redis：redis链接，配置不存在或为空时，不启用redis缓存
* cacheDay： 缓时天数
* free：是否开启免费查询功能，免费查询只能调用free标志为true的第三方接口
* third：启用第三方ip归属地查询顺序，目前支持：
> name、free、url为必须字段   
> free: 该第三方平台是否为免费查询用户使用  
> url: 第三方平台的请求地址  
> name：第三方接口
>> zipinfo: 同款软件之间获取IP   
>> ip138: iP138查询网
>> baidu： 百度  
>> taobao: 淘宝  
>> useragentinfo: ip.useragentinfo.com  
>> pconline： 太平洋  
>> ipapi： ip-api.com

## 2. 权限用户配置, user.json

kv模式的键值对，key是用户，value是验证密钥  
user.json不存在或无配置时，以无验证模式启动

```json
{
  "tmp": "f1ae62e1c76c5150f9b0d7e17db95dac"
}
```

# 返回数据

查询的结果都以UTF-8编码的json数据格式返回

```json
{
  "status":"success",
  "type":"tb",
  "country":"中国",
  "region":"浙江",
  "city":"杭州",
  "isp":"电信",
  "query":"36.18.115.152",
  "param":"id123"
}
```

* status: 查询状态，success为成功，fail是失败
* type: 调用的第三方查询接口类型。
* country: 国家
* region: 省
* city: 市
* isp： 运营商
* query: 查询的IP
* param: 透传参数

# 其他程序直接引用ip归属地

## 安装

```
go get -u github.com/zngw/zipinfo
```

## 导入

```go
import github.com/zngw/zipinfo/ipinfo
```

## 初始化

```go

// 使用
err, info = ipinfo.GetIpInfo(ip)
```