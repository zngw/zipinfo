// @Title
// @Description $
// @Author  55
// @Date  2022/5/12
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Conf struct {
	Third []string `json:"third"` // 三方接口类型

	Port     int    `json:"port"`     // 端口
	Redis    string `json:"redis"`    // redis
	CacheDay int    `json:"cacheDay"` // 本地缓存多少天
}

// 加载配置
func LoadConfig() (cfg *Conf) {
	raw, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("无本地配置文件：", err)
		return
	}

	cfg = new(Conf)
	err = json.Unmarshal(raw, cfg)
	if err != nil {
		fmt.Println("解析基本配置文件失败：%w", err)
		cfg = nil
		return
	}

	return
}
