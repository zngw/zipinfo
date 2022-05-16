// @Title 获取外网IP
// @Description $
// @Author  55
// @Date  2022/3/1
package main

import (
	"encoding/json"
	"fmt"
	"github.com/zngw/log"
	"github.com/zngw/zipinfo/config"
	"github.com/zngw/zipinfo/ipinfo"
	"github.com/zngw/zipinfo/listen/cache"
	"github.com/zngw/zipinfo/listen/cb"
	"github.com/zngw/zipinfo/listen/httputil"
	"github.com/zngw/zipinfo/listen/rdb"
	"net"
	"net/http"
	"net/url"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func main() {
	// 加载配置
	port := 8080
	cfg := config.LoadConfig()
	if cfg != nil {
		ipinfo.Init(cfg.Third)
		port = cfg.Port

		// 初始化Redis
		rdb.InitRedis(cfg.Redis, cfg.CacheDay)
	}

	// 监听
	r := &httputil.Router{}
	r.HandleFunc(httputil.UrlBase, query) // 获取IP位置信息

	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
		if err != nil {
			log.Error("sys", "端口 %d 被占用，服务器可能已经运行！！", port)
			panic(err)
		}
	}()
	log.Trace("sys", "----- Start Server @ %d -----", port)

	signal.Ignore(syscall.SIGHUP)
	runtime.Goexit()
}

// 获取IP位置信息
func query(w http.ResponseWriter, r *http.Request) {
	// 获取客户端参数
	param := r.URL.Query()
	ip := param.Get("ip")
	if ip == "" {
		// ip 不存在，查询当前访问ip，跳过验证
		ip = httputil.GetIP(r)
	} else {
		// 验证用户查询权限
		if !httputil.Authentication(param) {
			query, _ := url.QueryUnescape(r.URL.RawQuery)
			log.Trace("net", query+":验证失败！")
			_ = httputil.Send(w, []byte(`{"status":"fail","msg":"验证失败！"}`))
			return
		}
	}

	url := param.Get("url") // 回调地址
	p := param.Get("param") // 透传参数

	// 如果存在回调，先返回，再异步通知
	if url != "" {
		err := httputil.Send(w, []byte(`{"status":"success"}`))
		if err != nil {
			log.Error("sys", err.Error())
			return
		}
	}

	var data []byte = nil

	// 变量请求3次
	for i := 0; i < 3; i++ {
		err, info := getIpInfo(ip)
		if err == nil {
			info.Param = p
			data, err = json.Marshal(info)
			if err == nil {
				break
			}
		}

		// 所有接口查询失败，等待1s后继续查询
		time.Sleep(time.Second)
	}

	if data == nil {
		data = []byte(`{"status":"fail","msg":"获取数据失败"}`)
	}

	if url == "" {
		// 不存在回调，直接返回数据
		_ = httputil.Send(w, data)
	} else {
		// 添加到回调队列
		cb.Push(url, 1, data)
	}
}

// 获取ip数据
func getIpInfo(ip string) (err error, info *ipinfo.IpInfo) {
	// 验证IP是否合法
	address := net.ParseIP(ip)
	if address == nil {
		err = fmt.Errorf("ip地址格式不正确")
		return
	}

	err, info = cache.GetIpInfoByRedis(ip)
	if err == nil {
		return
	}

	err, info = ipinfo.GetIpInfo(ip)
	if err == nil {
		cache.SetIpInfoByRedis(ip, info)
	}

	return
}
