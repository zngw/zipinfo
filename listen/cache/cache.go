// @Title
// @Description $
// @Author  55
// @Date  2022/5/16
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zngw/zipinfo/ipinfo"
	"github.com/zngw/zipinfo/listen/rdb"
)

// 从Redis缓存在获取数据
func GetIpInfoByRedis(ip string) (err error, info *ipinfo.IpInfo) {
	if rdb.Redis == nil {
		err = fmt.Errorf("no cache")
		return
	}

	data, err := rdb.Redis.Get(context.Background(), "ipinfo:"+ip).Result()
	if err != nil {
		return
	}

	info = new(ipinfo.IpInfo)
	err = json.Unmarshal([]byte(data), info)
	return
}

// 写Redis缓存
func SetIpInfoByRedis(ip string, info *ipinfo.IpInfo) {
	if rdb.Redis == nil {
		return
	}

	info.Param = ""
	data, err := json.Marshal(info)
	if err != nil {
		return
	}

	// redis缓存一个月
	if rdb.CacheTime > 0 {
		rdb.Redis.Set(context.Background(), "ipinfo:"+ip, data, rdb.CacheTime)
	}
}
