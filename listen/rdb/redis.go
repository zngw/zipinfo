// @Title Redis
// @Description $
// @Author  55
// @Date  2021/9/17
package rdb

import (
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

// Redis
var Redis *redis.Client = nil
var CacheTime = time.Hour * 24 * 30

// 连接Redis
func InitRedis(url string, cacheDay int){
	if url == "" {
		return
	}

	// 解析url协议
	if strings.HasPrefix(url, "redis://") {
		url = url[8:]
	}

	// 解析密码
	var pwd string
	if c := strings.Index(url, "@"); c != -1 {
		pair := strings.SplitN(url[:c], ":", 2)
		if len(pair) > 1 {
			pwd = pair[1]
		}
		url = url[c+1:]
	}

	CacheTime = time.Hour * 24 * time.Duration(cacheDay)

	// 连接
	Redis = redis.NewClient(&redis.Options{
		Addr:     url,
		Password: pwd,
		DB:       0,
	})

	return
}
