// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var Cfg = new(Conf)

type Conf struct {
	Third    []interface{} `json:"third"`    // 三方接口类型
	Port     int           `json:"port"`     // 端口
	Redis    string        `json:"redis"`    // redis
	CacheDay int           `json:"cacheDay"` // 本地缓存多少天
	Free     bool          `json:"free"`     // 是否可以免费使用
}

// 加载默认配置
func init() {
	defCfg := `{
  "port": 8080,
  "free": true,
  "third":
    [
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
}`
	_ = json.Unmarshal([]byte(defCfg), Cfg)
}

// 加载配置
func LoadConfig() {
	raw, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println("无本地配置文件，使用默认配置：")
		return
	}

	err = json.Unmarshal(raw, Cfg)
	if err != nil {
		fmt.Println("解析基本配置文件失败，使用默认配置")
		Cfg = nil
		return
	}

	return
}
