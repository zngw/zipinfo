# 获取IP归属地接口

本程序集合了 淘宝、ip.useragentinfo.com、网易126.net、百度、太平洋、ip-api.com等免费IP归属地Api接口，将各网页提供的数据以json格式时实或异步返回。  

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

如果config.json配置文件不存在时，会有默认8080端口，无redis缓存运行

```json
{
  "port": 8080,
  "redis_": "redis://:123455@127.0.0.1:6379",
  "cacheDay": 30,
  "third": ["TaoBao","UserAgentInfo","Net126","BaiDu","PcOnline","IpApi"]
}
```

* port： 监听端口
* redis：redis链接，配置不存在或为空时，不启用redis缓存
* cacheDay： 缓时天数
* third：启用第三方ip归属地查询顺序，目前支持：
> TaoBao: 淘宝  
> UserAgentInfo: ip.useragentinfo.com  
> Net126: 网易126.net  
> BaiDu： 百度  
> PcOnline： 太平洋  
> IpApi： ip-api.com  

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
// 不初始化会安默认顺序判断，
var defaultEnable = []string{"TaoBao","UserAgentInfo","Net126","BaiDu","PcOnline","IpApi"}
ipinfo.Init(defaultEnable)

// 使用
err, info = ipinfo.GetIpInfo(ip)
```